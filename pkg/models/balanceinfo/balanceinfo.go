package balanceinfo

import (
	"fmt"

	"github.com/vx416/bitfinex-api-go/pkg/convert"
)

type BalanceInfo struct {
	TotalAUM float64
	NetAUM   float64
}

type Update BalanceInfo

func FromRaw(raw []interface{}) (o *BalanceInfo, err error) {
	if len(raw) < 2 {
		return o, fmt.Errorf("data slice too short for balance info: %#v", raw)
	}

	o = &BalanceInfo{
		TotalAUM: convert.F64ValOrZero(raw[0]),
		NetAUM:   convert.F64ValOrZero(raw[1]),
	}

	return
}

func UpdateFromRaw(raw []interface{}) (Update, error) {
	bi, err := FromRaw(raw)
	if err != nil {
		return Update{}, err
	}

	u := Update(*bi)
	return u, nil
}
