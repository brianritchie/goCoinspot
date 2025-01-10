# goCoinSpot - CoinSpot Market Data Collector

Objective: A modular Go application for collecting and storing market data from CoinSpot's public API. This application collects price and order book data at configurable intervals and stores it in a structured format for analysis.

## Features

- Configurable data collection intervals
- Modular collector design for easy extension
- Structured file storage with timestamp-based organization
- Graceful shutdown handling
- Environment-based configuration
- Concurrent data collection
- Error handling and logging

## Disclaimer

See [DISCLAIMER.md](DISCLAIMER.md) for more information.

## Project Structure

```
coinspot-tracker/
├── cmd/
│   └── tracker/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── models/
│   │   ├── api.go            # API models
│   │   ├── models.go         # Internal models
│   │   └── transformers.go   # Data transformation functions
│   ├── collector/
│   │   ├── collector.go      # Core collector logic
│   │   ├── orders.go         # Order book collection
│   │   └── prices.go         # Price data collection
│   └── storage/
│       └── filesystem.go     # File system storage implementation
├── .env.example              # Environment variable template
├── go.mod                    # Go module definition
└── README.md                 # Project documentation
```

## Configuration

Set the following environment variables in your `.env` file:

- `COINSPOT_API_KEY`: Your CoinSpot API key
- `COINSPOT_BASE_URL`: CoinSpot API base URL (default: https://www.coinspot.com.au/pubapi/v2)
- `TRACKED_COINS`: Comma-separated list of coins to track (e.g., "BTC,ETH,XRP")
- `COLLECTION_INTERVAL`: Data collection interval (e.g., "5m" for 5 minutes)
- `OUTPUT_DIR`: Directory for storing collected data

## Data Storage Structure

Data is stored in a hierarchical structure:
```
OUTPUT_DIR/
├── YYYY/
│   ├── MM/
│   │   ├── DD/
│   │   │   ├── prices/
│   │   │   │   └── prices_YYYY-MM-DDTHH-MM-SS.json
│   │   │   └── orders/
│   │   │       └── orders_YYYY-MM-DDTHH-MM-SS.json
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.