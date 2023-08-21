package main

import (
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/vx416/bitfinex-api-go/pkg/models/common"
	bfx "github.com/vx416/bitfinex-api-go/v2"
	"github.com/vx416/bitfinex-api-go/v2/rest"
)

func main() {
	c := rest.NewClient()

	last(c)
	history(c)
	historyWithQuery(c)
}

func last(c *rest.Client) {
	candle, err := c.Candles.Last(common.TradingPrefix+"BTCUSD", common.FiveMinutes)
	if err != nil {
		log.Fatalf("last: %s", err)
	}

	spew.Dump(candle)
}

func history(c *rest.Client) {
	candles, err := c.Candles.History(common.TradingPrefix+"BTCUSD", common.FiveMinutes)
	if err != nil {
		log.Fatalf("history: %s", err)
	}

	spew.Dump(candles)
}

func historyWithQuery(c *rest.Client) {
	now := time.Now()
	millis := now.UnixNano() / 1000000
	prior := now.Add(time.Duration(-24) * time.Hour)
	millisStart := prior.UnixNano() / 1000000
	start := common.Mts(millisStart)
	end := common.Mts(millis)

	candles, err := c.Candles.HistoryWithQuery(
		common.TradingPrefix+bfx.BTCUSD,
		common.FiveMinutes,
		start,
		end,
		200,
		common.OldestFirst,
	)

	if err != nil {
		log.Fatalf("historyWithQuery: %s", err)
	}

	spew.Dump(candles)
}
