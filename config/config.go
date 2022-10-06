package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config struct stores service starting params.
type Config struct {
	Debug    bool `envconfig:"DEBUG"`
	DevMode  bool `envconfig:"DEV"`
	HTTPPort int  `envconfig:"HTTP_PORT" default:"9080"`

	CryptologyMarket
}

// CryptologyMarket represents Cryptology exchange connection params set.
type CryptologyMarket struct {
	ConnURL     string        `envconfig:"CRYPTOLOGY_CONN_URL" default:"wss://octopus-sandbox.cryptology.com/v1/connect"`
	PingPeriod  time.Duration `envconfig:"CRYPTOLOGY_PING_PERIOD" default:"10s"`
	Instruments []string      `envconfig:"CRYPTOLOGY_INSTRUMENTS" default:"ETH_EUR"`
}

// New initiates a new config.
func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("failed to read a config: %w", err)
	}

	return cfg, nil
}
