package candle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vx416/bitfinex-api-go/pkg/models/candle"
	"github.com/vx416/bitfinex-api-go/pkg/models/common"
)

func TestNewCandleFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{1.5948918e+12}

		c, err := candle.FromRaw("tBTCUSD", common.FiveMinutes, payload)
		require.NotNil(t, err)
		require.Nil(t, c)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			1.5948918e+12,
			9100,
			9076.9,
			9100.1,
			9054.1,
			149.87216331,
		}

		w, err := candle.FromRaw("tBTCUSD", common.FiveMinutes, payload)
		require.Nil(t, err)

		expected := &candle.Candle{
			Symbol:     "tBTCUSD",
			Resolution: "5m",
			MTS:        1594891800000,
			Open:       9100,
			Close:      9076.9,
			High:       9100.1,
			Low:        9054.1,
			Volume:     149.87216331,
		}

		assert.Equal(t, expected, w)
	})
}

func TestSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := [][]interface{}{}

		c, err := candle.SnapshotFromRaw("tBTCUSD", common.FiveMinutes, payload)
		require.NotNil(t, err)
		require.Nil(t, c)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := [][]interface{}{
			{
				1.5948918e+12,
				9100,
				9076.9,
				9100.1,
				9054.1,
				149.87216331,
			},
			{
				1.5948918e+12,
				9200,
				9076.9,
				9100.1,
				9054.1,
				149.87216331,
			},
		}

		ss, err := candle.SnapshotFromRaw("tBTCUSD", common.FiveMinutes, payload)
		require.Nil(t, err)

		expected := &candle.Snapshot{
			Snapshot: []*candle.Candle{
				{
					Symbol:     "tBTCUSD",
					Resolution: "5m",
					MTS:        1594891800000,
					Open:       9100,
					Close:      9076.9,
					High:       9100.1,
					Low:        9054.1,
					Volume:     149.87216331,
				},
				{
					Symbol:     "tBTCUSD",
					Resolution: "5m",
					MTS:        1594891800000,
					Open:       9200,
					Close:      9076.9,
					High:       9100.1,
					Low:        9054.1,
					Volume:     149.87216331,
				},
			},
		}

		require.Nil(t, err)
		assert.Equal(t, expected, ss)
	})
}

func TestFromWSRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{1.5948918e+12}
		c, err := candle.FromWSRaw("trade:1m:tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, c)
	})

	t.Run("missing arguments", func(t *testing.T) {
		payload := []interface{}{}
		c, err := candle.FromWSRaw("trade:1m:tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, c)
	})

	t.Run("invalid key", func(t *testing.T) {
		payload := []interface{}{}
		c, err := candle.FromWSRaw("tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, c)
	})

	t.Run("valid candle arguments", func(t *testing.T) {
		payload := []interface{}{
			1.5948918e+12,
			9100,
			9076.9,
			9100.1,
			9054.1,
			149.87216331,
		}

		w, err := candle.FromWSRaw("trade:1m:tBTCUSD", payload)
		require.Nil(t, err)

		expected := &candle.Candle{
			Symbol:     "tBTCUSD",
			Resolution: "1m",
			MTS:        1594891800000,
			Open:       9100,
			Close:      9076.9,
			High:       9100.1,
			Low:        9054.1,
			Volume:     149.87216331,
		}

		assert.Equal(t, expected, w)
	})

	t.Run("valid snapshot arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				1.5948918e+12,
				9100,
				9076.9,
				9100.1,
				9054.1,
				149.87216331,
			},
			[]interface{}{
				1.5948918e+12,
				9200,
				9076.9,
				9100.1,
				9054.1,
				149.87216331,
			},
		}

		ss, err := candle.FromWSRaw("trade:1m:tBTCUSD", payload)
		require.Nil(t, err)

		expected := &candle.Snapshot{
			Snapshot: []*candle.Candle{
				{
					Symbol:     "tBTCUSD",
					Resolution: "1m",
					MTS:        1594891800000,
					Open:       9100,
					Close:      9076.9,
					High:       9100.1,
					Low:        9054.1,
					Volume:     149.87216331,
				},
				{
					Symbol:     "tBTCUSD",
					Resolution: "1m",
					MTS:        1594891800000,
					Open:       9200,
					Close:      9076.9,
					High:       9100.1,
					Low:        9054.1,
					Volume:     149.87216331,
				},
			},
		}

		require.Nil(t, err)
		assert.Equal(t, expected, ss)
	})
}
