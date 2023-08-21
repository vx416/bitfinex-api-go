package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vx416/bitfinex-api-go/pkg/models/common"
	"github.com/vx416/bitfinex-api-go/pkg/models/order"
	"github.com/vx416/bitfinex-api-go/v2"
)

func TestOrdersAll(t *testing.T) {
	httpDo := func(_ *http.Client, req *http.Request) (*http.Response, error) {
		msg := `
				[
					[33961681942,"1227",1337,"tBTCUSD",1573482478000,1573485373000,0.001,0.001,"EXCHANGE LIMIT",null,null,null,"0","CANCELED",null,null,15,0,0,0,null,null,null,0,0,null,null,null,"API>BFX",null,null,null],
					[33961681943,"1227",1337,"tBTCUSD",1573482478000,1573485373000,0.001,0.001,"EXCHANGE LIMIT",null,null,null,"0","CANCELED",null,null,15,0,0,0,null,null,null,0,0,null,null,null,"API>BFX",null,null,null],
					[33961681944,"1227",1337,"tBTCUSD",1573482478000,1573485373000,0.001,0.001,"EXCHANGE LIMIT",null,null,null,"0","CANCELED",null,null,15,0,0,0,null,null,null,0,0,null,null,null,"API>BFX",null,null,null]
				]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	orders, err := NewClientWithHttpDo(httpDo).Orders.All()

	if err != nil {
		t.Error(err)
	}

	if len(orders.Snapshot) != 3 {
		t.Fatalf("expected three orders but got %d", len(orders.Snapshot))
	}
}

func TestOrdersHistory(t *testing.T) {
	httpDo := func(_ *http.Client, req *http.Request) (*http.Response, error) {
		msg := `
				[
					[33961681942,"1227",1337,"tBTCUSD",1573482478000,1573485373000,0.001,0.001,"EXCHANGE LIMIT",null,null,null,"0","CANCELED",null,null,15,0,0,0,null,null,null,0,0,null,null,null,"API>BFX",null,null,null],
					[33961681943,"1227",1337,"tBTCUSD",1573482478000,1573485373000,0.001,0.001,"EXCHANGE LIMIT",null,null,null,"0","CANCELED",null,null,15,0,0,0,null,null,null,0,0,null,null,null,"API>BFX",null,null,null],
					[33961681944,"1227",1337,"tBTCUSD",1573482478000,1573485373000,0.001,0.001,"EXCHANGE LIMIT",null,null,null,"0","CANCELED",null,null,15,0,0,0,null,null,null,0,0,null,null,null,"API>BFX",null,null,null]
				]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	orders, err := NewClientWithHttpDo(httpDo).Orders.GetHistoryBySymbol(common.TradingPrefix + bitfinex.IOTBTC)

	if err != nil {
		t.Error(err)
	}

	if len(orders.Snapshot) != 3 {
		t.Errorf("expected three orders but got %d", len(orders.Snapshot))
	}
}

func TestCancelOrderMulti(t *testing.T) {
	t.Run("calls correct resource with correct payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/order/cancel/multi", r.RequestURI)
			assert.Equal(t, "POST", r.Method)

			gotReqPld := CancelOrderMultiRequest{}
			err := json.NewDecoder(r.Body).Decode(&gotReqPld)
			require.Nil(t, err)

			expectedReqPld := CancelOrderMultiRequest{
				OrderIDs:      OrderIDs{123},
				GroupOrderIDs: GroupOrderIDs{234},
				All:           1,
			}
			assert.Equal(t, expectedReqPld, gotReqPld)

			respMock := []interface{}{1568711312683, nil, nil, nil, nil, nil, nil, nil}
			payload, _ := json.Marshal(respMock)
			_, err = w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pld := CancelOrderMultiRequest{
			OrderIDs:      OrderIDs{123},
			GroupOrderIDs: GroupOrderIDs{234},
			All:           1,
		}

		rsp, err := c.Orders.CancelOrderMulti(pld)
		require.Nil(t, err)
		assert.Equal(t, int64(1568711312683), rsp.MTS)
	})
}

func TestCancelOrdersMultiOp(t *testing.T) {
	t.Run("calls correct resource with correct payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/order/multi", r.RequestURI)
			assert.Equal(t, "POST", r.Method)

			gotReqPld := OrderMultiOpsRequest{}
			err := json.NewDecoder(r.Body).Decode(&gotReqPld)
			require.Nil(t, err)

			expectedReqPld := []interface{}{
				"oc_multi",
				map[string]interface{}{
					"id": []interface{}{
						float64(1189428429),
						float64(1189428430),
					},
				},
			}
			assert.Equal(t, expectedReqPld, gotReqPld.Ops[0])

			respMock := []interface{}{1568711312683, nil, nil, nil, nil, nil, nil, nil}
			payload, _ := json.Marshal(respMock)
			_, err = w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		rsp, err := c.Orders.CancelOrdersMultiOp(OrderIDs{1189428429, 1189428430})
		require.Nil(t, err)
		assert.Equal(t, int64(1568711312683), rsp.MTS)
	})
}

