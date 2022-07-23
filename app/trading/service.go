package trading

import (
	"context"
	"github.com/Bartosz-D3V/binance-exercise/app/binance"
	"github.com/Bartosz-D3V/binance-exercise/app/calculation"
	"github.com/Bartosz-D3V/binance-exercise/app/transaction"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

type Service interface {
	ProcessStockTick(ctx context.Context, stream <-chan binance.BookTickerResponse, out chan<- []transaction.LogEntry)
}

func New(calcSvc calculation.Service, repository transaction.Repository) Service {
	return &service{
		calcSvc:    calcSvc,
		repository: repository,
	}
}

type service struct {
	calcSvc    calculation.Service
	repository transaction.Repository
}

func (s service) ProcessStockTick(ctx context.Context, stream <-chan binance.BookTickerResponse, out chan<- []transaction.LogEntry) {
	for tick := range stream {
		logEntries := s.repository.GetAll(ctx)

		sellAllowanceLeft := s.calcSvc.GetQuantityToSellLeft(logEntries)
		if sellAllowanceLeft.IsZero() {
			log.Println("All resources have been sold. Quitting.")
			out <- logEntries
			return
		}

		bestBidQuantity, err := decimal.NewFromString(tick.BestBidQuantity)
		if err != nil {
			log.Fatalf("Failed to convert bestBidQuantity value=%s to decimal.Decimal{}. Error=%s", tick.BestBidQuantity, err.Error())
		}

		bestBidPrice, err := decimal.NewFromString(tick.BestBidPrice)
		if err != nil {
			log.Fatalf("Failed to convert bestBidPrice value=%s to decimal.Decimal{}. Error=%s", tick.BestBidPrice, err.Error())
		}

		possibleQuantityToSell := s.calcSvc.GetQuantityToSellPerTick(bestBidQuantity, sellAllowanceLeft)

		if s.calcSvc.TickMeetsTransactionCriteria(bestBidPrice) {
			log.Printf("Binance tick with OrderBookUpdateId=%d meets trading criteria. Selling.\n", tick.OrderBookUpdateId)
			logEntry := transaction.LogEntry{
				OrderId:   tick.OrderBookUpdateId,
				Price:     tick.BestBidPrice,
				Quantity:  possibleQuantityToSell.String(),
				Timestamp: time.Now().Format(time.RFC3339Nano),
			}
			s.repository.Save(ctx, logEntry)
		} else {
			log.Printf("Binance tick with OrderBookUpdateId=%d does not meet trading criteria. Skipping.\n", tick.OrderBookUpdateId)
		}
	}
}
