# Xyzzy
Xyzzy is a service for <something>.

## Just running the client

The CLI client utility can be installed through Brew on macOS and Linux:
```
brew tap cyberax/tools
brew install cyberax-xyzzy-client
```

This will install `xyzzy` binary, that can be used to interact with the service.

# Server development

## Preparing the environment
You need the following software installed to build Xyzzy:

1. Go (1.14 or later)
2. POSIX make
3. Node.JS and yarn
4. PostgreSQL
5. PostGIS
6. Docker (optional)

To install dependencies on macOS do:

```bash
brew install go yarn yarn-completion postgresql postgis
```

Xyzzy uses PostgreSQL, so you need to create the database and the schema:
```
createdb xyzzy
cat schema/database/dbschema.sql | psql xyzzy
``` 

### Private packages

We use a private package on Github during the build, so make sure you have the following
environment variable set:
```
export GOPRIVATE=github.com/cyberax
```

If you're using Goland IDE then set this property in `GoLand -> Preferences... -> Go -> Go Modules`
in the environment box.

## Building and running the server

To prepare to run the server use `make generate-code`. This will create all the required
pre-requisites.

To run the server you need to run the `server/main/main.go` file. During development, you
can run it directly using `go run`. The minimal command line looks like this:

```
go run server/main/main.go --debug --db "<DATABASE_URL>"
```

A typical database URL for the local database is:
```
postgres://<USER>:<PASS>@localhost?dbname=<DB>&sslmode=disable
```

It's possible to use local socket authentication by specifying the URL like this:
`postgres://:?dbname=xyzzy&sslmode=disable`

## Running the client

The client can be run as following, and as usually it supports Bash autocompletion and has
fairly informative command-line help:
```
go run client/main/xyzzy.go
```

For example, to run the Ping command try this:
```
go run client/main/xyzzy.go --server http://localhost:8080 ping
```

## Running the integration tests

You can ensure that the server is running using one of the canaries. To run the Go-language
implementation, use the following command line:
```
go run canary/main/canary.go -debug -server-url http://localhost:8080
```

To run the Typescript-based canary (which would also test the TypeScript binding):
```
cd bindings/typescript
yarn run tsc --sourceMap *.ts
XYZZY_BASE_URL=http://localhost:8080 node examples/canary.js
```

# Tracing and monitoring
There is support for Datadog for metrics and tracing. You can activate it by specifying
the `DD_AGENT_HOST=localhost` environment variable (of course, you need to have the agent running).

DataDog normally connects to Heroku directly to ingest logs. It's cumbersome if you want to debug
the log-related functionality like log scanning or alarming. Fortunately, it's possible to
enable direct log submission to DataDog.

First, turn on the log collector in `/opt/datadog-agent/etc/datadog.yaml` by setting:
```
logs_enabled: true
``` 

Then you can run the server with `DD_TCP_SINK=localhost:10518` environment variable set, which
will redirect logs to DataDog.

The logs, metrics and traces submitted to DataDog are tagged with the environment name tag. This
tag by default is `dev-xyzzy-$USER`, you can override it by specifying the `--env-name` command
line option for the Xyzzy server. Please make sure to not step over the actual
*production/staging/alpha* environments when doing local development.

# Heroku deployment
The code is deployed on Heroku using Docker. There are two containers: Canary and Web. The
Web container is the main application, and the Canary is a test suite that constantly
runs in the background (or as a one-off integration test).

To build the containers and push them to Heroku's repository use:
```
heroku container:push --arg GITHUB_TOKEN --recursive -a "cyberax-xyzzy-<ENV>"
```
`<ENV>` is one of `alpha`/`staging`/`prod`.

After the containers are pushed, Heroku doesn't start them immediately. Use this command
line to do that:
```
heroku container:release web canary -a "cyberax-xyzzy-<ENV>"
```

Normally the Heroku build/released is triggered through the CircleCI workflow, so you
should use the manual release process only in emergency or during development in your own
private app branch.

# Pushing gems and NPMs
To push a Gem you need to do two steps.

Build it:
```
gem build cyberax-xyzzy-client
```

This will produce a file named `cyberax-xyzzy-client-X.Y.Z.gem`, so then you can push it:
```
gem push --key github \
    --host https://rubygems.pkg.github.com/cyberax \ 
    cyberax-xyzzy-client-X.Y.Z.gem
```
