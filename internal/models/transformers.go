package models

import (
	"fmt"
	"time"
)

// Transformers file handles all the conversion logic between our internal models and the API models

// PriceTransformer converts an API price responses to internal models
func TransformAPIPrice(apr *APIPriceResponse) PriceData {
	priceData := PriceData{
		Timestamp: time.Now().UTC(),
		Prices: make(map[string]Price),
	}

	for coinType, apiPrice := range apr.Prices {
		priceData.Prices[coinType] = Price{
			Bid: apiPrice.Bid,
			Ask: apiPrice.Ask,
			Last: apiPrice.Last,
		}
	}

	return priceData
}

// Helper function to transform API orders to internal orders
func transformOrders(apiOrders []APIOrder) []Order {
	orders := make([]Order, len(apiOrders))

	for i, order := range apiOrders {
		orders[i] = Order{
			Amount: order.Amount,
			Rate: order.Rate,
			Total: order.Total,
		}
	}

	return orders
}

// OrderbookTransformer converts an API orderbook responses to internal models
func TransformAPIOrderbook(aor *APIOrderbookResponse, coinType string) OrderbookData {
	return OrderbookData{
		Timestamp: time.Now().UTC(),
		BuyOrders: transformOrders(aor.BuyOrders),
		SellOrders: transformOrders(aor.SellOrders),
	}
}

// TradeTransformer converts API completed orders responses to internal models
func TransformAPICompletedOrders(acor *APICompletedOrdersResponse) TradeData {
	tradeData := TradeData{
		Timestamp: time.Now().UTC(),
		Trades: make([]Trade, 0, len(acor.BuyOrders) + len(acor.SellOrders)),
	}

	for _, buyOrder := range acor.BuyOrders {
		tradeData.Trades = append(tradeData.Trades, transformAPIOrderToTradeInfo(buy, "buy"))
	}

	for _, sellOrder := range acor.SellOrders {
		tradeData.Trades = append(tradeData.Trades, transformAPIOrderToTradeInfo(sell, "sell"))
	}

	return tradeData
}

// Helper function to transform individual orders
func transformAPIOrderToTradeInfo(order APIOrderDetail, ordertype string) TradeInfo {
	return TradeInfo{
		OrderType: ordertype,
		Amount: order.Amount,
		Rate: order.Rate,
		Total: order.Total,
		Market: order.Market,
		ExecutedAt: order.SoldDate,
		Fees: TradeFees{
			FeeExGst: order.AudFeeExGst,
			Gst: order.AudGst,
			Total: order.AudTotal,
		},
	}
}

// Validation helpers
func ValidateAPIResponse(resp interface{}) error {
	switch r := resp.(type) {
		case *APIPriceResponse:
			if r.Status != "ok" || r.Prices == nil {
				return fmt.Errorf("invalid price response: %v", r)
			}
		case *APIOrderbookResponse:
			if r.Status != "ok" || r.BuyOrders == nil || r.SellOrders == nil {
				return fmt.Errorf("invalid orderbook response: %v", r)
			}
		case *APICompletedOrdersResponse:
			if r.Status != "ok" || r.BuyOrders == nil || r.SellOrders == nil {
				return fmt.Errorf("invalid completed orders response: %v", r)
			}
		default:
			return fmt.Errorf("unsupported response type: %T", r)
	}
	return nil
}