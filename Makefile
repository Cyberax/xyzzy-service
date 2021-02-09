# The build code for the Apollo project
GO:=go
PROTOC:=protoc

ALL: generate-code

PWD=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

validate_path:=$(shell $(GO) list -f '{{.Dir}}' -m github.com/envoyproxy/protoc-gen-validate)
gogo_path:=$(shell $(GO) list -f '{{.Dir}}' -m github.com/gogo/protobuf)

# Command to copy the proto files
rsync=rsync -r --no-perms  --include='*.proto' --include='*/' --exclude='*' --chmod=ugo=rwX

# Map the standard Protobuf packages to their Gogo versions
mappings += gogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto
mappings += google/api/annotations.proto=github.com/gogo/googleapis/google/api
mappings += google/api/http.proto=github.com/gogo/googleapis/google/api
mappings += google/protobuf/any.proto=github.com/gogo/protobuf/types
mappings += google/protobuf/duration.proto=github.com/gogo/protobuf/types
mappings += google/protobuf/empty.proto=github.com/gogo/protobuf/types
mappings += google/protobuf/struct.proto=github.com/gogo/protobuf/types
mappings += google/protobuf/timestamp.proto=github.com/gogo/protobuf/types
mappings += google/protobuf/wrappers.proto=github.com/gogo/protobuf/types

joinlist=$(if $(word 2,$2),$(firstword $2)$1$(call joinlist,$1,$(wordlist 2,$(words $2),$2)),$2)
comma = ,

generate-code: generate-rpc generate-bindings generate-db
	# Generate code

gen/exe/.timestamp: go.mod go.sum
	go build -o gen/exe/protoc-gen-gogofast github.com/gogo/protobuf/protoc-gen-gogofast
	go build -o gen/exe/protoc-gen-lv github.com/cyberax/go-dd-service-base/twirpwrap
	go build -o gen/exe/protoc-twirp github.com/twitchtv/twirp/protoc-gen-twirp
	go build -o gen/exe/protoc-twirp-ts go.larrymyers.com/protoc-gen-twirp_typescript
	go build -o gen/exe/protoc-validate github.com/envoyproxy/protoc-gen-validate
	go build -o gen/exe/sqlc github.com/kyleconroy/sqlc/cmd/sqlc
	go build -o gen/exe/protoc-twirp-ruby github.com/twitchtv/twirp-ruby/protoc-gen-twirp_ruby
	touch gen/exe/.timestamp

# Golang RPC bindings
generate-rpc := gen/xyzzy/.timestamp
generate-rpc: $(generate-rpc)
gen/xyzzy/.timestamp: rpc/xyzzy/service.proto gen/exe/.timestamp
	mkdir -p gen
	echo "Using validate from $(validate_path), gogo from $(gogo_path)"
	$(PROTOC) --proto_path=rpc \
	    -I $(gogo_path) -I $(validate_path) \
		--plugin=protoc-gen-gogo=gen/exe/protoc-gen-gogofast \
		--plugin=protoc-gen-twirp=gen/exe/protoc-twirp \
		--plugin=protoc-gen-validate=gen/exe/protoc-validate \
		--plugin=protoc-gen-lv=gen/exe/protoc-gen-lv \
		--twirp_out=gen \
		--gogo_out='$(call joinlist,$(comma),plugins=-grpc $(addprefix M,$(mappings))):$(PWD)/gen' \
		--validate_out="lang=gogo:gen" \
		--lv_out="gen" \
		rpc/xyzzy/service.proto
	touch gen/xyzzy/.timestamp

# Typescript bindings
generate-bindings: bindings/typescript/.timestamp bindings/ruby/.timestamp
bindings/typescript/.timestamp: rpc/xyzzy/service.proto gen/exe/.timestamp
	# Now generate the Twirp top-level wrapper
	$(PROTOC) --proto_path=rpc \
	    -I $(gogo_path) -I $(validate_path) \
		--plugin=protoc-gen-typescript=gen/exe/protoc-twirp-ts \
		--typescript_out=library=pbjs:$(PWD)/bindings/typescript \
		rpc/xyzzy/service.proto
	# Then generate the Protobuf TypeScript wrappers
	cd bindings/typescript && yarn install
	cd bindings/typescript && yarn run pbjs \
		-w commonjs --sparse --no-delimited --no-verify --no-convert --sparse \
		-p $(gogo_path) -p $(validate_path) \
		-t static-module -o service.pb.js $(PWD)/rpc/xyzzy/service.proto
	cd bindings/typescript && yarn run pbts -o service.pb.d.ts service.pb.js
	cd bindings/typescript && yarn run tsc --sourceMap
	touch bindings/typescript/.timestamp

bindings/ruby/.timestamp: rpc/xyzzy/service.proto gen/exe/.timestamp
	mkdir -p $(PWD)/bindings/ruby/lib/cyberax-xyzzy-client
	$(PROTOC) --proto_path=rpc/xyzzy \
	    -I $(gogo_path) -I $(validate_path) \
		--plugin=protoc-gen-twirp_ruby=gen/exe/protoc-twirp-ruby \
		--ruby_out=$(PWD)/bindings/ruby/lib/cyberax-xyzzy-client \
		--twirp_ruby_out=$(PWD)/bindings/ruby/lib/cyberax-xyzzy-client \
		rpc/xyzzy/service.proto

client: xyzzy

run-ts-canary:
	cd $(PWD)/bindings/typescript && node examples/canary.js

xyzzy: $(generate-rpc)
	$(GO) build client/main/xyzzy.go

# Generated type-safe queries
generate-db := gen/db/.timestamp
generate-db: $(generate-db)
db_deps := $(shell find $(PWD)/schema -name \*.sql -print)
gen/db/.timestamp: $(PWD)/schema/sqlc.yaml $(db_deps)
	echo "Generating typesafe queries"
	cd $(PWD)/schema; $(PWD)/gen/exe/sqlc generate
	cd $(PWD)/schema; $(PWD)/gen/exe/sqlc compile
	touch $(PWD)/gen/db/.timestamp

test:
	@$(GO) test -race ./... | $(GO) run github.com/jstemmer/go-junit-report

clean:
	rm -Rf gen
