package models

import "time"

type PriceData struct {
	Timestamp 	time.Time			`json:"timestamp"`
	Prices		map[string]Price	`json:"prices"`
}