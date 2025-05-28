// Example for testing GetOrderBook from the Blofin public API.
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
	params.Set("sz", "5") // Optional: order book depth. The maximum is 400. The default is 1

	orderBook, err := client.GetOrderBook(context.Background(), params)
	if err != nil {
		slog.Error("failed to get order book", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Order Book: %+v\n", orderBook)
}
