package models

import "time"

type ApiResponse struct {
	Status		string	`json:"status"`
	Message		string	`json:"message"`
	BuyOrders	[]Order	`json:"buyorders"`
	SellOrders	[]Order	`json:"sellorders"`
}

type OrderbookSnapshot struct {
	Timestamp	time.Time		`json:"timestamp"`
	Orderbooks	[]OrderbookData	`json:"orderbooks"`
}
