package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/cyberax/go-dd-service-base/utils"
	. "github.com/cyberax/go-dd-service-base/visibility"
	"github.com/cyberax/go-dd-service-base/visibility/zaputils"
	"github.com/stretchr/testify/assert"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"net/http"
	"os"
	"os/user"
	"time"
	"xyzzy/gen/xyzzy"
)

const AppName = "Xyzzy.Canary"

type TestFailed struct{}

type Client struct {
	Api   xyzzy.XyzzyData
	Stats statsd.ClientInterface
}

type tester struct {
	ctx context.Context
}

var _ assert.TestingT = &tester{}

func (c *tester) Errorf(format string, args ...interface{}) {
	CLS(c.ctx, zap.AddCallerSkip(3)).Errorf(format, args)
	panic(&TestFailed{})
}

func T(ctx context.Context) assert.TestingT {
	return &tester{ctx: ctx}
}

func MakeClient(ctx context.Context, envName, serverUrl string) *Client {
	res := &Client{}

	hooks := &twirp.ClientHooks{}
	key := os.Getenv("XYZZY_TOKEN")
	if key != "" {
		hooks.RequestPrepared = func(ctx context.Context, request *http.Request) (context.Context, error) {
			request.Header.Add("Authorization", "Bearer "+key)
			return ctx, nil
		}
	}

	cli := WrapTwirpClientDef(http.DefaultClient, utils.ToSnakeCase(AppName, '-'))
	res.Api = xyzzy.NewXyzzyDataProtobufClient(serverUrl, cli, twirp.WithClientHooks(hooks))
	res.Stats = &statsd.NoOpClient{}

	if os.Getenv("DD_AGENT_HOST") != "" {
		// Start the tracer
		options := []tracer.StartOption{
			tracer.WithAnalytics(true),
			tracer.WithServiceName(utils.ToSnakeCase(AppName, '-')),
			tracer.WithGlobalTag("env", envName),
		}
		// Hostname is not always pulled
		if os.Getenv("DD_HOSTNAME") != "" {
			options = append(options,
				tracer.WithGlobalTag("host", os.Getenv("DD_HOSTNAME")))
		}
		tracer.Start(options...)

		// Start the metrics submitter
		statsTags := []statsd.Option{
			statsd.WithNamespace(AppName + "."),
			statsd.WithTags([]string{"env:" + envName}),
		}
		statsCli, err := statsd.New("", statsTags...)
		if err == nil {
			res.Stats = statsCli
		} else {
			CL(ctx).Fatal("Skipped initialization of the stats daemon (no DD_AGENT_HOST?)",
				zap.Error(err))
		}
	}

	return res
}

func (cli *Client) Stop() {
	tracer.Stop()
	_ = cli.Stats.Flush()
	_ = cli.Stats.Close()
}

func ping(ctx context.Context, cli *Client) {
	_, err := cli.Api.Ping(ctx, &xyzzy.PingRequest{})
	assert.NoError(T(ctx), err)

	CL(ctx).Info("Ping done")
}

func (cli *Client) runTests(ctx context.Context) (err error) {
	met := GetMetricsFromContext(ctx)
	met.AddCount("Failure", 0)
	met.AddCount("Success", 0)
	bench := met.Benchmark("Time")
	defer bench.Done()

	defer func() {
		p := recover()
		if p == nil {
			met.AddCount("Success", 1)
			return
		}
		met.AddCount("Failure", 1)

		if _, ok := p.(*TestFailed); ok {
			err = fmt.Errorf("test failed")
		} else {
			panic(p)
		}
	}()

	_ = RunInstrumented(ctx, "ping", cli.Stats, CL(ctx),
		func(c context.Context) error {
			ping(c, cli)
			return nil
		})

	return err
}

func testMain(logger *zap.Logger, envName, serverUrl string, period int64) {
	ctx := ImbueContext(context.Background(), logger)
	cli := MakeClient(ctx, envName, serverUrl)
	defer cli.Stop()

	for ; ; {
		err := RunInstrumented(ctx, "MainTests", cli.Stats, logger, cli.runTests)
		// Only exit with an error if we're not doing a periodical run
		if period == 0 {
			if err != nil {
				os.Exit(1)
			}
			break
		}
		time.Sleep(time.Duration(period) * time.Second)
	}
}

func main() {
	userName := "unknown"
	if u, err := user.Current(); err == nil {
		userName = u.Username
	}

	envNameDef := os.Getenv("ENV_NAME")
	if envNameDef == "" {
		envNameDef = "dev-canary-" + userName
	}

	var debug = flag.Bool("debug", false, "Debug mode")
	var envName = flag.String("env-name", envNameDef, "Environment name")
	var serverUrl = flag.String("server-url", "", "Server URL")
	var period = flag.Int64("period", 0,
		"The period for the scheduled runs (seconds)")

	flag.Parse()

	// Create the logger first, nothing can live without logging!
	var logger *zap.Logger
	if *debug {
		logger = zaputils.ConfigureDevLogger()
	} else {
		logger = zaputils.ConfigureProdLogger()
	}
	//noinspection GoUnhandledErrorResult
	defer logger.Sync()
	old := zap.RedirectStdLog(logger)
	defer old()

	testMain(logger, *envName, *serverUrl, *period)
}
