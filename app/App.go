package app

import (
	"context"
	"database/sql"
	"github.com/Bartosz-D3V/binance-exercise/app/binance"
	"github.com/Bartosz-D3V/binance-exercise/app/config"
	"github.com/Bartosz-D3V/binance-exercise/app/trading"
	"github.com/Bartosz-D3V/binance-exercise/app/transaction"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

type App struct {
	Cfg           *config.Config
	Db            *sql.DB
	Conn          *websocket.Conn
	BinanceClient binance.Client
	Repository    transaction.Repository
	TradingSvc    trading.Service
}

func (app App) Start(ctx context.Context, cancelCtx context.CancelFunc) {
	tickerResponses := make(chan binance.BookTickerResponse)
	out := make(chan []transaction.LogEntry)
	interrupt := make(chan os.Signal)

	go app.BinanceClient.ReceiveHandler(ctx, tickerResponses)
	go app.TradingSvc.ProcessStockTick(ctx, tickerResponses, out)

	for {
		select {
		case logEntries := <-out:
			app.printResults(logEntries)
			cancelCtx()
			return
		case <-interrupt:
			cancelCtx()
			return
		}
	}
}

func (app App) printResults(entries []transaction.LogEntry) {
	log.Println("Successfully sold all coins. Please find transaction details below.")
	for _, logEntry := range entries {
		log.Printf("OrderId: %d, Price: %s, Quantity: %s, Timestamp: %s\n", logEntry.OrderId, logEntry.Price, logEntry.Quantity, logEntry.Timestamp)
	}
}
