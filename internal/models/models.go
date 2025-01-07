package models

import "time"

// PriceData is the data structure for internal price tracking
type PriceData struct {
	Timestamp 	time.Time			`json:"timestamp"`
	Prices		map[string]Price	`json:"prices"`
}

// Price represents our internal price structure
type Price struct {
	Bid		float64	`json:"bid"`
	Ask		float64	`json:"ask"`
	Last	float64	`json:"last"`
}

// OrderbookData represents our internal orderbook structure
type OrderbookData struct {
	Timestamp	time.Time	`json:"timestamp"`
	BuyOrders	[]Order		`json:"buyOrders"`
	SellOrders	[]Order		`json:"sellOrders"`
}

// Order represents our internal order structure
type Order struct {
	Amount		float64	`json:"amount"`
	Rate		float64	`json:"rate"`
	Total		float64	`json:"total"`
}

// TradeData represents our internal trade tracking structure
