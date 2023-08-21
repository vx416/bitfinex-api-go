package status_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vx416/bitfinex-api-go/pkg/models/status"
)

func TestLiqFromRaw(t *testing.T) {
	testCases := map[string]struct {
		data    []interface{}
		err     func(*testing.T, error)
		success func(*testing.T, interface{})
	}{
		"empty pld": {
			data: []interface{}{},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"invalid pld": {
			data: []interface{}{1},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"valid pld": {
			data: []interface{}{
				"pos", 145400868, 1609144352338, nil, "tETHF0:USTF0",
				-1.67288094, 730.96, nil, 1, 1, nil, 736.13,
			},
			err: func(t *testing.T, got error) {
				assert.Nil(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Equal(t, got, &status.Liquidation{
					Symbol:        "tETHF0:USTF0",
					PositionID:    145400868,
					MTS:           1609144352338,
					Amount:        -1.67288094,
					BasePrice:     730.96,
					IsMatch:       1,
					IsMarketSold:  1,
					PriceAcquired: 736.13,
				})
			},
		},
	}

	for k, v := range testCases {
		t.Run(k, func(t *testing.T) {
			got, err := status.LiqFromRaw(v.data)
			v.err(t, err)
			v.success(t, got)
		})
	}
}

func TestLiqSnapshotFromRaw(t *testing.T) {
	testCases := map[string]struct {
		data    [][]interface{}
		err     func(*testing.T, error)
		success func(*testing.T, interface{})
	}{
		"empty slice": {
			data: [][]interface{}{},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"invalid pld": {
			data: [][]interface{}{{1}},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"valid pld": {
			data: [][]interface{}{
				{
					"pos", 145400868, 1609144352338, nil, "tETHF0:USTF0",
					-1.67288094, 730.96, nil, 1, 1, nil, 736.13,
				},
				{
					"pos", 145400869, 1609144352338, nil, "tETHF0:USTF0",
					-1.67288094, 730.96, nil, 1, 1, nil, 736.13, nil, 1000,
				},
			},
			err: func(t *testing.T, got error) {
				assert.Nil(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Equal(t, got, &status.LiquidationsSnapshot{
					Snapshot: []*status.Liquidation{
						{
							Symbol:        "tETHF0:USTF0",
							PositionID:    145400868,
							MTS:           1609144352338,
							Amount:        -1.67288094,
							BasePrice:     730.96,
							IsMatch:       1,
							IsMarketSold:  1,
							PriceAcquired: 736.13,
						},
						{
							Symbol:        "tETHF0:USTF0",
							PositionID:    145400869,
							MTS:           1609144352338,
							Amount:        -1.67288094,
							BasePrice:     730.96,
							IsMatch:       1,
							IsMarketSold:  1,
							PriceAcquired: 736.13,
						},
					},
				})
			},
		},
	}

	for k, v := range testCases {
		t.Run(k, func(t *testing.T) {
			got, err := status.LiqSnapshotFromRaw(v.data)
			v.err(t, err)
			v.success(t, got)
		})
	}
}
