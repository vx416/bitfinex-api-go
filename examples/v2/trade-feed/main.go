package main

import (
	"context"
	"log"
	_ "net/http/pprof"

	"github.com/vx416/bitfinex-api-go/pkg/models/trade"
	"github.com/vx416/bitfinex-api-go/v2/websocket"
)

func main() {
	client := websocket.New()
	err := client.Connect()
	if err != nil {
		log.Printf("could not connect: %s", err.Error())
		return
	}

	for obj := range client.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %s", obj)
			return
		case *trade.Trade:
			log.Printf("New trade: %+v\n", obj)
		case *websocket.InfoEvent:
			// Info event confirms connection to the bfx websocket
			log.Printf("Subscribing to tBTCUSD")
			_, err := client.SubscribeTrades(context.Background(), "tBTCUSD")
			if err != nil {
				log.Printf("could not subscribe to trades: %s", err.Error())
			}
		default:
			log.Printf("MSG RECV: %#v", obj)
		}
	}
}
