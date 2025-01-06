package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)


type Config struct {
	APIKey		string
	BaseURL		string
	Coins		[]string
	Interval	time.Duration
	OutputDir	string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	interval, err := time.ParseDuration(os.Getenv("COLLECTION_INTERVAL"))
	if err != nil {
		interval = 5 * time.Minute // default interval
	}

	return &Config{
		APIKey: os.Getenv("COINSPOT_API_KEY"),
		BaseURL: os.Getenv("COINSPOT_BASE_URL"),
		Coins: strings.Split(os.Getenv("TRACKED_COINS"), ","),
		Interval: interval,
		OutputDir: os.Getenv("OUTPUT_DIR"),
	}, nil
}