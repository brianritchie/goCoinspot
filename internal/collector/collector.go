package collector

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/brianritchie/goCoinspot/internal/config"
	"github.com/brianritchie/goCoinspot/internal/storage"
)

type Collector struct {
	config *config.Config
	storage storage.Storage
	client *APIClient
}

// APIClient wraps HTTP Client with rate limits and retries
type APIClient struct {
	baseURL string
	apiKey string
	rateLimiter *time.Ticker
	httpClient *http.Client
}

func NewCollector(cfg *config.Config, s storage.Storage) *Collector {
	client := &APIClient{
		baseURL: cfg.BaseURL,
		apiKey: cfg.APIKey,
		rateLimiter: time.NewTicker(time.Minute / time.Duration(cfg.RequestsPerMinute)),
		httpClient: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
	}


	return &Collector{
		config: cfg, 
		storage: s, 
		client: client,
	}
}

func (c *Collector) Start(ctx context.Context) error {
	ticker := time.NewTicker(c.config.Interval)
	defer ticker.Stop()

	//Initial collection on start
	if err := c.collect(ctx); err != nil {
		fmt.Printf("Initial collection error: %v\n", err)
	}

	for {
		select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				if err := c.collect(ctx); err != nil {
					fmt.Printf("Error collecting data: %v\n", err)
				}
		}
	}
}

func (c *Collector) collect(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(c.config.Coins)*2) // Space for price and order related errors per coin

	// Collect prices for all coins
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.collectPrices(ctx); err != nil {
			errCh <- fmt.Errorf("Price collection error: %w", err)
		}
	}()

	// Collect orders for each coin
	for _, coin := range c.config.Coins {
		wg.Add(1)
		go func(coinType string) {
			defer wg.Done()
			if err := c.collectOrders(ctx, coinType); err != nil {
				errCh <- fmt.Errorf("Order collection error for %s: %w", coinType, err)
			}
		}(coin)
	}

	// Wait for all collections to complete
	wg.Wait()
	close(errCh)

	// Collect all errors
	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("Errors during collection: %v", errors)
	}

	return nil
}

