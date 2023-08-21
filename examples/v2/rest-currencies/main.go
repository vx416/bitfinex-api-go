package main

import (
	"fmt"
	"log"

	"github.com/vx416/bitfinex-api-go/v2/rest"
)

func main() {
	c := rest.NewClient()

	cc, err := c.Currencies.Conf(true, true, true, true, true)
	if err != nil {
		log.Fatalf("getting currency config: %s", err)
	}

	for _, config := range cc {
		fmt.Println(config)
	}
}
