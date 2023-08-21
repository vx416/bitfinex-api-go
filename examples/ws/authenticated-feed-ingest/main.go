package main

import (
	"log"

	"github.com/vx416/bitfinex-api-go/pkg/models/balanceinfo"
	"github.com/vx416/bitfinex-api-go/pkg/models/event"
	"github.com/vx416/bitfinex-api-go/pkg/models/fundingcredit"
	"github.com/vx416/bitfinex-api-go/pkg/models/fundingloan"
	"github.com/vx416/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/vx416/bitfinex-api-go/pkg/models/fundingtrade"
	"github.com/vx416/bitfinex-api-go/pkg/models/margin"
	"github.com/vx416/bitfinex-api-go/pkg/models/notification"
	"github.com/vx416/bitfinex-api-go/pkg/models/order"
	"github.com/vx416/bitfinex-api-go/pkg/models/position"
	"github.com/vx416/bitfinex-api-go/pkg/models/trades"
	"github.com/vx416/bitfinex-api-go/pkg/models/wallet"
	"github.com/vx416/bitfinex-api-go/pkg/mux"
)

func main() {
	m := mux.New().
		TransformRaw().
		WithAPIKEY("YOUR_API_KEY").
		WithAPISEC("YOUR_API_SECRET").
		Start()

	crash := make(chan error)

	go func() {
		crash <- m.Listen(func(msg interface{}, err error) {
			if err != nil {
				log.Printf("error received: %s\n", err)
			}

			switch v := msg.(type) {
			case event.Info:
				log.Printf("%T: %+v\n", v, v)
			case order.New:
				log.Printf("%T: %+v\n", v, v)
			case *order.Snapshot:
				log.Printf("%T: %+v\n", v, v)
			case order.Update:
				log.Printf("%T: %+v\n", v, v)
			case order.Cancel:
				log.Printf("%T: %+v\n", v, v)
			case wallet.Update:
				log.Printf("%T: %+v\n", v, v)
			case *wallet.Snapshot:
				log.Printf("%T: %+v\n", v, v)
			case balanceinfo.Update:
				log.Printf("%T: %+v\n", v, v)
			case fundingoffer.New:
				log.Printf("%T: %+v\n", v, v)
			case fundingoffer.Cancel:
				log.Printf("%T: %+v\n", v, v)
			case fundingoffer.Update:
				log.Printf("%T: %+v\n", v, v)
			case *fundingoffer.Snapshot:
				log.Printf("%T: %+v\n", v, v)
			case fundingcredit.New:
				log.Printf("%T: %+v\n", v, v)
			case fundingcredit.Update:
				log.Printf("%T: %+v\n", v, v)
			case fundingcredit.Cancel:
				log.Printf("%T: %+v\n", v, v)
			case trades.AuthFundingTradeUpdate:
				log.Printf("%T: %+v\n", v, v)
			case trades.AuthFundingTradeExecuted:
				log.Printf("%T: %+v\n", v, v)
			case trades.AuthTradeExecution:
				log.Printf("%T: %+v\n", v, v)
			case trades.AuthTradeExecutionUpdate:
				log.Printf("%T: %+v\n", v, v)
			case fundingtrade.HistoricalSnapshot:
				log.Printf("%T: %+v\n", v, v)
			case *fundingcredit.Snapshot:
				log.Printf("%T: %+v\n", v, v)
			case *position.Snapshot:
				log.Printf("%T: %+v\n", v, v)
			case position.New:
				log.Printf("%T: %+v\n", v, v)
			case position.Update:
				log.Printf("%T: %+v\n", v, v)
			case position.Cancel:
				log.Printf("%T: %+v\n", v, v)
			case *fundingloan.Snapshot:
				log.Printf("%T: %+v\n", v, v)
			case fundingloan.New:
				log.Printf("%T: %+v\n", v, v)
			case fundingloan.Update:
				log.Printf("%T: %+v\n", v, v)
			case fundingloan.Cancel:
				log.Printf("%T: %+v\n", v, v)
			case notification.Notification:
				log.Printf("%T: %+v\n", v, v)
			case *margin.InfoUpdate:
				log.Printf("%T: %+v\n", v, v)
			case *margin.InfoBase:
				log.Printf("%T: %+v\n", v, v)
			default:
				log.Printf("raw/unhandled: %T: %+v\n", v, v)
			}
		})
	}()

	log.Fatal(<-crash)
}
