package ws

import "time"

// Config is a structure holding url and a header params needed to configure ws realization.
type Config struct {
	URL        string
	Header     map[string][]string
	PingPeriod time.Duration
}

// New is a constructor function masking web socket realization.
func New(cnf Config) *Gorilla {
	return &Gorilla{
		config:     cnf,
		pingPeriod: cnf.PingPeriod,
		doneCh:     make(chan error, 1),
	}
}
