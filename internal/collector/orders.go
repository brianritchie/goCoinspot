package collector

import (
	"net/http"

	"gocoinspot/internal/config"
	"gocoinspot/internal/storage"
)

type OrderCollector struct {
	config	*config.Config
	storage	storage.Storage
	limiter *RateLimiter
	client	*http.Client
}

func NewOrderCollector(cfg *config.Config, store storage.Storage) *OrderCollector {
	return &OrderCollector{
		config: cfg,
		storage: store,
		limiter:  NewRateLimiter(cfg.RequestsPerMinute),
		client: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
	}
}