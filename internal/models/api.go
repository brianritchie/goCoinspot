package models

import "time"

// Latest Prices endpoint
type APIPriceResponse struct {
	Status		string	`json:"status"`
	Message		string	`json:"message"`
	Prices		map[string]Price	`json:"prices"`
}

type APIPrice struct {
	Bid		float64	`json:"bid"`
	Ask		float64	`json:"ask"`
	Last	float64	`json:"last"`
}

// Open Orders endpoint
type APIOrderbookResponse struct {
	Status		string	`json:"status"`
	Message		string	`json:"message"`
	BuyOrders	[]APIOrder	`json:"buyorders"`
	SellOrders	[]APIOrder	`json:"sellorders"`
}

type APIOrder struct {
	Amount		float64	`json:"amount"`
	Rate		float64	`json:"rate"`
	Total		float64	`json:"total"`
	Coin		string	`json:"coin,omitempty"`
	Market		string	`json:"market,omitempty"`
}

// Completed Orders endpoint
type APICompletedOrdersResponse struct {
	Status		string	`json:"status"`
	Message		string	`json:"message"`
	BuyOrders	[]APIOrderDetail	`json:"buyorders"`
	SellOrders	[]APIOrderDetail	`json:"sellorders"`
}

type APIOrderDetail struct {
	Amount		float64	`json:"amount"`
	Rate		float64	`json:"rate"`
	Total		float64	`json:"total"`
	Coin		string	`json:"coin"`
	Market		string	`json:"market"`
	SoldDate	time.Time	`json:"solddate"`
	AudFeeExGst	float64	`json:"audfeeexgst"`
	AudGst		float64	`json:"audgst"`
	AudTotal	float64	`json:"audtotal"`
}
