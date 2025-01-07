package models

import "time"

type PriceData struct {
	Timestamp 	time.Time			`json:"timestamp"`
	Prices		map[string]Price	`json:"prices"`
}

type Price struct {
	Bid		float64	`json:"bid"`
	Ask		float64	`json:"ask"`
	Last	float64	`json:"last"`
}

type OrderbookData struct {
	Timestamp	time.Time	`json:"timestamp"`
	BuyOrders	[]Order		`json:"buyOrders"`
	SellOrders	[]Order		`json:"sellOrders"`
}

type Order struct {
	Amount		float64	`json:"amount"`
	Rate		float64	`json:"rate"`
	Total		float64	`json:"total"`
}