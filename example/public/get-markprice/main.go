// Example for testing GetMarkPrice from the Blofin public API.
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/mmavka/go-blofin/internal/rest"
)

func main() {
	client := rest.NewClient() // Uses BaseURLProd by default

	params := url.Values{}
	// params.Set("instId", "BTC-USDT") // Optional: filter by instrument

	prices, err := client.GetMarkPrice(context.Background(), params)
	if err != nil {
		log.Fatalf("failed to get mark price: %v", err)
	}

	for _, p := range prices {
		fmt.Printf("%+v\n", p)
	}
}
