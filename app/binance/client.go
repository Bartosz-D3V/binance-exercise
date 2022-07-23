package binance

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
)

type Client interface {
	ReceiveHandler(ctx context.Context, out chan BookTickerResponse)
}

func New(connection *websocket.Conn) Client {
	return &client{
		connection: connection,
	}
}

type client struct {
	connection *websocket.Conn
}

func (c client) ReceiveHandler(ctx context.Context, out chan BookTickerResponse) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.consumeStream(out)
		}
	}
}

func (c client) consumeStream(out chan BookTickerResponse) {
	var resp = BookTickerResponse{}
	err := c.connection.ReadJSON(&resp)
	if err != nil {
		log.Fatalf("Failed to read JSON. Error=%s", err.Error())
	}

	out <- resp

	log.Printf("Received binance tick: %v\n", resp.String())
}
