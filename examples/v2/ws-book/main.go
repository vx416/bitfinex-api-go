package main

import (
	"context"
	"log"
	"time"

	"github.com/vx416/bitfinex-api-go/pkg/models/common"
	"github.com/vx416/bitfinex-api-go/v2"
	"github.com/vx416/bitfinex-api-go/v2/websocket"
)

func main() {
	p := websocket.NewDefaultParameters()
	// Enable orderbook checksum verification
	p.ManageOrderbook = true
	c := websocket.NewWithParams(p)

	err := c.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}

	// subscribe to BTCUSD book
	ctx, cxl2 := context.WithTimeout(context.Background(), time.Second*5)
	defer cxl2()
	_, err = c.SubscribeBook(ctx, common.TradingPrefix+bitfinex.BTCUSD, common.Precision0, common.FrequencyRealtime, 25)
	if err != nil {
		log.Fatal(err)
	}

	// subscribe to BTCUSD trades
	ctx, cxl3 := context.WithTimeout(context.Background(), time.Second*5)
	defer cxl3()
	_, err = c.SubscribeTrades(ctx, common.TradingPrefix+bitfinex.BTCUSD)
	if err != nil {
		log.Fatal(err)
	}

	for obj := range c.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %+v\n", obj)
		default:
		}
		log.Printf("MSG RECV: %+v\n", obj)

		// Load the latest orderbook
		ob, _ := c.GetOrderbook(common.TradingPrefix + bitfinex.BTCUSD)
		if ob != nil {
			log.Printf("Orderbook asks: %v", ob.Asks())
			log.Printf("Orderbook bids: %v", ob.Bids())
		}
	}
}
