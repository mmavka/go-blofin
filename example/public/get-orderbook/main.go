// Example for testing GetOrderBook from the Blofin public API.
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/mmavka/go-blofin/rest"
)

func main() {
	client := rest.NewClient() // Uses BaseURLProd by default

	params := url.Values{}
	params.Set("instId", "BTC-USDT") // required
	params.Set("size", "5")          // Optional: Order book depth per side. Maximum 100, e.g. 100 bids + 100 asks. Default returns to 1 depth data

	orderbook, err := client.GetOrderBook(context.Background(), params)
	if err != nil {
		log.Fatalf("failed to get order book: %v", err)
	}
	if orderbook == nil {
		fmt.Println("no order book data returned")
		return
	}

	fmt.Println("Asks:")
	for _, ask := range orderbook.Asks {
		fmt.Printf("%+v\n", ask)
	}
	fmt.Println("Bids:")
	for _, bid := range orderbook.Bids {
		fmt.Printf("%+v\n", bid)
	}
	fmt.Printf("Timestamp: %s\n", orderbook.Ts)
}
