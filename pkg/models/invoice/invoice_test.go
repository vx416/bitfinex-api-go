package invoice_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vx416/bitfinex-api-go/pkg/models/invoice"
)

func TestNewInvoiceFromRaw(t *testing.T) {
	t.Run("insufficient arguments", func(t *testing.T) {
		payload := []interface{}{"invoicehash"}

		invc, err := invoice.NewFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, invc)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			"invoicehash",
			"invoice",
			nil,
			nil,
			"0.00016",
		}

		invc, err := invoice.NewFromRaw(payload)
		require.Nil(t, err)

		expected := &invoice.Invoice{
			InvoiceHash: "invoicehash",
			Invoice:     "invoice",
			Amount:      "0.00016",
		}

		assert.Equal(t, expected, invc)
	})
}
