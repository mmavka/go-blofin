// Package rest provides REST API client for public endpoints.
package rest

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"github.com/mmavka/go-blofin/internal/models"
)

// GetInstruments fetches the list of trading instruments from Blofin public API.
// Params can be used to filter instruments (see API docs).
// Returns a slice of Instrument or an error (including rate limit and API errors).
// Params: instId (optional)
func (c *Client) GetInstruments(ctx context.Context, params url.Values) ([]models.Instrument, error) {
	var instruments []models.Instrument
	err := c.doGet(ctx, EndpointInstruments, params, &instruments)
	if err != nil {
		return nil, err
	}
	return instruments, nil
}

// GetCandlesticks fetches candlestick data for a given instrument.
// Params: instId (required), bar, after, before, limit (see API docs)
func (c *Client) GetCandlesticks(ctx context.Context, params url.Values) ([]models.Candlestick, error) {
	// check required params
	if params.Get("instId") == "" {
		return nil, errors.New("instId is required")
	}
	if params.Get("limit") != "" {
		limit, err := strconv.Atoi(params.Get("limit"))
		if err != nil {
			return nil, errors.New("limit must be a number")
		}
		if limit > 1440 {
			return nil, errors.New("limit must be less than 1440")
		}
	}
	var rawData [][]string
	err := c.doGet(ctx, EndpointCandlesticks, params, &rawData)
	if err != nil {
		return nil, err
	}

	candles := make([]models.Candlestick, 0, len(rawData))
	for _, arr := range rawData {
		if len(arr) < 9 {
			continue // skip malformed
		}
		candles = append(candles, models.Candlestick{
			Ts:             arr[0],
			Open:           arr[1],
			High:           arr[2],
			Low:            arr[3],
			Close:          arr[4],
			Volume:         arr[5],
			VolumeCurrency: arr[6],
			VolumeQuote:    arr[7],
			Confirm:        arr[8],
		})
	}
	return candles, nil
}

// GetTickers fetches the latest price snapshot, best bid/ask, and 24h volume.
// Params: instId (optional)
func (c *Client) GetTickers(ctx context.Context, params url.Values) ([]models.Ticker, error) {
	var tickers []models.Ticker
	err := c.doGet(ctx, EndpointTickers, params, &tickers)
	if err != nil {
		return nil, err
	}
	return tickers, nil
}

// GetOrderBook fetches the order book for a given instrument.
// Params: instId (required), size (optional)
func (c *Client) GetOrderBook(ctx context.Context, params url.Values) (*models.OrderBook, error) {
	// check required params
	if params.Get("instId") == "" {
		return nil, errors.New("instId is required")
	}
	if params.Get("size") != "" {
		size, err := strconv.Atoi(params.Get("size"))
		if err != nil {
			return nil, errors.New("size must be a number")
		}
		if size > 100 {
			return nil, errors.New("size must be less than 100")
		}
	}
	var rawBooks []struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
		Ts   string     `json:"ts"`
	}
	err := c.doGet(ctx, EndpointOrderBook, params, &rawBooks)
	if err != nil {
		return nil, err
	}
	if len(rawBooks) == 0 {
		return nil, nil
	}
	ob := &models.OrderBook{
		Asks: make([]models.OrderBookLevel, 0, len(rawBooks[0].Asks)),
		Bids: make([]models.OrderBookLevel, 0, len(rawBooks[0].Bids)),
		Ts:   rawBooks[0].Ts,
	}
	for _, arr := range rawBooks[0].Asks {
		if len(arr) < 2 {
			continue
		}
		ob.Asks = append(ob.Asks, models.OrderBookLevel{Price: arr[0], Quantity: arr[1]})
	}
	for _, arr := range rawBooks[0].Bids {
		if len(arr) < 2 {
			continue
		}
		ob.Bids = append(ob.Bids, models.OrderBookLevel{Price: arr[0], Quantity: arr[1]})
	}
	return ob, nil
}

// GetTrades fetches recent transactions for a given instrument.
// Params: instId (required), limit (optional)
func (c *Client) GetTrades(ctx context.Context, params url.Values) ([]models.Trade, error) {
	// check required params
	if params.Get("instId") == "" {
		return nil, errors.New("instId is required")
	}
	if params.Get("limit") != "" {
		limit, err := strconv.Atoi(params.Get("limit"))
		if err != nil {
			return nil, errors.New("limit must be a number")
		}
		if limit > 100 {
			return nil, errors.New("limit must be less than 100")
		}
	}
	var trades []models.Trade
	err := c.doGet(ctx, EndpointTrades, params, &trades)
	if err != nil {
		return nil, err
	}
	return trades, nil
}

// GetMarkPrice fetches index and mark price for an instrument.
// Params: instId (optional)
func (c *Client) GetMarkPrice(ctx context.Context, params url.Values) ([]models.MarkPrice, error) {
	var prices []models.MarkPrice
	err := c.doGet(ctx, EndpointMarkPrice, params, &prices)
	if err != nil {
		return nil, err
	}
	return prices, nil
}

// GetFundingRate fetches funding rate for an instrument.
// Params: instId (optional)
func (c *Client) GetFundingRate(ctx context.Context, params url.Values) ([]models.FundingRate, error) {
	var rates []models.FundingRate
	err := c.doGet(ctx, EndpointFundingRate, params, &rates)
	if err != nil {
		return nil, err
	}
	return rates, nil
}

// GetFundingRateHistory fetches funding rate history for an instrument.
// Params: instId (required), before, after, limit (optional)
func (c *Client) GetFundingRateHistory(ctx context.Context, params url.Values) ([]models.FundingRate, error) {
	// check required params
	if params.Get("instId") == "" {
		return nil, errors.New("instId is required")
	}
	if params.Get("limit") != "" {
		limit, err := strconv.Atoi(params.Get("limit"))
		if err != nil {
			return nil, errors.New("limit must be a number")
		}
		if limit > 100 {
			return nil, errors.New("limit must be less than 100")
		}
	}
	var rates []models.FundingRate
	err := c.doGet(ctx, EndpointFundingRateHist, params, &rates)
	if err != nil {
		return nil, err
	}
	return rates, nil
}
