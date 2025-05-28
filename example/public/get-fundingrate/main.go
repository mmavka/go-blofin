// Example for testing GetFundingRate from the Blofin public API.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/mmavka/go-blofin/rest"
)

func main() {
	client := rest.NewClient() // Uses BaseURLProd by default

	params := url.Values{}
	params.Set("instId", "BTC-USDT")

	fundingRates, err := client.GetFundingRate(context.Background(), params)
	if err != nil {
		slog.Error("failed to get funding rate", "error", err)
		os.Exit(1)
	}

	for _, fr := range fundingRates {
		fmt.Printf("%+v\n", fr)
	}
}
