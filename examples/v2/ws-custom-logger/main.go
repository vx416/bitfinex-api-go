package main

import (
	"context"
	_ "net/http/pprof"
	"os"

	"github.com/op/go-logging"
	"github.com/vx416/bitfinex-api-go/pkg/models/trade"
	"github.com/vx416/bitfinex-api-go/v2/websocket"
)

func main() {
	// create a new go-logger instance
	var log = logging.MustGetLogger("bfx-websocket")
	// create string formatter
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	// apply to logging instance
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	// create websocket client and pass logger
	p := websocket.NewDefaultParameters()
	p.Logger = log
	client := websocket.NewWithParams(p)
	err := client.Connect()
	if err != nil {
		log.Errorf("could not connect: %s", err.Error())
		return
	}

	for obj := range client.Listen() {
		switch obj.(type) {
		case error:
			log.Errorf("channel closed: %s", obj)
			return
		case *trade.Trade:
			log.Infof("New trade: %+v\n", obj)
		case *websocket.InfoEvent:
			// Info event confirms connection to the bfx websocket
			log.Info("Subscribing to tBTCUSD")
			_, err := client.SubscribeTrades(context.Background(), "tBTCUSD")
			if err != nil {
				log.Infof("could not subscribe to trades: %s", err.Error())
			}
		default:
			log.Infof("MSG RECV: %#v", obj)
		}
	}
}
