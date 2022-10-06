package handlers

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type DataType uint8

const (
	UnknownDataType DataType = iota
	QuoteDataType
	OrderBookDataType
)

type Response struct {
	Type            string `json:"type,omitempty"`
	ServerTimestamp int64  `json:"server_timestamp,omitempty"`
	RequestID       string `json:"request_id,omitempty"`
	Payload         struct {
		Channel      string `json:"channel,omitempty"`
		Message      string `json:"message,omitempty"`
		ConnectionID string `json:"connection_id,omitempty"`
		Data         struct {
			TradePair string `json:"trade_pair,omitempty"`

			// quote's data
			BestBid string `json:"best_bid,omitempty"`
			BestAsk string `json:"best_ask,omitempty"`

			// orderbook's data
			Volumes Volumes `json:"volumes,omitempty"`
		} `json:"data,omitempty"`
	} `json:"payload,omitempty"`
}

type Volumes struct {
	Asks map[string]string `json:"asks,omitempty"`
}

func (r *Response) GetDataType() DataType {
	if r.Payload.Data.BestAsk != "" && r.Payload.Data.BestBid != "" {
		return QuoteDataType
	}

	if r.Payload.Data.Volumes.Asks != nil {
		return OrderBookDataType
	}

	return UnknownDataType
}
