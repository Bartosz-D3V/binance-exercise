package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Bartosz-D3V/binance-exercise/app"
	"github.com/Bartosz-D3V/binance-exercise/app/binance"
	"github.com/Bartosz-D3V/binance-exercise/app/calculation"
	"github.com/Bartosz-D3V/binance-exercise/app/config"
	"github.com/Bartosz-D3V/binance-exercise/app/trading"
	"github.com/Bartosz-D3V/binance-exercise/app/transaction"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shopspring/decimal"
	"log"
	"net/url"
	"strings"
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config file with system environment variables: %v", err)
	}

	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		log.Fatalf("failed to open sqllite DB file. Error=%s", err.Error())
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}

	socketUrl := url.URL{Scheme: "wss", Host: "stream.binancefuture.com", Path: fmt.Sprintf("ws/%s@bookTicker", strings.ToLower(cfg.Symbol))}
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, socketUrl.String(), nil)
	if err != nil {
		log.Fatalf("failed to connect to WSS. Error=%s", err.Error())
	}

	quantityToSell, err := decimal.NewFromString(cfg.QuantityToSell)
	if err != nil {
		log.Fatalf("failed to convert quantityToSell to decimal.decimal{}. Error=%s", err.Error())
	}
	minimumBid, err := decimal.NewFromString(cfg.MinimumBid)
	if err != nil {
		log.Fatalf("failed to convert minimumBid to decimal.decimal{}: Error=%s", err.Error())
	}

	binanceClient := binance.New(conn)
	transactionRepo := transaction.New(db)
	calcSvc := calculation.New(quantityToSell, minimumBid)
	tradingSvc := trading.New(calcSvc, transactionRepo)

	application := app.App{
		Cfg:           cfg,
		Db:            db,
		Conn:          conn,
		BinanceClient: binanceClient,
		Repository:    transactionRepo,
		TradingSvc:    tradingSvc,
	}

	application.Start(ctx, cancelCtx)
}
