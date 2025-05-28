// Example for testing GetCandlesticks from the Blofin public API.
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
	params.Set("bar", "1m")    // Required: bar size
	params.Set("limit", "100") // Optional: number of candlesticks. The maximum is 300. The default is 100

	candlesticks, err := client.GetCandlesticks(context.Background(), params)
	if err != nil {
		slog.Error("failed to get candlesticks", "error", err)
		os.Exit(1)
	}

	for _, c := range candlesticks {
		fmt.Printf("%+v\n", c)
	}
}
