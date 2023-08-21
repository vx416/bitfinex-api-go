package tests

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/vx416/bitfinex-api-go/pkg/models/book"
	"github.com/vx416/bitfinex-api-go/pkg/models/candle"
	"github.com/vx416/bitfinex-api-go/pkg/models/common"
	"github.com/vx416/bitfinex-api-go/pkg/models/ticker"
	"github.com/vx416/bitfinex-api-go/pkg/models/trade"
	"github.com/vx416/bitfinex-api-go/pkg/models/tradeexecution"
	"github.com/vx416/bitfinex-api-go/pkg/models/tradeexecutionupdate"
	"github.com/vx416/bitfinex-api-go/v2"
	"github.com/vx416/bitfinex-api-go/v2/websocket"
)

// wait2 will wait for at least "count" messages on channel "ch" within time "t", or return an error
func wait2(ch <-chan interface{}, count int, bc <-chan error, t time.Duration) error {
	c := make(chan interface{})
	go func() {
		<-ch
		close(c)
	}()
	select {
	case <-bc:
		return fmt.Errorf("transport closed while waiting")
	case <-c:
		return nil // normal
	case <-time.After(t):
		return fmt.Errorf("timed out waiting")
	}
}

func wait(wg *sync.WaitGroup, bc <-chan error, to time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-bc:
		return fmt.Errorf("websocket closed while waiting") // timed out
	case <-c:
		return nil // completed normally
	case <-time.After(to):
		return fmt.Errorf("timed out waiting") // timed out
	}
}

