package msg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"unicode"

	"github.com/vx416/bitfinex-api-go/pkg/convert"
	"github.com/vx416/bitfinex-api-go/pkg/models/balanceinfo"
	"github.com/vx416/bitfinex-api-go/pkg/models/book"
	"github.com/vx416/bitfinex-api-go/pkg/models/candle"
	"github.com/vx416/bitfinex-api-go/pkg/models/event"
	"github.com/vx416/bitfinex-api-go/pkg/models/fundingcredit"
	"github.com/vx416/bitfinex-api-go/pkg/models/fundingloan"
	"github.com/vx416/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/vx416/bitfinex-api-go/pkg/models/fundingtrade"
	"github.com/vx416/bitfinex-api-go/pkg/models/margin"
	"github.com/vx416/bitfinex-api-go/pkg/models/notification"
	"github.com/vx416/bitfinex-api-go/pkg/models/order"
	"github.com/vx416/bitfinex-api-go/pkg/models/position"
	"github.com/vx416/bitfinex-api-go/pkg/models/status"
	"github.com/vx416/bitfinex-api-go/pkg/models/ticker"
	"github.com/vx416/bitfinex-api-go/pkg/models/trades"
	"github.com/vx416/bitfinex-api-go/pkg/models/wallet"
)

type Msg struct {
	Data     []byte
	Err      error
	CID      int
	IsPublic bool
}

func (m Msg) IsEvent() bool {
	t := bytes.TrimLeftFunc(m.Data, unicode.IsSpace)
	return bytes.HasPrefix(t, []byte("{"))
}

func (m Msg) IsRaw() bool {
	t := bytes.TrimLeftFunc(m.Data, unicode.IsSpace)
	return bytes.HasPrefix(t, []byte("["))
}

// PreprocessRaw takes raw slice of bytes and splits it into:
// 1. raw payload data - always last element of the slice
// 2. chanID - always 1st element of the slice
// 3. msg type - in 3 element msg slice, type is always at index 1
func (m Msg) PreprocessRaw() (raw []interface{}, pld interface{}, chID int64, msgType string, err error) {
	err = json.Unmarshal(m.Data, &raw)
	pld = raw[len(raw)-1]
	chID = convert.I64ValOrZero(raw[0])
	if len(raw) == 3 {
		msgType = convert.SValOrEmpty(raw[1])
	}
	return
}

func (m Msg) ProcessPublic(raw []interface{}, pld interface{}, chID int64, inf event.Info) (interface{}, error) {
	switch data := pld.(type) {
	case string:
		return event.Info{
			ChanID:    chID,
			Subscribe: event.Subscribe{Event: data},
		}, nil
	case []interface{}:
		switch inf.Channel {
		case "trades":
			return trades.FromWSRaw(inf.Symbol, raw, data)
		case "ticker":
			return ticker.FromWSRaw(inf.Symbol, data)
		case "book":
			return book.FromWSRaw(inf.Symbol, inf.Precision, data)
		case "candles":
			return candle.FromWSRaw(inf.Key, data)
		case "status":
			return status.FromWSRaw(inf.Key, data)
		}
	}

	return raw, nil
}

func (m Msg) ProcessPrivate(raw []interface{}, pld interface{}, chID int64, op string) (interface{}, error) {
	switch data := pld.(type) {
	case string:
		return event.Info{
			ChanID:    chID,
			Subscribe: event.Subscribe{Event: data},
		}, nil
	case []interface{}:
		switch op {
		case "bu":
			return balanceinfo.UpdateFromRaw(data)
		case "ps":
			return position.SnapshotFromRaw(data)
		case "pn":
			return position.NewFromRaw(data)
		case "pu":
			return position.UpdateFromRaw(data)
		case "pc":
			return position.CancelFromRaw(data)
		case "ws":
			return wallet.SnapshotFromRaw(data)
		case "wu":
			return wallet.UpdateFromRaw(data)
		case "os":
			return order.SnapshotFromRaw(data)
		case "on":
			return order.NewFromRaw(data)
		case "ou":
			return order.UpdateFromRaw(data)
		case "oc":
			return order.CancelFromRaw(data)
		case "te":
			return trades.ATEFromRaw(data)
		case "tu":
			return trades.ATEUFromRaw(data)
		case "fte":
			return trades.AFTEFromRaw(data)
		case "ftu":
			return trades.AFTUFromRaw(data)
		case "mis":
			return nil, errors.New("mis msg type no longer supported")
		case "miu":
			return margin.FromRaw(data)
		case "n":
			return notification.FromRaw(data)
		case "fos":
			return fundingoffer.SnapshotFromRaw(data)
		case "fon":
			return fundingoffer.NewFromRaw(data)
		case "fou":
			return fundingoffer.UpdateFromRaw(data)
		case "foc":
			return fundingoffer.CancelFromRaw(data)
		case "fcs":
			return fundingcredit.SnapshotFromRaw(data)
		case "fcn":
			return fundingcredit.NewFromRaw(data)
		case "fcu":
			return fundingcredit.UpdateFromRaw(data)
		case "fcc":
			return fundingcredit.CancelFromRaw(data)
		case "fls":
			return fundingloan.SnapshotFromRaw(data)
		case "fln":
			return fundingloan.NewFromRaw(data)
		case "flu":
			return fundingloan.UpdateFromRaw(data)
		case "flc":
			return fundingloan.CancelFromRaw(data)
		case "hfts":
			return fundingtrade.HistoricalSnapshotFromRaw(data)
		case "uac":
			return nil, errors.New("uac msg type no longer supported")
		}
	}

	return raw, nil
}

func (m Msg) ProcessEvent() (i event.Info, err error) {
	if err = json.Unmarshal(m.Data, &i); err != nil {
		return i, fmt.Errorf("parsing msg: %s, err: %s", m.Data, err)
	}
	return
}
