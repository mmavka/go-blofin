// Example for testing GetMarkPrice from the Blofin public API.
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
	params.Set("instType", "SWAP") // Required: SWAP, FUTURES, OPTION

	markPrices, err := client.GetMarkPrice(context.Background(), params)
	if err != nil {
		slog.Error("failed to get mark price", "error", err)
		os.Exit(1)
	}

	for _, mp := range markPrices {
		fmt.Printf("%+v\n", mp)
	}
}