func TestCancelOrderMultiOp(t *testing.T) {
	t.Run("calls correct resource with correct payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/order/multi", r.RequestURI)
			assert.Equal(t, "POST", r.Method)

			gotReqPld := OrderMultiOpsRequest{}
			err := json.NewDecoder(r.Body).Decode(&gotReqPld)
			require.Nil(t, err)

			expectedReqPld := []interface{}{
				"oc",
				map[string]interface{}{"id": float64(1189428429)},
			}
			assert.Equal(t, expectedReqPld, gotReqPld.Ops[0])

			respMock := []interface{}{1568711312683, nil, nil, nil, nil, nil, nil, nil}
			payload, _ := json.Marshal(respMock)
			_, err = w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		rsp, err := c.Orders.CancelOrderMultiOp(1189428429)
		require.Nil(t, err)
		assert.Equal(t, int64(1568711312683), rsp.MTS)
	})
}

func TestOrderNewMultiOp(t *testing.T) {
	t.Run("calls correct resource with correct payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/order/multi", r.RequestURI)
			assert.Equal(t, "POST", r.Method)

			gotReqPld := OrderMultiOpsRequest{}
			err := json.NewDecoder(r.Body).Decode(&gotReqPld)
			require.Nil(t, err)

			expectedReqPld := []interface{}{
				"on",
				map[string]interface{}{
					"amount": "0.002",
					"cid":    float64(119),
					"gid":    float64(118),
					"price":  "12",
					"symbol": "tBTCUSD",
					"type":   "EXCHANGE LIMIT",
					"flags":  float64(512),
					"meta": map[string]interface{}{
						"aff_code": "abc",
					},
				},
			}
			assert.Equal(t, expectedReqPld, gotReqPld.Ops[0])

			respMock := []interface{}{1568711312683, nil, nil, nil, nil, nil, nil, nil}
			payload, _ := json.Marshal(respMock)
			_, err = w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		o := order.NewRequest{
			CID:           119,
			GID:           118,
			Type:          "EXCHANGE LIMIT",
			Symbol:        "tBTCUSD",
			Price:         12,
			Amount:        0.002,
			AffiliateCode: "abc",
			Close:         true,
		}

		rsp, err := c.Orders.OrderNewMultiOp(o)
		require.Nil(t, err)
		assert.Equal(t, int64(1568711312683), rsp.MTS)
	})
}

func TestOrderUpdateMultiOp(t *testing.T) {
	t.Run("calls correct resource with correct payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/order/multi", r.RequestURI)
			assert.Equal(t, "POST", r.Method)

			gotReqPld := OrderMultiOpsRequest{}
			err := json.NewDecoder(r.Body).Decode(&gotReqPld)
			require.Nil(t, err)

			expectedReqPld := []interface{}{
				"ou",
				map[string]interface{}{
					"amount": "0.002",
					"price":  "12",
					"id":     float64(1189503586),
					"flags":  float64(64),
				},
			}
			assert.Equal(t, expectedReqPld, gotReqPld.Ops[0])

			respMock := []interface{}{1568711312683, nil, nil, nil, nil, nil, nil, nil}
			payload, _ := json.Marshal(respMock)
			_, err = w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		o := order.UpdateRequest{
			ID:     1189503586,
			Price:  12,
			Amount: 0.002,
			Hidden: true,
		}

		rsp, err := c.Orders.OrderUpdateMultiOp(o)
		require.Nil(t, err)
		assert.Equal(t, int64(1568711312683), rsp.MTS)
	})
}

func TestOrderMultiOp(t *testing.T) {
	t.Run("calls correct resource with correct payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/order/multi", r.RequestURI)
			assert.Equal(t, "POST", r.Method)

			gotReqPld := OrderMultiOpsRequest{}
			err := json.NewDecoder(r.Body).Decode(&gotReqPld)
			require.Nil(t, err)

			expectedReqPld := map[string][]interface{}{
				"on": {
					"on",
					map[string]interface{}{
						"amount": "0.001",
						"cid":    float64(987),
						"flags":  float64(4096),
						"gid":    float64(876),
						"meta": map[string]interface{}{
							"aff_code": "abc",
						},
						"price":  "13",
						"symbol": "tBTCUSD",
						"type":   "EXCHANGE LIMIT",
					},
				},
				"ou": {
					"ou",
					map[string]interface{}{
						"amount": "0.002",
						"price":  "15",
						"id":     float64(1189503342),
						"flags":  float64(64),
					},
				},
				"oc": {
					"oc",
					map[string]interface{}{"id": float64(1189502430)},
				},
				"oc_multi": {
					"oc_multi",
					map[string]interface{}{
						"id": []interface{}{
							float64(1189502431),
							float64(1189502432),
						},
					},
				},
			}

			for _, v := range gotReqPld.Ops {
				key := v[0].(string)
				assert.Equal(t, expectedReqPld[key][1], v[1])
			}

			respMock := []interface{}{1568711312683, nil, nil, nil, nil, nil, nil, nil}
			payload, _ := json.Marshal(respMock)
			_, err = w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pld := OrderOps{
			{
				"on",
				order.NewRequest{
					CID:           987,
					GID:           876,
					Type:          "EXCHANGE LIMIT",
					Symbol:        "tBTCUSD",
					Price:         13,
					Amount:        0.001,
					PostOnly:      true,
					AffiliateCode: "abc",
				},
			},
			{
				"oc",
				map[string]int{"id": 1189502430},
			},
			{
				"oc_multi",
				map[string][]int{"id": OrderIDs{1189502431, 1189502432}},
			},
			{
				"ou",
				order.UpdateRequest{
					ID:     1189503342,
					Price:  15,
					Amount: 0.002,
					Hidden: true,
				},
			},
		}
		rsp, err := c.Orders.OrderMultiOp(pld)
		require.Nil(t, err)
		assert.Equal(t, int64(1568711312683), rsp.MTS)
	})
}
