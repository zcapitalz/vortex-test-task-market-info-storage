package domain

import "time"

type HistoryOrder struct {
	ClientName          string    `json:"clinet"`
	ExchangeName        string    `json:"exchangeName"`
	Label               string    `json:"label"`
	Pair                string    `json:"pair"`
	Side                string    `json:"side"`
	Type                string    `json:"type"`
	BaseQty             float64   `json:"baseQty"`
	Price               float64   `json:"price"`
	AlgorithmNamePlaced string    `json:"algorithmNamePlaced"`
	LowestSellPrc       float64   `json:"lowestSellPrc"`
	HighestBuyPrc       float64   `json:"highestBuyPrc"`
	CommissionQuoteQty  float64   `json:"commissionQuoteQty"`
	TimePlaced          time.Time `json:"timePlaced"`
}
