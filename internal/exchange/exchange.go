package exchange

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/zaffka/testio/config"
	"github.com/zaffka/testio/internal/exchange/handlers"
	"go.uber.org/zap"
)

var ErrNoExchange = errors.New("no exchange matching a config")

type Opts struct {
	Ctx           context.Context
	Logger        *zap.Logger
	Configuration interface{}
}

func New(opts Opts) (func() error, error) {
	var teardownFn func() error
	switch cfg := opts.Configuration.(type) {
	case config.CryptologyMarket:
		u, err := url.Parse(cfg.ConnURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse connection url: %w", err)
		}

		c := handlers.Cryptology{
			ConnURL:      u.String(),
			WSPingPeriod: cfg.PingPeriod,
			Instruments:  cfg.Instruments,
			L:            opts.Logger,
		}

		teardownFn, err = c.Dial(opts.Ctx)
		if err != nil {
			return nil, fmt.Errorf("dial failed: %w", err)
		}

		go c.Handle(opts.Ctx)
	default:
		return nil, ErrNoExchange
	}

	return teardownFn, nil
}
