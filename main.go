package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/zaffka/testio/config"
	"github.com/zaffka/testio/internal/exchange"
	"github.com/zaffka/testio/pkg/metrics"
	"github.com/zaffka/testio/pkg/zaplog"
	"go.uber.org/zap"
)

const (
	serviceName      = "testio"
	microServiceName = "quote-acquirer"

	httpShutdownTimeout = 5 * time.Second
)

var (
	ver = "dev"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err.Error())
	}

	cfg, err := config.New()
	if err != nil {
		panic(err.Error())
	}

	rootCtx, rootCancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer rootCancel()

	logger := zaplog.New(os.Stderr, zaplog.Opts{
		Host:             hostname,
		Service:          serviceName,
		MicroService:     microServiceName,
		Version:          ver,
		Debug:            cfg.Debug,
		IsDevEnvironment: cfg.DevMode,
	})

	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("failed to flush logger data", zap.Error(err))
		}

		logger.Info("the app is finished")
	}()

	logger.Info("starting the app")

	metricServ := metrics.Run(cfg.HTTPPort, logger, rootCancel)
	teardownFn, err := exchange.New(exchange.Opts{
		Ctx:           rootCtx,
		Logger:        logger,
		Configuration: cfg.CryptologyMarket,
	})
	if err != nil {
		logger.Error("failed to init an exchange", zap.Error(err))
	}
	defer func() {
		if err := teardownFn(); err != nil {
			logger.Error("failed to close exchange connection", zap.Error(err))
		}
	}()

	<-rootCtx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer cancel()

	if err := metricServ.Shutdown(ctx); err != nil {
		logger.Error("failed to gracefully shutdown an http server", zap.Error(err))
	}
}
