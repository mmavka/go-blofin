/**
 * @file: private.go
 * @description: Приватные сервисы RestClient для BloFin (стиль go-binance)
 * @dependencies: client.go, models.go, auth/signature.go
 * @created: 2025-05-19
 */

package rest

import (
	"context"
	"encoding/json"
)

// --- GetBalance Service ---
type GetBalanceService struct {
	c *RestClient
}

func (c *RestClient) NewGetBalanceService() *GetBalanceService {
	return &GetBalanceService{c: c}
}

func (s *GetBalanceService) Do(ctx context.Context) (*GetBalanceResponse, error) {
	resp := &GetBalanceResponse{}
	path := "/api/v1/account/balance"
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp)
	s.c.addAuthHeaders(req, "GET", path, "")
	_, err := req.Get(path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// --- GetPositions Service ---
type GetPositionsService struct {
	c *RestClient
}

func (c *RestClient) NewGetPositionsService() *GetPositionsService {
	return &GetPositionsService{c: c}
}

func (s *GetPositionsService) Do(ctx context.Context) (*PositionsResponse, error) {
	resp := &PositionsResponse{}
	path := "/api/v1/account/positions"
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp)
	s.c.addAuthHeaders(req, "GET", path, "")
	_, err := req.Get(path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// --- PlaceOrder Service ---
type PlaceOrderService struct {
	c     *RestClient
	order OrderRequest
}

func (c *RestClient) NewPlaceOrderService() *PlaceOrderService {
	return &PlaceOrderService{c: c}
}

func (s *PlaceOrderService) Order(order OrderRequest) *PlaceOrderService {
	s.order = order
	return s
}

func (s *PlaceOrderService) Do(ctx context.Context) (*OrderResponse, error) {
	resp := &OrderResponse{}
	path := "/api/v1/trade/order"
	bodyBytes, _ := json.Marshal(s.order)
	bodyStr := string(bodyBytes)
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp).SetBody(bodyStr)
	s.c.addAuthHeaders(req, "POST", path, bodyStr)
	_, err := req.Post(path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// --- CancelOrder Service ---
type CancelOrderService struct {
	c       *RestClient
	orderId string
}

func (c *RestClient) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c}
}

func (s *CancelOrderService) OrderId(orderId string) *CancelOrderService {
	s.orderId = orderId
	return s
}

func (s *CancelOrderService) Do(ctx context.Context) (*OrderResponse, error) {
	resp := &OrderResponse{}
	path := "/api/v1/trade/cancel"
	body := map[string]string{"orderId": s.orderId}
	bodyBytes, _ := json.Marshal(body)
	bodyStr := string(bodyBytes)
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp).SetBody(bodyStr)
	s.c.addAuthHeaders(req, "POST", path, bodyStr)
	_, err := req.Post(path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// --- GetActiveOrders Service ---
type GetActiveOrdersService struct {
	c *RestClient
}

func (c *RestClient) NewGetActiveOrdersService() *GetActiveOrdersService {
	return &GetActiveOrdersService{c: c}
}

func (s *GetActiveOrdersService) Do(ctx context.Context) (*OrderResponse, error) {
	resp := &OrderResponse{}
	path := "/api/v1/trade/orders"
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp)
	s.c.addAuthHeaders(req, "GET", path, "")
	_, err := req.Get(path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// --- GetWithdrawHistory Service ---
type GetWithdrawHistoryService struct {
	c          *RestClient
	currency   string
	withdrawId string
	txId       string
	state      string
	before     string
	after      string
	limit      string
}

func (c *RestClient) NewGetWithdrawHistoryService() *GetWithdrawHistoryService {
	return &GetWithdrawHistoryService{c: c}
}

func (s *GetWithdrawHistoryService) Currency(currency string) *GetWithdrawHistoryService {
	s.currency = currency
	return s
}
func (s *GetWithdrawHistoryService) WithdrawId(withdrawId string) *GetWithdrawHistoryService {
	s.withdrawId = withdrawId
	return s
}
func (s *GetWithdrawHistoryService) TxId(txId string) *GetWithdrawHistoryService {
	s.txId = txId
	return s
}
func (s *GetWithdrawHistoryService) State(state string) *GetWithdrawHistoryService {
	s.state = state
	return s
}
func (s *GetWithdrawHistoryService) Before(before string) *GetWithdrawHistoryService {
	s.before = before
	return s
}
func (s *GetWithdrawHistoryService) After(after string) *GetWithdrawHistoryService {
	s.after = after
	return s
}
func (s *GetWithdrawHistoryService) Limit(limit string) *GetWithdrawHistoryService {
	s.limit = limit
	return s
}

func (s *GetWithdrawHistoryService) Do(ctx context.Context) (*WithdrawHistoryResponse, error) {
	endpoint := "/api/v1/asset/withdrawal-history"
	params := map[string]string{}
	if s.currency != "" {
		params["currency"] = s.currency
	}
	if s.withdrawId != "" {
		params["withdrawId"] = s.withdrawId
	}
	if s.txId != "" {
		params["txId"] = s.txId
	}
	if s.state != "" {
		params["state"] = s.state
	}
	if s.before != "" {
		params["before"] = s.before
	}
	if s.after != "" {
		params["after"] = s.after
	}
	if s.limit != "" {
		params["limit"] = s.limit
	}
	req := s.c.httpClient.R().SetContext(ctx).SetQueryParams(params)
	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}
	var result WithdrawHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// --- GetTransferHistory Service ---
type GetTransferHistoryService struct {
	c           *RestClient
	currency    string
	fromAccount string
	toAccount   string
	before      string
	after       string
	limit       string
}

func (c *RestClient) NewGetTransferHistoryService() *GetTransferHistoryService {
	return &GetTransferHistoryService{c: c}
}

func (s *GetTransferHistoryService) Currency(currency string) *GetTransferHistoryService {
	s.currency = currency
	return s
}
func (s *GetTransferHistoryService) FromAccount(from string) *GetTransferHistoryService {
	s.fromAccount = from
	return s
}
func (s *GetTransferHistoryService) ToAccount(to string) *GetTransferHistoryService {
	s.toAccount = to
	return s
}
func (s *GetTransferHistoryService) Before(before string) *GetTransferHistoryService {
	s.before = before
	return s
}
func (s *GetTransferHistoryService) After(after string) *GetTransferHistoryService {
	s.after = after
	return s
}
func (s *GetTransferHistoryService) Limit(limit string) *GetTransferHistoryService {
	s.limit = limit
	return s
}

func (s *GetTransferHistoryService) Do(ctx context.Context) (*TransferHistoryResponse, error) {
	endpoint := "/api/v1/asset/bills"
	params := map[string]string{}
	if s.currency != "" {
		params["currency"] = s.currency
	}
	if s.fromAccount != "" {
		params["fromAccount"] = s.fromAccount
	}
	if s.toAccount != "" {
		params["toAccount"] = s.toAccount
	}
	if s.before != "" {
		params["before"] = s.before
	}
	if s.after != "" {
		params["after"] = s.after
	}
	if s.limit != "" {
		params["limit"] = s.limit
	}
	req := s.c.httpClient.R().SetContext(ctx).SetQueryParams(params)
	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}
	var result TransferHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// --- TransferFunds Service ---
type TransferFundsService struct {
	c        *RestClient
	currency string
	amount   string
	from     string
	to       string
	clientId string
}

func (c *RestClient) NewTransferFundsService() *TransferFundsService {
	return &TransferFundsService{c: c}
}

func (s *TransferFundsService) Currency(currency string) *TransferFundsService {
	s.currency = currency
	return s
}
func (s *TransferFundsService) Amount(amount string) *TransferFundsService {
	s.amount = amount
	return s
}
func (s *TransferFundsService) FromAccount(from string) *TransferFundsService {
	s.from = from
	return s
}
func (s *TransferFundsService) ToAccount(to string) *TransferFundsService {
	s.to = to
	return s
}
func (s *TransferFundsService) ClientId(clientId string) *TransferFundsService {
	s.clientId = clientId
	return s
}

func (s *TransferFundsService) Do(ctx context.Context) (*TransferResponse, error) {
	endpoint := "/api/v1/asset/transfer"
	body := TransferRequest{
		Currency:    s.currency,
		Amount:      s.amount,
		FromAccount: s.from,
		ToAccount:   s.to,
		ClientId:    s.clientId,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req := s.c.httpClient.R().SetContext(ctx).SetBody(bodyBytes).SetHeader("Content-Type", "application/json")
	s.c.addAuthHeaders(req, "POST", endpoint, string(bodyBytes))
	resp, err := req.Post(endpoint)
	if err != nil {
		return nil, err
	}
	var result TransferResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// --- GetBalances Service ---
type GetBalancesService struct {
	c           *RestClient
	accountType string
	currency    string
}

func (c *RestClient) NewGetBalancesService() *GetBalancesService {
	return &GetBalancesService{c: c}
}

func (s *GetBalancesService) AccountType(accountType string) *GetBalancesService {
	s.accountType = accountType
	return s
}
func (s *GetBalancesService) Currency(currency string) *GetBalancesService {
	s.currency = currency
	return s
}

func (s *GetBalancesService) Do(ctx context.Context) (*GetBalanceResponse, error) {
	endpoint := "/api/v1/asset/balances"
	params := map[string]string{"accountType": s.accountType}
	if s.currency != "" {
		params["currency"] = s.currency
	}
	req := s.c.httpClient.R().SetContext(ctx).SetQueryParams(params)
	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}
	var result GetBalanceResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
