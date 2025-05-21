// Example for testing GetFundingRate from the Blofin public API.
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

	rates, err := client.GetFundingRate(context.Background(), params)
	if err != nil {
		log.Fatalf("failed to get funding rate: %v", err)
	}

	for _, r := range rates {
		fmt.Printf("%+v\n", r)
	}
}
