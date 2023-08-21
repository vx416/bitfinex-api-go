package fundingloan_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vx416/bitfinex-api-go/pkg/models/fundingloan"
)

func TestCancelRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		flcr := fundingloan.CancelRequest{ID: 123}
		got, err := flcr.MarshalJSON()

		require.Nil(t, err)

		expected := "[0, \"flc\", null, {\"id\":123}]"
		assert.Equal(t, expected, string(got))
	})
}
