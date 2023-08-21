package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/vx416/bitfinex-api-go/pkg/models/order"
	"github.com/vx416/bitfinex-api-go/v2/rest"
)

// Set BFX_API_KEY and BFX_API_SECRET:
//
// export BFX_API_KEY=<your-api-key>
// export BFX_API_SECRET=<your-api-secret>
//
// you can obtain it from https://www.bitfinex.com/api
//
// Below you can see different variations of using Order Multi ops

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")

	c := rest.
		NewClient().
		Credentials(key, secret)

	cancelOrdersMultiOp(c)
	cancelOrderMultiOp(c)
	orderNewMultiOp(c)
	orderUpdateMultiOp(c)
	orderMultiOp(c)
}

func cancelOrdersMultiOp(c *rest.Client) {
	resp, err := c.Orders.CancelOrdersMultiOp(rest.OrderIDs{1189452506, 1189452507})
	if err != nil {
		log.Fatalf("CancelOrdersMultiOp error: %s", err)
	}

	spew.Dump(resp)
}

func cancelOrderMultiOp(c *rest.Client) {
	resp, err := c.Orders.CancelOrderMultiOp(1189502586)
	if err != nil {
		log.Fatalf("CancelOrderMultiOp error: %s", err)
	}

	spew.Dump(resp)
}

func orderNewMultiOp(c *rest.Client) {
	o := order.NewRequest{
		CID:    119,
		GID:    118,
		Type:   "EXCHANGE LIMIT",
		Symbol: "tBTCUSD",
		Price:  12,
		Amount: 0.002,
	}

	resp, err := c.Orders.OrderNewMultiOp(o)
	if err != nil {
		log.Fatalf("OrderNewMultiOp error: %s", err)
	}

	spew.Dump(resp)
}

func orderUpdateMultiOp(c *rest.Client) {
	o := order.UpdateRequest{
		ID:     1189503586,
		Price:  12,
		Amount: 0.002,
	}

	resp, err := c.Orders.OrderUpdateMultiOp(o)
	if err != nil {
		log.Fatalf("OrderUpdateMultiOp error: %s", err)
	}

	spew.Dump(resp)
}

func orderMultiOp(c *rest.Client) {
	pld := rest.OrderOps{
		{
			"on",
			order.NewRequest{
				CID:    987,
				GID:    876,
				Type:   "EXCHANGE LIMIT",
				Symbol: "tBTCUSD",
				Price:  13,
				Amount: 0.001,
			},
		},
		{
			"oc",
			map[string]int{"id": 1189502587},
		},
		{
			"oc_multi",
			map[string][]int{"id": rest.OrderIDs{1189502588, 1189503341}},
		},
		{
			"ou",
			order.UpdateRequest{
				ID:     1189503342,
				Price:  15,
				Amount: 0.002,
			},
		},
	}

	resp, err := c.Orders.OrderMultiOp(pld)
	if err != nil {
		log.Fatalf("OrderMultiOp error: %s", err)
	}

	spew.Dump(resp)
}
