package models

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