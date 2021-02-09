package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cyberax/go-dd-service-base/visibility/zaputils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"runtime/pprof"
	"time"
	"xyzzy/server"
)

func main() {
	userName := "unknown"
	if u, err := user.Current(); err == nil {
		userName = u.Username
	}

	envNameDef := os.Getenv("ENV_NAME")
	if envNameDef == "" {
		envNameDef = "dev-xyzzy-"+userName
	}

	dbUrlDef := os.Getenv("DATABASE_URL")
	if dbUrlDef == "" {
		dbUrlDef = fmt.Sprintf("postgres://%s:@?dbname=xyzzy&sslmode=disable", userName)
	}

	var listenAddr = flag.String("listen-interface",
		"[::0]:8080", "The listen address")
	var hostName = flag.String("external-host-name",
		userName, "The external host name of the load balancer")
	var envName = flag.String("env-name", envNameDef,
		"The environment name for the installation")
	var dbUrl = flag.String("db", dbUrlDef,
		"The database connection string (RDS-type or regular)")

	var debug = flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	if *envName == "" {
		panic("Environment name is not set!")
	}

	labelCtx := pprof.WithLabels(context.Background(), pprof.Labels("from", "init"))
	pprof.SetGoroutineLabels(labelCtx)

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

	// Initialize the registry dependencies
	reg := &server.Registry{}
	err := reg.Init(*debug, *envName, *hostName, *dbUrl, logger)
	if err != nil {
		logger.Fatal("failed to initialize the server", zap.Error(err))
	}

	reg.StartBackground()
	defer reg.StopBackground()

	// Create the Echo server
	e := reg.MakeServer()

	// And run it!
	runServer(e, *listenAddr, logger)
}

func runServer(muxer *mux.Router, addr string, logger *zap.Logger) {
	srv := http.Server{Addr: addr, Handler: muxer}

	quit := make(chan os.Signal)
	defer close(quit)

	labelCtx := pprof.WithLabels(context.Background(), pprof.Labels("from", "server"))
	pprof.SetGoroutineLabels(labelCtx)

	// And we serve HTTP until the world ends.
	serverDone := make(chan bool)
	go func() {
		<-quit
		defer close(serverDone)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal(err.Error())
		}
	}()
	signal.Notify(quit, os.Interrupt)

	logger.Info("Starting to listen on", zap.String("addr", addr))
	err := srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	<-serverDone
}
