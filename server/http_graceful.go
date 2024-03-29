package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ehwjh2010/viper"
	"github.com/ehwjh2010/viper/constant"
	"github.com/ehwjh2010/viper/log"
	"github.com/ehwjh2010/viper/verror"
	wrapErrs "github.com/pkg/errors"
)

var (
	InvalidHttpConf   = errors.New("invalid http config")
	InvalidHttpEngine = errors.New("invalid http engine")
)

func getStopChan() chan os.Signal {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, constant.ListenSignals...)
	return stopChan
}

func getErrChan() chan error {
	errChan := make(chan error)
	return errChan
}

func GraceHttpServer(graceHttp *GraceHttp) error {
	log.Infof(viper.SIGN + "\n" + "Viper Version: " + viper.VERSION)

	if graceHttp == nil {
		return InvalidHttpConf
	}

	if graceHttp.Engine == nil {
		return InvalidHttpEngine
	}

	//Invoke OnStartUp
	log.Debugf("execute on startup functions")
	if err := graceHttp.ExecuteStartUp(); err != nil {
		return wrapErrs.Wrap(err, "on start function occur err")
	}

	defer func() {
		log.Debugf("execute on shutdown functions")
		if closeErrs := graceHttp.ExecuteStartUp(); closeErrs != nil {
			log.Errors(closeErrs)
		}
	}()

	srv := &http.Server{
		Addr:    graceHttp.Addr,
		Handler: graceHttp.Engine,
	}

	stopChan := getStopChan()
	errChan := getErrChan()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	var grpcFlag bool
	if graceHttp.GraceGrpc != nil {
		graceGrpcServer(graceHttp.GraceGrpc, errChan)
		grpcFlag = true
	}

	select {
	case <-stopChan:
		var multiErr verror.MultiErr
		log.Infof("Shutting down gracefully")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(graceHttp.WaitSecond)*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			multiErr.AddErr(wrapErrs.Wrap(err, "stop server failed!!!"))
		}

		if grpcFlag {
			graceHttp.GraceGrpc.Server.GracefulStop()
		}

		return multiErr.AsStdErr()
	case e := <-errChan:
		return e
	}
}
