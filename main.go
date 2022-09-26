package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tickTockPeriodStr = "5s"
	httpPort          = ":9080"
	metricsPath       = "/metrics"
	tickTockMsgFmt    = "tick tock #%d"
)

func main() {
	tickTockPeriod, err := time.ParseDuration(tickTockPeriodStr)
	if err != nil {
		log.Printf("failed to parse period, tickTockPeriodStr variable holds %s", tickTockPeriodStr)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go runHTTPServ(cancel)                 // with prometheus default handler
	go tickTockLogger(ctx, tickTockPeriod) // just to see something in logs

	<-ctx.Done()
}

func tickTockLogger(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-ticker.C:
			log.Printf(tickTockMsgFmt, i)
		case <-ctx.Done():
			return
		}

		i++
	}
}

func runHTTPServ(ctxCancelFn context.CancelFunc) {
	log.Printf("Starting HTTP server, port %s, metrics path %s", httpPort, metricsPath)

	http.Handle(metricsPath, promhttp.Handler())
	err := http.ListenAndServe(httpPort, nil)
	if errors.Is(err, http.ErrServerClosed) {
		return
	}

	if err != nil {
		ctxCancelFn()
	}
}
