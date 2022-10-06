package handlers_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zaffka/testio/internal/exchange/handlers"
)

//go:embed test/order_book.json
var orderBook []byte

//go:embed test/ticker.json
var ticker []byte

//go:embed test/unknown_data.json
var unknown []byte

func Test_Response(t *testing.T) {
	t.Run("order_book", func(t *testing.T) {
		r := handlers.Response{}
		err := json.Unmarshal(orderBook, &r)
		require.NoError(t, err)

		dt := r.GetDataType()
		require.Equal(t, handlers.OrderBookDataType, dt)
	})

	t.Run("ticker", func(t *testing.T) {
		r := handlers.Response{}
		err := json.Unmarshal(ticker, &r)
		require.NoError(t, err)

		dt := r.GetDataType()
		require.Equal(t, handlers.QuoteDataType, dt)
	})

	t.Run("unknown", func(t *testing.T) {
		r := handlers.Response{}
		err := json.Unmarshal(unknown, &r)
		require.NoError(t, err)

		dt := r.GetDataType()
		require.Equal(t, handlers.UnknownDataType, dt)

		require.Nil(t, r.Payload.Data.Volumes.Asks)
	})
}
