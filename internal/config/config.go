package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// ConfigValue represents a configuration value with its parsing logic
type configValue[T any] struct {
	key				string
	defaultValue	string
	parser			func(string) (T, error)
	fallback	T // Used when parsing fails
}

type Config struct {
	APIKey				string
	BaseURL				string
	Coins				[]string
	Interval			time.Duration
	OutputDir			string
	RequestsPerMinute	int			// Rate limit
	RequestTimeout		time.Duration	// HTTP client timeout
	RetryAttempts		int			// Number of retry attempts
	RetryBaseDelay		time.Duration	// Base delay for retry attempts
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	// Define duration configurations
	durationConfigs := []configValue[time.Duration]{
		{
			key: "COLLECTION_INTERVAL",
			defaultValue: "5m",
			parser: time.ParseDuration,
		},
		{
			key: "REQUEST_TIMEOUT",
			defaultValue: "30s",
			parser: time.ParseDuration,
		},
		{
			key: "RETRY_BASE_DELAY",
			defaultValue: "100ms",
			parser: time.ParseDuration,
		},
	}

	// Define integer configurations
	intConfigs := []configValue[int]{
		{
			key: "REQUESTS_PER_MINUTE",
			defaultValue: "30",
			parser: strconv.Atoi,
			fallback: 30,
		},
		{
			key: "RETRY_ATTEMPTS",
			defaultValue: "3",
			parser: strconv.Atoi,
			fallback: 3,
		},
	}

	// Parse duration configurations
	durations := make(map[string]time.Duration)
	for _, cfg := range durationConfigs {
		value, err := parseConfigValue(cfg)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %w", cfg.key, err)
		}
		durations[cfg.key] = value
	}

	// Parse integer configurations
	integers := make(map[string]int)
	for _, cfg := range intConfigs {
		value, err := parseConfigValue(cfg)
		if err != nil {
			integers[cfg.key] = cfg.fallback
		} else {
			integers[cfg.key] = value
		}
	}

	return &Config{
		APIKey:          os.Getenv("COINSPOT_API_KEY"),
        BaseURL:         getEnvWithDefault("COINSPOT_BASE_URL", "https://www.coinspot.com.au/pubapi/v2"),
        Coins:           strings.Split(os.Getenv("TRACKED_COINS"), ","),
        Interval:        durations["COLLECTION_INTERVAL"],
        OutputDir:       os.Getenv("OUTPUT_DIR"),
        RequestsPerMinute:  integers["REQUESTS_PER_MINUTE"],
        RequestTimeout:  durations["REQUEST_TIMEOUT"],
        RetryAttempts:   integers["RETRY_ATTEMPTS"],
        RetryBaseDelay:  durations["RETRY_BASE_DELAY"],
	}, nil
}

func parseConfigValue[T any](cfg configValue[T]) (T, error) {
	value := getEnvWithDefault(cfg.key, cfg.defaultValue)
	return cfg.parser(value)
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}