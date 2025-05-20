/**
 * @file: public.go
 * @description: Публичные сервисы RestClient для BloFin (стиль go-binance)
 * @dependencies: client.go, models.go
 * @created: 2025-05-19
 */

package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

var ErrMissingInstID = errors.New("instId is required")

// --- Instruments Service ---
type GetInstrumentsService struct {
	c      *RestClient
	instId string
}

func (c *RestClient) NewGetInstrumentsService() *GetInstrumentsService {
	return &GetInstrumentsService{c: c}
}

func (s *GetInstrumentsService) InstId(instId string) *GetInstrumentsService {
	s.instId = instId
	return s
}

func (s *GetInstrumentsService) Do(ctx context.Context) (*InstrumentsResponse, error) {
	resp := &InstrumentsResponse{}
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp)
	if s.instId != "" {
		req.SetQueryParam("instId", s.instId)
	}
	r, err := s.c.Request(ctx, "GET", "/api/v1/market/instruments", nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(r.Body(), resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return resp, nil
}

// --- Tickers Service ---
type GetTickersService struct {
	c      *RestClient
	instId string
}

func (c *RestClient) NewGetTickersService() *GetTickersService {
	return &GetTickersService{c: c}
}

func (s *GetTickersService) InstId(instId string) *GetTickersService {
	s.instId = instId
	return s
}

func (s *GetTickersService) Do(ctx context.Context) (*TickersResponse, error) {
	resp := &TickersResponse{}
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp)
	if s.instId != "" {
		req.SetQueryParam("instId", s.instId)
	}
	r, err := s.c.Request(ctx, "GET", "/api/v1/market/tickers", nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(r.Body(), resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return resp, nil
}

// --- OrderBook Service ---
type GetOrderBookService struct {
	c      *RestClient
	instId string
	size   int
}

func (c *RestClient) NewGetOrderBookService() *GetOrderBookService {
	return &GetOrderBookService{c: c}
}

func (s *GetOrderBookService) InstId(instId string) *GetOrderBookService {
	s.instId = instId
	return s
}

func (s *GetOrderBookService) Size(size int) *GetOrderBookService {
	s.size = size
	return s
}

func (s *GetOrderBookService) Do(ctx context.Context) (*OrderBookResponse, error) {
	if s.instId == "" {
		return nil, ErrMissingInstID
	}
	resp := &OrderBookResponse{}
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp)
	req.SetQueryParam("instId", s.instId)
	if s.size > 0 {
		req.SetQueryParam("size", fmt.Sprintf("%d", s.size))
	}
	r, err := req.Get("/api/v1/market/books")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(r.Body(), resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return resp, nil
}

// --- Trades Service ---
type GetTradesService struct {
	c      *RestClient
	instId string
	limit  int
}

func (c *RestClient) NewGetTradesService() *GetTradesService {
	return &GetTradesService{c: c}
}

func (s *GetTradesService) InstId(instId string) *GetTradesService {
	s.instId = instId
	return s
}

func (s *GetTradesService) Limit(limit int) *GetTradesService {
	s.limit = limit
	return s
}

func (s *GetTradesService) Do(ctx context.Context) (*TradesResponse, error) {
	if s.instId == "" {
		return nil, ErrMissingInstID
	}
	resp := &TradesResponse{}
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp).SetQueryParam("instId", s.instId)
	if s.limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", s.limit))
	}
	r, err := s.c.Request(ctx, "GET", "/api/v1/market/trades", nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(r.Body(), resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return resp, nil
}

// --- MarkPrice Service ---
type GetMarkPriceService struct {
	c      *RestClient
	instId string
}

func (c *RestClient) NewGetMarkPriceService() *GetMarkPriceService {
	return &GetMarkPriceService{c: c}
}

func (s *GetMarkPriceService) InstId(instId string) *GetMarkPriceService {
	s.instId = instId
	return s
}

func (s *GetMarkPriceService) Do(ctx context.Context) (*MarkPriceResponse, error) {
	resp := &MarkPriceResponse{}
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp)
	if s.instId != "" {
		req.SetQueryParam("instId", s.instId)
	}
	r, err := s.c.Request(ctx, "GET", "/api/v1/market/mark-price", nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(r.Body(), resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return resp, nil
}

// --- FundingRateHistory Service ---
type GetFundingRateHistoryService struct {
	c      *RestClient
	instId string
	before string
	after  string
	limit  int
}

func (c *RestClient) NewGetFundingRateHistoryService() *GetFundingRateHistoryService {
	return &GetFundingRateHistoryService{c: c}
}

func (s *GetFundingRateHistoryService) InstId(instId string) *GetFundingRateHistoryService {
	s.instId = instId
	return s
}

func (s *GetFundingRateHistoryService) Before(before string) *GetFundingRateHistoryService {
	s.before = before
	return s
}

func (s *GetFundingRateHistoryService) After(after string) *GetFundingRateHistoryService {
	s.after = after
	return s
}

func (s *GetFundingRateHistoryService) Limit(limit int) *GetFundingRateHistoryService {
	s.limit = limit
	return s
}

func (s *GetFundingRateHistoryService) Do(ctx context.Context) (*FundingRateResponse, error) {
	if s.instId == "" {
		return nil, ErrMissingInstID
	}
	resp := &FundingRateResponse{}
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp).SetQueryParam("instId", s.instId)
	if s.before != "" {
		req.SetQueryParam("before", s.before)
	}
	if s.after != "" {
		req.SetQueryParam("after", s.after)
	}
	if s.limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", s.limit))
	}
	r, err := req.Get("/api/v1/market/funding-rate-history")
	if err != nil {
		return nil, err
	}
	return r.Result().(*FundingRateResponse), nil
}

// --- Candles Service ---
type GetCandlesService struct {
	c      *RestClient
	instId string
	bar    string
	after  string
	before string
	limit  int
}

func (c *RestClient) NewGetCandlesService() *GetCandlesService {
	return &GetCandlesService{c: c}
}

func (s *GetCandlesService) InstId(instId string) *GetCandlesService {
	s.instId = instId
	return s
}

func (s *GetCandlesService) Bar(bar string) *GetCandlesService {
	s.bar = bar
	return s
}

func (s *GetCandlesService) After(after string) *GetCandlesService {
	s.after = after
	return s
}

func (s *GetCandlesService) Before(before string) *GetCandlesService {
	s.before = before
	return s
}

func (s *GetCandlesService) Limit(limit int) *GetCandlesService {
	s.limit = limit
	return s
}

func (s *GetCandlesService) Do(ctx context.Context) (*CandlesResponse, error) {
	if s.instId == "" {
		return nil, ErrMissingInstID
	}
	resp := &CandlesResponse{}
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp).SetQueryParam("instId", s.instId)
	if s.bar != "" {
		req.SetQueryParam("bar", s.bar)
	}
	if s.after != "" {
		req.SetQueryParam("after", s.after)
	}
	if s.before != "" {
		req.SetQueryParam("before", s.before)
	}
	if s.limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", s.limit))
	}
	r, err := s.c.Request(ctx, "GET", "/api/v1/market/candles", nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(r.Body(), resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return resp, nil
}
