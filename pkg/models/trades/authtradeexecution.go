package trades

import (
	"fmt"

	"github.com/vx416/bitfinex-api-go/pkg/convert"
)

// AuthTradeExecution used for mapping authenticated trade execution raw messages
type AuthTradeExecution struct {
	ID            int64
	Pair          string
	MTS           int64
	OrderID       int64
	ExecAmount    float64
	ExecPrice     float64
	OrderType     string
	OrderPrice    float64
	Maker         int
	ClientOrderID int64
}

// ATEFromRaw - authenticated trade execution mapping to data type
func ATEFromRaw(raw []interface{}) (e AuthTradeExecution, err error) {
	if len(raw) < 12 {
		return AuthTradeExecution{}, fmt.Errorf("data slice too short for auth trade execution: %#v", raw)
	}

	e = AuthTradeExecution{
		ID:            convert.I64ValOrZero(raw[0]),
		Pair:          convert.SValOrEmpty(raw[1]),
		MTS:           convert.I64ValOrZero(raw[2]),
		OrderID:       convert.I64ValOrZero(raw[3]),
		ExecAmount:    convert.F64ValOrZero(raw[4]),
		ExecPrice:     convert.F64ValOrZero(raw[5]),
		OrderType:     convert.SValOrEmpty(raw[6]),
		OrderPrice:    convert.F64ValOrZero(raw[7]),
		Maker:         convert.ToInt(raw[8]),
		ClientOrderID: convert.I64ValOrZero(raw[11]),
	}

	return
}
