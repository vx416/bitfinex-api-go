package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vx416/bitfinex-api-go/v2/rest"
)

// Set BFX_API_KEY and BFX_API_SECRET:
//
// export BFX_API_KEY=<your-api-key>
// export BFX_API_SECRET=<your-api-secret>
//
// you can obtain it from https://www.bitfinex.com/api

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")

	c := rest.
		NewClient().
		Credentials(key, secret)

	deleted, err := c.Pulse.DeletePulse("437b5b44-0f7d-4638-baff-3bbf6966482d")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted: %d\n", deleted)
}
