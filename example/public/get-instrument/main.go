// Example for testing GetInstruments from the Blofin public API.
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/mmavka/go-blofin/rest"
)

func main() {
	client := rest.NewClient() // Uses BaseURLProd by default

	params := url.Values{}
	// Example: params.Set("instId", "BTC-USDT") // Optional: filter by instrument

	instruments, err := client.GetInstruments(context.Background(), params)
	if err != nil {
		log.Fatalf("failed to get instruments: %v", err)
	}

	for _, inst := range instruments {
		fmt.Printf("%+v\n", inst)
	}
}
