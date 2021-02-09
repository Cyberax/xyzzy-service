package server

import (
	"compress/gzip"
	"context"
	"database/sql"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/NYTimes/gziphandler"
	"github.com/cyberax/go-dd-service-base/utils"
	. "github.com/cyberax/go-dd-service-base/visibility"
	"github.com/cyberax/go-dd-service-base/visibility/tracedsql"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/lib/pq"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"xyzzy/gen/db"
	xyzzy_api "xyzzy/gen/xyzzy"
)

const AppName = "Xyzzy"

type Registry struct {
	// General settings
	Debug            bool
	EnvName          string
	HostNameOverride string
	MagicKey         string
	BrowserPath      string

	// Logging and metrics
	RootLogger *zap.Logger
	Metrics    statsd.ClientInterface

	// Connectivity
	Database  *sql.DB

	// Background stuff
	Processes *ProcessRegistry

	// API
	ApiImpl xyzzy_api.XyzzyData

	// Managers
	//AuthManager *managers.AuthManager
}

func (r *Registry) Init(debug bool, envName string, extAddr string, dbUrl string,
	rootLogger *zap.Logger) error {

	r.Debug = debug
	r.EnvName = envName
	r.RootLogger = rootLogger
	r.HostNameOverride = extAddr

	cli, err := SetupTracing(context.Background(), AppName, envName, rootLogger)
	if err != nil {
		return err
	}
	r.Metrics = cli

	err = r.InitDatabase(dbUrl)
	if err != nil {
		return err
	}

	// Managers
	//r.AuthManager = managers.NewAuthManager(r.Database)

	r.MakeApi()

	r.Processes = NewProcessRegistry(r.Metrics, r.RootLogger)

	r.MagicKey = os.Getenv("XYZZY_TOKEN")

	return nil
}

func (r *Registry) MakeApi() {
	dbx := sqlx.NewDb(r.Database, "postgres")
	dbx.Mapper = reflectx.NewMapperFunc("", func(s string) string {
		return utils.ToSnakeCase(s, '_')
	})

	tapi := &XyzzyApi{
		Database: dbx,
		Q:        db.New(dbx),
	}
	r.ApiImpl = xyzzy_api.NewXyzzyDataLogValidate(tapi)
}

func (r *Registry) InitDatabase(dbUrl string) error {
	sqltrace.Register("pq", &pq.Driver{},
		sqltrace.WithServiceName(utils.ToSnakeCase(AppName, '-')+".db"))

	return RunInstrumented(context.Background(), "DatabaseInit",
		r.Metrics, r.RootLogger, func(ctx context.Context) error {
			CL(ctx).Info("Preparing the database connection")

			sslPath, err := tracedsql.MakeCaCertFile(tracedsql.Rds2019)
			if err != nil {
				return err
			}

			connector, err := tracedsql.MakePgConnector(ctx, dbUrl, sslPath, aws.Config{})
			if err != nil {
				return err
			}
			r.Database = sqltrace.OpenDB(connector)
			r.Database.SetConnMaxLifetime(10 * time.Minute)

			// Ping to make sure!
			CL(ctx).Info("Checking database connection")
			for i := 0; i<10; i++ {
				err = r.Database.PingContext(ctx)
				if err == nil {
					return nil
				} else {
					CL(ctx).Info("Connection failed, retrying")
					time.Sleep(2 * time.Second)
				}
			}
			return err
		})
}

func (r *Registry) MakeServer() *mux.Router {
	// Create a basic router and set up logging hooks
	twirpHandler := xyzzy_api.NewXyzzyDataServer(r.ApiImpl,
		MakeTraceHooks(AppName, r.Metrics))

	gorilla := NewTracedGorilla(twirpHandler, r.RootLogger, aws.Float64(1.0), aws.Float64(1.0))

	muxer := mux.NewRouter()
	gorilla.AttachGorillaToMuxer(muxer)
	// Attach the compression handler after the main HTTP middleware
	gz, _ := gziphandler.NewGzipLevelAndMinSize(gzip.BestSpeed, 3*1024*1024)
	muxer.Use(func(handler http.Handler) http.Handler {
		return gz(handler)
	})

	// For the healthcheck we really want to allow simple GETs, so we fake the
	// body and the method.
	muxer.Path("/twirp/xyzzy.XyzzyData/Ping").Methods("GET").HandlerFunc(func(
		writer http.ResponseWriter, request *http.Request) {
		request.Method = "POST"
		request.Header.Set("content-type", "application/json")
		request.Body = ioutil.NopCloser(strings.NewReader("{}"))
		twirpHandler.ServeHTTP(writer, request)
	})

	// The default handler
	if r.BrowserPath != "" {
		muxer.PathPrefix("/").Handler(http.FileServer(http.Dir(r.BrowserPath)))
	}
	return muxer
}

func (r *Registry) runAuth(ctx context.Context) (context.Context, error) {
	if r.MagicKey == "" {
		return ctx, nil
	}

	headers, b := GetHttpRequestHeader(ctx)
	if !b {
		panic("No request headers")
	}
	auth := headers.Get("Authorization")
	if !strings.Contains(auth, r.MagicKey) {
		return ctx, twirp.NewError(twirp.Unauthenticated, "no authorization header")
	}
	return ctx, nil
}

func (r *Registry) StartBackground() {
}

func (r *Registry) StopBackground() {
	r.Processes.Close()
	TearDownTracing(context.Background(), r.Metrics)
}
