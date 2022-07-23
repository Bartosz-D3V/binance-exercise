package calculation

import (
	"github.com/Bartosz-D3V/binance-exercise/app/transaction"
	"github.com/shopspring/decimal"
	"log"
)

type Service interface {
	GetQuantityToSellPerTick(bestBidQuantity, sellAllowanceLeft decimal.Decimal) decimal.Decimal
	GetQuantityToSellLeft(logEntries []transaction.LogEntry) decimal.Decimal
	TickMeetsTransactionCriteria(bestBidPrice decimal.Decimal) bool
}

func New(quantityToSell, minimumBid decimal.Decimal) Service {
	return &service{
		quantityToSell: quantityToSell,
		minimumBid:     minimumBid,
	}
}

type service struct {
	quantityToSell decimal.Decimal
	minimumBid     decimal.Decimal
}

func (s service) GetQuantityToSellPerTick(bestBidQuantity, sellAllowanceLeft decimal.Decimal) decimal.Decimal {
	if sellAllowanceLeft.GreaterThanOrEqual(bestBidQuantity) {
		return bestBidQuantity
	}
	return sellAllowanceLeft
}

func (s service) GetQuantityToSellLeft(logEntries []transaction.LogEntry) decimal.Decimal {
	totalQuantitySold := decimal.Zero

	for _, entry := range logEntries {
		quantitySold, err := decimal.NewFromString(entry.Quantity)
		if err != nil {
			log.Fatalf("Failed to convert value %s to decimal.Decimal{}", entry.Quantity)
		}
		totalQuantitySold = totalQuantitySold.Add(quantitySold)
	}

	return s.quantityToSell.Sub(totalQuantitySold)
}

func (s service) TickMeetsTransactionCriteria(bestBidPrice decimal.Decimal) bool {
	return bestBidPrice.GreaterThanOrEqual(s.minimumBid)
}
