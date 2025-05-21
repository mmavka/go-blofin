// Example for testing GetTickers from the Blofin public API.
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

	tickers, err := client.GetTickers(context.Background(), params)
	if err != nil {
		log.Fatalf("failed to get tickers: %v", err)
	}

	for _, t := range tickers {
		fmt.Printf("%+v\n", t)
	}
}