func TestPublicTicker(t *testing.T) {
	c := websocket.New()

	err := c.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Close()

	subs := make(chan interface{}, 10)
	unsubs := make(chan interface{}, 10)
	infos := make(chan interface{}, 10)
	tick := make(chan interface{}, 100)

	errch := make(chan error)
	go func() {
		// nolint:megacheck
		for {
			select {
			case msg := <-c.Listen():
				if msg == nil {
					return
				}
				log.Printf("recv msg: %#v", msg)
				switch m := msg.(type) {
				case error:
					errch <- msg.(error)
				case *websocket.UnsubscribeEvent:
					unsubs <- m
				case *websocket.SubscribeEvent:
					subs <- m
				case *websocket.InfoEvent:
					infos <- m
				case *ticker.Snapshot:
					tick <- m
				case *ticker.Ticker:
					tick <- m
				default:
					t.Logf("test recv: %#v", msg)
				}
			}
		}
	}()

	ctx, cxl := context.WithTimeout(context.Background(), time.Second*5)
	defer cxl()
	id, err := c.SubscribeTicker(ctx, common.TradingPrefix+bitfinex.BTCUSD)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait2(tick, 1, errch, 2*time.Second); err != nil {
		t.Fatalf("failed to receive ticker message from websocket: %s", err)
	}

	err = c.Unsubscribe(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait2(unsubs, 1, errch, 2*time.Second); err != nil {
		t.Errorf("failed to receive unsubscribe message from websocket: %s", err)
	}
}

func TestPublicTrades(t *testing.T) {
	c := websocket.New()
	wg := sync.WaitGroup{}
	wg.Add(3) // 1. Info with version, 2. Subscription event, 3. 3 x data message

	err := c.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Close()

	subs := make(chan interface{}, 10)
	unsubs := make(chan interface{}, 10)
	infos := make(chan interface{}, 10)
	trades := make(chan interface{}, 100)

	errch := make(chan error)
	go func() {
		// nolint:megacheck
		for {
			select {
			case msg := <-c.Listen():
				if msg == nil {
					return
				}
				log.Printf("recv msg: %#v", msg)
				switch m := msg.(type) {
				case error:
					errch <- msg.(error)
				case *websocket.UnsubscribeEvent:
					unsubs <- m
				case *websocket.SubscribeEvent:
					subs <- m
				case *websocket.InfoEvent:
					infos <- m
				case *tradeexecutionupdate.Snapshot:
					trades <- m
				case *trade.Trade:
					trades <- m
				case *tradeexecutionupdate.TradeExecutionUpdate:
					trades <- m
				case *tradeexecution.TradeExecution:
					trades <- m
				case *trade.Snapshot:
					trades <- m
				default:
					t.Logf("test recv: %#v", msg)
				}
			}
		}
	}()

	ctx, cxl := context.WithTimeout(context.Background(), time.Second*5)
	defer cxl()
	id, err := c.SubscribeTrades(ctx, common.TradingPrefix+bitfinex.BTCUSD)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait2(trades, 1, errch, 2*time.Second); err != nil {
		t.Errorf("failed to receive trade message from websocket: %s", err)
	}

	err = c.Unsubscribe(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait2(unsubs, 1, errch, 2*time.Second); err != nil {
		t.Errorf("failed to receive unsubscribe message from websocket: %s", err)
	}
}

func TestPublicBooks(t *testing.T) {
	c := websocket.New()
	wg := sync.WaitGroup{}
	wg.Add(3) // 1. Info with version, 2. Subscription event, 3. data message

	err := c.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Close()

	subs := make(chan interface{}, 10)
	unsubs := make(chan interface{}, 10)
	infos := make(chan interface{}, 10)
	books := make(chan interface{}, 100)

	errch := make(chan error)
	go func() {
		// nolint:megacheck
		for {
			select {
			case msg := <-c.Listen():
				if msg == nil {
					return
				}
				log.Printf("recv msg: %#v", msg)
				switch m := msg.(type) {
				case error:
					errch <- msg.(error)
				case *websocket.UnsubscribeEvent:
					unsubs <- m
				case *websocket.SubscribeEvent:
					subs <- m
				case *websocket.InfoEvent:
					infos <- m
				case *book.Snapshot:
					books <- m
				case *book.Book:
					books <- m
				default:
					t.Logf("test recv: %#v", msg)
				}
			}
		}
	}()

	ctx, cxl := context.WithTimeout(context.Background(), time.Second*5)
	defer cxl()
	id, err := c.SubscribeBook(ctx, common.TradingPrefix+bitfinex.BTCUSD, common.Precision0, common.FrequencyRealtime, 1)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait2(books, 1, errch, 5*time.Second); err != nil {
		t.Fatalf("failed to receive book update message from websocket: %s", err)
	}

	err = c.Unsubscribe(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait2(unsubs, 1, errch, 5*time.Second); err != nil {
		t.Errorf("failed to receive unsubscribe message from websocket: %s", err)
	}
}

func TestPublicCandles(t *testing.T) {
	c := websocket.New()
	wg := sync.WaitGroup{}
	wg.Add(3) // 1. Info with version, 2. Subscription event, 3. data message

	err := c.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Close()

	subs := make(chan interface{}, 10)
	unsubs := make(chan interface{}, 10)
	infos := make(chan interface{}, 10)
	candles := make(chan interface{}, 100)

	errch := make(chan error)
	go func() {
		// nolint:megacheck
		for {
			select {
			case msg := <-c.Listen():
				if msg == nil {
					return
				}
				log.Printf("recv msg: %#v", msg)
				switch m := msg.(type) {
				case error:
					errch <- msg.(error)
				case *websocket.UnsubscribeEvent:
					unsubs <- m
				case *websocket.SubscribeEvent:
					subs <- m
				case *websocket.InfoEvent:
					infos <- m
				case *candle.Candle:
					candles <- m
				case *candle.Snapshot:
					candles <- m
				default:
					t.Logf("test recv: %#v", msg)
				}
			}
		}
	}()

	ctx, cxl := context.WithTimeout(context.Background(), time.Second*5)
	defer cxl()
	id, err := c.SubscribeCandles(ctx, common.TradingPrefix+bitfinex.BTCUSD, common.OneMonth)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait2(candles, 1, errch, 2*time.Second); err != nil {
		t.Errorf("failed to receive a candle message from websocket: %s", err)
	}

	err = c.Unsubscribe(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait2(unsubs, 1, errch, 2*time.Second); err != nil {
		t.Errorf("failed to receive an unsubscribe message from websocket: %s", err)
	}
}
