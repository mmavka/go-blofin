// Example for testing GetFundingRateHistory from the Blofin public API.
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
	params.Set("instId", "BTC-USDT") // required
	// params.Set("limit", "10")        // Optional: number of records (default 100) max 100
	// params.Set("before", "<timestamp>") // Optional
	// params.Set("after", "<timestamp>") // Optional

	rates, err := client.GetFundingRateHistory(context.Background(), params)
	if err != nil {
		log.Fatalf("failed to get funding rate history: %v", err)
	}

	for _, r := range rates {
		fmt.Printf("%+v\n", r)
	}
}
