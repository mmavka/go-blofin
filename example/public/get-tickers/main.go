// Example for testing GetTickers from the Blofin public API.
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
	params.Set("instType", "SPOT") // Optional: SPOT, SWAP, FUTURES, OPTION

	tickers, err := client.GetTickers(context.Background(), params)
	if err != nil {
		slog.Error("failed to get tickers", "error", err)
		os.Exit(1)
	}

	for _, t := range tickers {
		fmt.Printf("%+v\n", t)
	}
}
