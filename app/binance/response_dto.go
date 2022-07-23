package binance

import "fmt"

type BookTickerResponse struct {
	OrderBookUpdateId int    `json:"u"`
	Symbol            string `json:"s"`
	BestBidPrice      string `json:"b"`
	BestBidQuantity   string `json:"B"`
	BestAskPrice      string `json:"a"`
	BestAskQuantity   string `json:"A"`
}

func (r BookTickerResponse) String() string {
	return fmt.Sprintf(
		"BookTickerResponse(Symbol: %s, BestBidPrice: %s, BestBidQuantity: %s)",
		r.Symbol, r.BestBidPrice, r.BestBidQuantity,
	)
}
