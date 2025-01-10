package collector

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/brianritchie/goCoinspot/internal/models"
	"github.com/brianritchie/gocoinspot/internal/models"
)

// collect Prices handles the collection of all latest prices for all specified coins
func (c *Collector) collectPrices(ctx context.Context) error {
	prices, err := c.client.GetLatestPrices(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch latest prices: %w", err)
	}

	priceData := models.TransformAPIPrices(prices)
	if err := c.storage.StorePriceData(priceData); err != nil {
		return fmt.Errorf("failed to store price data: %w", err)
	}

	return nil
}

// Fetch latest prices for specified coins from API
func (c *APIClient) GetLatestPrices(ctx context.Context) (*models.APIPriceResponse, error) {
	<-c.rateLimiter.C // Rate limiting

	url := fmt.Sprintf("%s/latest", c.baseURL)
	resp,err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to make price request: %w", err)
	}
	defer resp.Body.Close()

	var priceResponse models.APIPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResponse); err != nil {
		return nil, fmt.Errorf("failed to decode price response: %w", err)
	}

	if err := models.ValidateAPIPriceResponse(&priceResponse); err != nil {
		return nil, fmt.Errorf("invalid price response: %w", err)
	}

	// Additional validation for specific coins
	if err := c.validatePriceData(&priceResponse); err != nil {
		return nil, fmt.Errorf("failed to validate specific coins: %w", err)
	}

	return &priceResponse, nil
}

func (c *APIClient) validatePriceData(priceResponse *models.APIPriceResponse) error {
	if prices.Prices == nil {
		return fmt.Errorf("no price data received")
	}

	// Log warning for any zero prices
	for coin, price := range prices.Prices {
		if price.Bid == 0 && price.Ask == 0 && price.Last == 0 {
			fmt.Printf("Warning: Zero price for %s\n", coin)
		}
	}

	return nil
}

func (c *APIClient) GetCoinPrice(ctx context.Context, coin string) (*models.APIPrice, error) {
	<-c.rateLimiter.C // Rate limiting

	url := fmt.Sprintf("%s/latest/%s", c.baseURL, coinType)
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to make single coin price request: %w", err)
	}
	defer resp.Body.Close()

	var priceResp struct {
		Status string `json:"status"`
		Message string `json:"message"`
		Prices models.APIPrice `json:"prices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&priceResp); err != nil {
		return nil, fmt.Errorf("failed to decode single coin price response: %w", err)
	}

	if priceResp.Status != "ok" {
		return nil, fmt.Errorf("API error: %s", priceResp.Message)
	}

	return &priceResp.Prices, nil
}
