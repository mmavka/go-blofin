// Example for testing GetTrades from the Blofin public API.
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
	params.Set("limit", "10") // Optional: number of trades. The maximum is 100. The default is 100

	trades, err := client.GetTrades(context.Background(), params)
	if err != nil {
		slog.Error("failed to get trades", "error", err)
		os.Exit(1)
	}

	for _, t := range trades {
		fmt.Printf("%+v\n", t)
	}
}
