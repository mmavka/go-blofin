// Example for testing GetCandlesticks from the Blofin public API.
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
	params.Set("instId", "BTC-USDT")
	params.Set("bar", "1m")  // Optional: 1m/3m/5m/15m/30m/1H/2H/4H/6H/8H/12H/1D/3D/1W/1M. Default 1m
	params.Set("limit", "5") // Optional: number of candles (default 500) max 1440
	//params.Set("after", "1716151200000") // Optional: Pagination of data to return records earlier than the requested ts
	//params.Set("before", "1716151200000") // Optional: Pagination of data to return records newer than the requested ts. The latest data will be returned when using before individually

	candles, err := client.GetCandlesticks(context.Background(), params)
	if err != nil {
		log.Fatalf("failed to get candlesticks: %v", err)
	}

	for _, c := range candles {
		fmt.Printf("%+v\n", c)
	}
}
