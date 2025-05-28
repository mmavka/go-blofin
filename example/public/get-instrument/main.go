// Example for testing GetInstruments from the Blofin public API.
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
	params.Set("instType", "SPOT") // Required: SPOT, SWAP, FUTURES, OPTION

	instruments, err := client.GetInstruments(context.Background(), params)
	if err != nil {
		slog.Error("failed to get instruments", "error", err)
		os.Exit(1)
	}

	for _, inst := range instruments {
		fmt.Printf("%+v\n", inst)
	}
}
