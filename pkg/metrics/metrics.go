package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	ReadTO       = 10 * time.Second
	WriteTO      = 10 * time.Second
	MaxHeaderBts = 1 << 20

	defaultHTTPAddrFmt = "0.0.0.0:%d"
)

func Run(port int, log *zap.Logger, rootCtxCnFn context.CancelFunc) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.Handle("/metrics", promhttp.Handler())

	httpServ := &http.Server{
		Addr:           fmt.Sprintf(defaultHTTPAddrFmt, port),
		Handler:        mux,
		ReadTimeout:    ReadTO,
		WriteTimeout:   WriteTO,
		MaxHeaderBytes: MaxHeaderBts,
	}

	go func() {
		err := httpServ.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return
		}

		if err != nil {
			log.Error("failed to start HTTP server, cancelling root context", zap.Error(err))
			rootCtxCnFn()
		}
	}()

	return httpServ
}
