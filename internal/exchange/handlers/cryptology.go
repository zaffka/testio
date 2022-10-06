package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zaffka/testio/pkg/ws"
	"go.uber.org/zap"
)

const (
	connTimeout     = 5 * time.Second
	spotChannel     = "spotTickers"
	subscriptionFmt = `{"type": "subscribe","channel":"%s.%s","request_id": "%s"}`
)

var subscriptionChannels = [...]string{
	spotChannel,
}

type Cryptology struct {
	ConnURL      string
	WSPingPeriod time.Duration
	Instruments  []string

	L *zap.Logger

	conn *ws.Gorilla
}

func (c *Cryptology) Dial(ctx context.Context) (func() error, error) {
	c.conn = ws.New(ws.Config{URL: c.ConnURL, PingPeriod: c.WSPingPeriod})

	ctx, cancel := context.WithTimeout(ctx, connTimeout)
	defer cancel()

	return c.conn.Dial(ctx)
}

func (c *Cryptology) Handle(ctx context.Context) {
	go c.responseManager(ctx)
	go c.subscriptionManager(ctx)
}

func (c *Cryptology) subscriptionManager(ctx context.Context) {
	for _, instrument := range c.Instruments {
		for _, channel := range subscriptionChannels {
			c.L.Info("subscribing", zap.String("instrument", instrument), zap.String("channel", channel))
			subscriptionID := uuid.New()
			_, err := c.conn.Write([]byte(fmt.Sprintf(subscriptionFmt, channel, instrument, subscriptionID)))
			if err != nil {
				c.L.Warn("failed to send a subscription", zap.String("instrument", instrument))
			}
		}
	}
}

func (c *Cryptology) responseManager(ctx context.Context) {
	log := c.L.With(zap.String("routine_name", "response_manager"))
	log.Info("starting response manager")
	defer log.Info("finishing response manager")

	responses := c.conn.Run(ctx)

	for respBts := range responses {
		resp := Response{}
		err := json.Unmarshal(respBts, &resp)
		if err != nil {
			c.L.Warn("failed to unmarshal response", zap.ByteString("payload", respBts))
		}

		switch resp.Type {
		case "welcomeMessage":
			log = log.With(zap.String("connection_id", resp.Payload.ConnectionID))
			log.Info("successfully connected")
		case "subscribedSuccessful":
			log.Info("subscribed",
				zap.Int64("server_timestamp", resp.ServerTimestamp),
				zap.String("request_id", resp.RequestID),
				zap.String("channel", resp.Payload.Channel))
		case "subscriptionData":
			switch resp.GetDataType() {
			case QuoteDataType:
				c.L.Info("got a quote", zap.Any("data", resp))
			case OrderBookDataType:
				// TODO
			default:
				log.Warn("got a message with unexpected structure", zap.ByteString("payload", respBts))
			}
		default:
			log.Warn("got a message of unknown type", zap.ByteString("payload", respBts))
		}
	}
}
