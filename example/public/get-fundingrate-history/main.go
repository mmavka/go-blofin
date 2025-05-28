// Example for testing GetFundingRateHistory from the Blofin public API.
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
	params.Set("limit", "10") // Optional: number of records. The maximum is 100. The default is 100

	history, err := client.GetFundingRateHistory(context.Background(), params)
	if err != nil {
		slog.Error("failed to get funding rate history", "error", err)
		os.Exit(1)
	}

	for _, h := range history {
		fmt.Printf("%+v\n", h)
	}
}
