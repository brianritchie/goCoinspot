package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brianritchie/gocoinspot/internal/models"
)

func (c *Collector) collectOrders(ctx context.Context, coinType string) error {
	// Collect open orders
	openOrders, err := c.client.GetOpenOrders(ctx, coinType)
	if err != nil {
		return fmt.Errorf("failed to fetch open orders: %w", err)
	}

	orderbookData := models.TransformAPIOrderbook(openOrders, coinType)
	if err := c.storage.StoreOrderbookData(orderbookData); err != nil {
		return fmt.Errorf("failed to store orderbook data: %w", err)
	}

	// Collect completed orders
	completedOrders, err := c.client.GetCompletedOrders(ctx, coinType)
	if err != nil {
		return fmt.Errorf("failed to fetch completed orders: %w", err)
	}

	tradeData := models.TransformAPICompletedOrders(completedOrders)
	if err := c.storage.StoreTradeData(tradeData); err != nil {
		return fmt.Errorf("failed to store trade data: %w", err)
	}

	return nil
}

// API Client methods
func (c *APIClient) GetOpenOrders(ctx context.Context, coinType string) (*models.APIOrderbookResponse, error) {
	<-c.rateLimiter.C // Rate limiting

	url := fmt.Sprintf("%s/orders/open/%s", c.baseURL, coinType)
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orderbook models.APIOrderbookResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderbook); err != nil {
		return nil, fmt.Errorf("failed to decode orderbook response: %w", err)
	}

	if err := models.ValidateAPIResponse(&orderbook); err != nil {
		return nil, err
	}

	return &orderbook, nil
}

func (c *APIClient) GetCompletedOrders(ctx context.Context, coinType string) (*models.APICompletedOrdersResponse, error) {
	<-c.rateLimiter.C // Rate limiting

	url := fmt.Sprintf("%s/orders/completed/%s", c.baseURL, coinType)
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var completedOrders models.APICompletedOrdersResponse
	if err := json.NewDecoder(resp.Body).Decode(&completedOrders); err != nil {
		return nil, fmt.Errorf("failed to decode completed orders response: %w", err)
	}

	if err := models.ValidateAPIResponse(&completedOrders); err != nil {
		return nil, err
	}

	return &completedOrders, nil
}

// Helper method for making HTTP requests
func (c *APIClient) makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return resp, nil
}

