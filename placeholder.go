// +build os360

package placeholder

// A placeholder file to prevent go mod tidy from removing this module.
// It's tagged with a non-existent build qualifier to prevent it from fouling tests.

//noinspection GoInvalidPackageImport
import (
	_ "github.com/iancoleman/strcase"
	_ "github.com/jstemmer/go-junit-report"
	_ "github.com/kyleconroy/sqlc"
	_ "github.com/twitchtv/twirp-ruby/protoc-gen-twirp_ruby"
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/twpayne/go-geom"
	_ "go.larrymyers.com/protoc-gen-twirp_typescript"
	_ "github.com/golang/protobuf"
)
