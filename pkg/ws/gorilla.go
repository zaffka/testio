package ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var keepAlivePayload = []byte{}

// Gorilla uses https://github.com/gorilla/websocket library to handle ws connection.
type Gorilla struct {
	config     Config
	conn       *websocket.Conn
	pingPeriod time.Duration
	mu         sync.Mutex
	doneCh     chan error
}

// Dial calls the remote ws server.
func (g *Gorilla) Dial(ctx context.Context) (func() error, error) {
	var (
		err  error
		resp *http.Response
	)
	g.conn, resp, err = websocket.DefaultDialer.DialContext(ctx, g.config.URL, g.config.Header)
	if errors.Is(err, websocket.ErrBadHandshake) {
		return nil, fmt.Errorf("connection failed: %w: code: %s", err, resp.Status)
	}

	if err != nil {
		return nil, err
	}

	return g.conn.Close, nil
}

// Write satisfies io.Writer interface.
// Use it to write a messages to ws.
func (g *Gorilla) Write(payload []byte) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	wcl, err := g.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, fmt.Errorf("failed to write to ws: %w", err)
	}

	i, err := wcl.Write(payload)
	if err != nil {
		return 0, err
	}

	err = wcl.Close()
	if err != nil {
		return 0, err
	}

	return i, nil
}

// Run executes a handling routine and returns a channel form where raw messages to be read.
func (g *Gorilla) Run(ctx context.Context) <-chan []byte {
	res := make(chan []byte)
	go g.handle(ctx, res)
	go g.ping(ctx)

	return res
}

// Done returns a channel where any handling error to be written.
// In a case of normal interrupting the channel holds a nil-error.
func (g *Gorilla) Done() <-chan error {
	return g.doneCh
}

func (g *Gorilla) handle(ctx context.Context, res chan []byte) {
	defer close(g.doneCh)
	defer g.unsubscribe()

handleLoop:
	for {
		_, message, err := g.conn.ReadMessage()
		if err != nil {
			g.doneCh <- fmt.Errorf("read error: %w", err)

			break handleLoop
		}

		select {
		case <-ctx.Done():
			g.doneCh <- ctx.Err()

			break handleLoop
		default:
			res <- message
		}
	}
}

func (g *Gorilla) ping(ctx context.Context) {
	ticker := time.NewTicker(g.pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := g.writeMsg(websocket.PingMessage, keepAlivePayload); err != nil {
				return
			}
		}
	}
}

func (g *Gorilla) writeMsg(messageType int, data []byte) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.conn.WriteMessage(messageType, data)
}

func (g *Gorilla) unsubscribe() {
	_ = g.writeMsg(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
