package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gocoinspot/internal/config"
	"gocoinspot/internal/storage"
)

type Collector struct {
	config *config.Config
	storage storage.Storage
}

func NewCollector(cfg *config.Config, s storage.Storage) *Collector {
	return &Collector{
		config: cfg, 
		storage: s,
	}
}

func (c *Collector) Start(ctx context.Context) error {
	ticker := time.NewTicker(c.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := c.collect(); err != nil {
				fmt.Printf("Error collecting data: %v\n", err)
			}
		}
	}
}

func (c *Collector) collect() error {
	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := c.collectPrices(); err != nil {
			errCh <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := c.collectOrders(); err != nil {
			errCh <- err
		}
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}