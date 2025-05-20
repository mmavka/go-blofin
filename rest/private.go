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
	"errors"
	"fmt"
	"strings"

	"github.com/mmavka/go-blofin"
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

// GetAccountBalance returns service for getting futures account balance
type GetAccountBalanceService struct {
	c *RestClient
}

// Do executes the request
func (s *GetAccountBalanceService) Do(ctx context.Context) (*AccountBalanceResponse, error) {
	resp := &AccountBalanceResponse{}
	path := "/api/v1/account/balance"
	req := s.c.httpClient.R().SetContext(ctx).SetResult(resp)
	s.c.addAuthHeaders(req, "GET", path, "")
	_, err := req.Get(path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAccountBalance возвращает сервис для получения баланса фьючерсного аккаунта
func (c *RestClient) GetAccountBalance() *GetAccountBalanceService {
	return &GetAccountBalanceService{c: c}
}

// GetPositionsService service for getting position information
// InstId sets instrument ID (optional)
// Do executes the request
// GetPositions returns service for getting position information
type GetPositionsService struct {
	c      *RestClient
	instId string
}

// InstId устанавливает ID инструмента (опционально)
func (s *GetPositionsService) InstId(instId string) *GetPositionsService {
	s.instId = instId
	return s
}

// Do выполняет запрос
func (s *GetPositionsService) Do(ctx context.Context) (*GetPositionsResponse, error) {
	endpoint := "/api/v1/account/positions"
	params := make(map[string]string)

	if s.instId != "" {
		params["instId"] = s.instId
	}

	req := s.c.httpClient.R().SetContext(ctx)
	if len(params) > 0 {
		req.SetQueryParams(params)
	}

	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetPositionsResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPositions возвращает сервис для получения информации о позициях
func (c *RestClient) GetPositions() *GetPositionsService {
	return &GetPositionsService{c: c}
}

// GetMarginModeService service for getting margin mode
// Do executes the request
// GetMarginMode returns service for getting margin mode
type GetMarginModeService struct {
	c *RestClient
}

// Do выполняет запрос
func (s *GetMarginModeService) Do(ctx context.Context) (*GetMarginModeResponse, error) {
	endpoint := "/api/v1/account/margin-mode"
	req := s.c.httpClient.R().SetContext(ctx)
	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetMarginModeResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMarginMode возвращает сервис для получения режима маржи
func (c *RestClient) GetMarginMode() *GetMarginModeService {
	return &GetMarginModeService{c: c}
}

// SetMarginModeService service for setting margin mode
// MarginMode sets margin mode (cross/isolated)
// Do executes the request
// SetMarginMode returns service for setting margin mode
// SetMarginModeService сервис для установки режима маржи
type SetMarginModeService struct {
	c          *RestClient
	marginMode string
}

// MarginMode устанавливает режим маржи (cross/isolated)
func (s *SetMarginModeService) MarginMode(marginMode string) *SetMarginModeService {
	s.marginMode = marginMode
	return s
}

// Do выполняет запрос
func (s *SetMarginModeService) Do(ctx context.Context) (*GetMarginModeResponse, error) {
	if s.marginMode == "" {
		return nil, fmt.Errorf("marginMode required")
	}

	endpoint := "/api/v1/account/set-margin-mode"
	req := SetMarginModeRequest{
		MarginMode: s.marginMode,
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetMarginModeResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetMarginMode возвращает сервис для установки режима маржи
func (c *RestClient) SetMarginMode() *SetMarginModeService {
	return &SetMarginModeService{c: c}
}

// GetPositionModeService service for getting position mode
// Do executes the request
// GetPositionMode returns service for getting position mode
type GetPositionModeService struct {
	c *RestClient
}

// Do выполняет запрос
func (s *GetPositionModeService) Do(ctx context.Context) (*GetPositionModeResponse, error) {
	endpoint := "/api/v1/account/position-mode"
	req := s.c.httpClient.R().SetContext(ctx)
	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetPositionModeResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPositionMode возвращает сервис для получения режима позиций
func (c *RestClient) GetPositionMode() *GetPositionModeService {
	return &GetPositionModeService{c: c}
}

// SetPositionModeService service for setting position mode
// PositionMode sets position mode (net_mode/long_short_mode)
// Do executes the request
// SetPositionMode returns service for setting position mode
type SetPositionModeService struct {
	c            *RestClient
	positionMode string
}

// PositionMode устанавливает режим позиций (net_mode/long_short_mode)
func (s *SetPositionModeService) PositionMode(positionMode string) *SetPositionModeService {
	s.positionMode = positionMode
	return s
}

// Do выполняет запрос
func (s *SetPositionModeService) Do(ctx context.Context) (*GetPositionModeResponse, error) {
	if s.positionMode == "" {
		return nil, fmt.Errorf("positionMode required")
	}

	endpoint := "/api/v1/account/set-position-mode"
	req := SetPositionModeRequest{
		PositionMode: s.positionMode,
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetPositionModeResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetPositionMode возвращает сервис для установки режима позиций
func (c *RestClient) SetPositionMode() *SetPositionModeService {
	return &SetPositionModeService{c: c}
}

// GetLeverageInfoService service for getting leverage information (deprecated)
type GetLeverageInfoService struct {
	c          *RestClient
	instId     string
	marginMode string
}

// InstId sets instrument ID
func (s *GetLeverageInfoService) InstId(instId string) *GetLeverageInfoService {
	s.instId = instId
	return s
}

// MarginMode sets margin mode (cross/isolated)
func (s *GetLeverageInfoService) MarginMode(marginMode string) *GetLeverageInfoService {
	s.marginMode = marginMode
	return s
}

// Do выполняет запрос
func (s *GetLeverageInfoService) Do(ctx context.Context) (*GetLeverageInfoResponse, error) {
	if s.instId == "" {
		return nil, fmt.Errorf("instId required")
	}
	if s.marginMode == "" {
		return nil, fmt.Errorf("marginMode required")
	}

	endpoint := "/api/v1/account/leverage-info"
	params := map[string]string{
		"instId":     s.instId,
		"marginMode": s.marginMode,
	}

	req := s.c.httpClient.R().SetContext(ctx).SetQueryParams(params)
	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetLeverageInfoResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetLeverageInfo returns service for getting leverage information (deprecated)
func (c *RestClient) GetLeverageInfo() *GetLeverageInfoService {
	return &GetLeverageInfoService{c: c}
}

// GetBatchLeverageInfoService service for getting leverage information for multiple instruments
type GetBatchLeverageInfoService struct {
	c          *RestClient
	instIds    []string
	marginMode string
}

// NewGetBatchLeverageInfoService creates a new service
func (c *RestClient) NewGetBatchLeverageInfoService() *GetBatchLeverageInfoService {
	return &GetBatchLeverageInfoService{c: c}
}

// InstIds sets list of instruments (not more than 20)
func (s *GetBatchLeverageInfoService) InstIds(instIds []string) *GetBatchLeverageInfoService {
	s.instIds = instIds
	return s
}

// MarginMode sets margin mode
func (s *GetBatchLeverageInfoService) MarginMode(marginMode string) *GetBatchLeverageInfoService {
	s.marginMode = marginMode
	return s
}

// Do executes request
func (s *GetBatchLeverageInfoService) Do(ctx context.Context) (*GetBatchLeverageInfoResponse, error) {
	if len(s.instIds) == 0 {
		return nil, fmt.Errorf("instIds required")
	}
	if len(s.instIds) > 20 {
		return nil, fmt.Errorf("too many instIds (max 20)")
	}
	if s.marginMode == "" {
		return nil, fmt.Errorf("marginMode required")
	}

	resp := &GetBatchLeverageInfoResponse{}
	params := map[string]string{
		"instId":     strings.Join(s.instIds, ","),
		"marginMode": s.marginMode,
	}

	path := "/api/v1/account/batch-leverage-info"
	req := s.c.httpClient.R().SetContext(ctx).SetQueryParams(params).SetResult(resp)
	s.c.addAuthHeaders(req, "GET", path, "")
	_, err := req.Get(path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SetLeverageService service for setting leverage
type SetLeverageService struct {
	c            *RestClient
	instId       string
	leverage     string
	marginMode   string
	positionSide string
}

// NewSetLeverageService creates a new service
func (c *RestClient) NewSetLeverageService() *SetLeverageService {
	return &SetLeverageService{c: c}
}

// InstId sets instrument ID
func (s *SetLeverageService) InstId(instId string) *SetLeverageService {
	s.instId = instId
	return s
}

// Leverage sets leverage value
func (s *SetLeverageService) Leverage(leverage string) *SetLeverageService {
	s.leverage = leverage
	return s
}

// MarginMode sets margin mode
func (s *SetLeverageService) MarginMode(marginMode string) *SetLeverageService {
	s.marginMode = marginMode
	return s
}

// PositionSide sets position side (optional)
func (s *SetLeverageService) PositionSide(positionSide string) *SetLeverageService {
	s.positionSide = positionSide
	return s
}

// Do executes request
func (s *SetLeverageService) Do(ctx context.Context) (*SetLeverageResponse, error) {
	if s.instId == "" {
		return nil, fmt.Errorf("instId required")
	}
	if s.leverage == "" {
		return nil, fmt.Errorf("leverage required")
	}
	if s.marginMode == "" {
		return nil, fmt.Errorf("marginMode required")
	}

	req := SetLeverageRequest{
		InstId:     s.instId,
		Leverage:   s.leverage,
		MarginMode: s.marginMode,
	}

	if s.positionSide != "" {
		req.PositionSide = s.positionSide
	}

	endpoint := "/api/v1/account/set-leverage"
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result SetLeverageResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PlaceOrderService service for placing an order
type PlaceOrderService struct {
	c     *RestClient
	order PlaceOrderRequest
}

// NewPlaceOrderService creates a new service
func (c *RestClient) NewPlaceOrderService() *PlaceOrderService {
	return &PlaceOrderService{c: c}
}

// InstId sets instrument ID
func (s *PlaceOrderService) InstId(instId string) *PlaceOrderService {
	s.order.InstID = instId
	return s
}

// MarginMode sets margin mode
func (s *PlaceOrderService) MarginMode(marginMode string) *PlaceOrderService {
	s.order.MarginMode = marginMode
	return s
}

// PositionSide sets position side
func (s *PlaceOrderService) PositionSide(positionSide string) *PlaceOrderService {
	s.order.PositionSide = positionSide
	return s
}

// Side sets order side
func (s *PlaceOrderService) Side(side string) *PlaceOrderService {
	s.order.Side = side
	return s
}

// OrderType sets order type
func (s *PlaceOrderService) OrderType(orderType string) *PlaceOrderService {
	s.order.OrderType = orderType
	return s
}

// Price sets order price
func (s *PlaceOrderService) Price(price string) *PlaceOrderService {
	s.order.Price = price
	return s
}

// Size sets order size
func (s *PlaceOrderService) Size(size string) *PlaceOrderService {
	s.order.Size = size
	return s
}

// ReduceOnly sets reduceOnly flag
func (s *PlaceOrderService) ReduceOnly(reduceOnly string) *PlaceOrderService {
	s.order.ReduceOnly = reduceOnly
	return s
}

// ClientOrderId sets client order ID
func (s *PlaceOrderService) ClientOrderId(clientOrderId string) *PlaceOrderService {
	s.order.ClientOrderId = clientOrderId
	return s
}

// TakeProfitParams sets take-profit parameters
func (s *PlaceOrderService) TakeProfitParams(triggerPrice, orderPrice string) *PlaceOrderService {
	s.order.TpTriggerPrice = triggerPrice
	s.order.TpOrderPrice = orderPrice
	return s
}

// StopLossParams sets stop-loss parameters
func (s *PlaceOrderService) StopLossParams(triggerPrice, orderPrice string) *PlaceOrderService {
	s.order.SlTriggerPrice = triggerPrice
	s.order.SlOrderPrice = orderPrice
	return s
}

// BrokerId sets broker ID
func (s *PlaceOrderService) BrokerId(brokerId string) *PlaceOrderService {
	s.order.BrokerId = brokerId
	return s
}

// Do executes request
func (s *PlaceOrderService) Do(ctx context.Context) (*PlaceOrderResponse, error) {
	if s.order.InstID == "" {
		return nil, fmt.Errorf("instId required")
	}
	if s.order.MarginMode == "" {
		return nil, fmt.Errorf("marginMode required")
	}
	if s.order.PositionSide == "" {
		return nil, fmt.Errorf("positionSide required")
	}
	if s.order.Side == "" {
		return nil, fmt.Errorf("side required")
	}
	if s.order.OrderType == "" {
		return nil, fmt.Errorf("orderType required")
	}
	if s.order.Size == "" {
		return nil, fmt.Errorf("size required")
	}
	if s.order.OrderType != blofin.OrderTypeMarket && s.order.Price == "" {
		return nil, fmt.Errorf("price required for non-market orders")
	}

	endpoint := "/api/v1/trade/order"
	bodyBytes, err := json.Marshal(s.order)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result PlaceOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PlaceBatchOrdersService service for placing multiple orders
type PlaceBatchOrdersService struct {
	c      *RestClient
	orders []PlaceOrderRequest
}

// NewPlaceBatchOrdersService creates a new service
func (c *RestClient) NewPlaceBatchOrdersService() *PlaceBatchOrdersService {
	return &PlaceBatchOrdersService{c: c}
}

// Orders sets list of orders
func (s *PlaceBatchOrdersService) Orders(orders []PlaceOrderRequest) *PlaceBatchOrdersService {
	s.orders = orders
	return s
}

// AddOrder adds an order to the list
func (s *PlaceBatchOrdersService) AddOrder(order PlaceOrderRequest) *PlaceBatchOrdersService {
	s.orders = append(s.orders, order)
	return s
}

// Do executes request
func (s *PlaceBatchOrdersService) Do(ctx context.Context) (*BatchOrdersResponse, error) {
	if len(s.orders) == 0 {
		return nil, fmt.Errorf("at least one order required")
	}
	if len(s.orders) > 20 {
		return nil, fmt.Errorf("too many orders (max 20)")
	}

	// Check that all orders are for the same instrument
	firstInstId := s.orders[0].InstID
	for _, order := range s.orders[1:] {
		if order.InstID != firstInstId {
			return nil, fmt.Errorf("all orders must have the same instId")
		}
	}

	// Check that all orders have required parameters
	for i, order := range s.orders {
		if order.InstID == "" {
			return nil, fmt.Errorf("order %d: instId required", i)
		}
		if order.MarginMode == "" {
			return nil, fmt.Errorf("order %d: marginMode required", i)
		}
		if order.PositionSide == "" {
			return nil, fmt.Errorf("order %d: positionSide required", i)
		}
		if order.Side == "" {
			return nil, fmt.Errorf("order %d: side required", i)
		}
		if order.OrderType == "" {
			return nil, fmt.Errorf("order %d: orderType required", i)
		}
		if order.Size == "" {
			return nil, fmt.Errorf("order %d: size required", i)
		}
		if order.OrderType != blofin.OrderTypeMarket && order.Price == "" {
			return nil, fmt.Errorf("order %d: price required for non-market orders", i)
		}
	}

	endpoint := "/api/v1/trade/batch-orders"
	bodyBytes, err := json.Marshal(s.orders)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result BatchOrdersResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PlaceTPSLOrderService service for placing TPSL order
type PlaceTPSLOrderService struct {
	c     *RestClient
	order PlaceTPSLOrderRequest
}

// NewPlaceTPSLOrderService creates a new service
func (c *RestClient) NewPlaceTPSLOrderService() *PlaceTPSLOrderService {
	return &PlaceTPSLOrderService{c: c}
}

// InstId sets instrument ID
func (s *PlaceTPSLOrderService) InstId(instId string) *PlaceTPSLOrderService {
	s.order.InstID = instId
	return s
}

// MarginMode sets margin mode
func (s *PlaceTPSLOrderService) MarginMode(marginMode string) *PlaceTPSLOrderService {
	s.order.MarginMode = marginMode
	return s
}

// PositionSide sets position side
func (s *PlaceTPSLOrderService) PositionSide(positionSide string) *PlaceTPSLOrderService {
	s.order.PositionSide = positionSide
	return s
}

// Side sets order side
func (s *PlaceTPSLOrderService) Side(side string) *PlaceTPSLOrderService {
	s.order.Side = side
	return s
}

// TakeProfitParams sets take-profit parameters
func (s *PlaceTPSLOrderService) TakeProfitParams(triggerPrice, orderPrice string) *PlaceTPSLOrderService {
	s.order.TpTriggerPrice = triggerPrice
	s.order.TpOrderPrice = orderPrice
	return s
}

// StopLossParams sets stop-loss parameters
func (s *PlaceTPSLOrderService) StopLossParams(triggerPrice, orderPrice string) *PlaceTPSLOrderService {
	s.order.SlTriggerPrice = triggerPrice
	s.order.SlOrderPrice = orderPrice
	return s
}

// Size sets order size
func (s *PlaceTPSLOrderService) Size(size string) *PlaceTPSLOrderService {
	s.order.Size = size
	return s
}

// ReduceOnly sets reduceOnly flag
func (s *PlaceTPSLOrderService) ReduceOnly(reduceOnly string) *PlaceTPSLOrderService {
	s.order.ReduceOnly = reduceOnly
	return s
}

// ClientOrderId sets client order ID
func (s *PlaceTPSLOrderService) ClientOrderId(clientOrderId string) *PlaceTPSLOrderService {
	s.order.ClientOrderId = clientOrderId
	return s
}

// BrokerId sets broker ID
func (s *PlaceTPSLOrderService) BrokerId(brokerId string) *PlaceTPSLOrderService {
	s.order.BrokerId = brokerId
	return s
}

// Do executes request
func (s *PlaceTPSLOrderService) Do(ctx context.Context) (*PlaceTPSLOrderResponse, error) {
	if s.order.InstID == "" {
		return nil, fmt.Errorf("instId required")
	}
	if s.order.MarginMode == "" {
		return nil, fmt.Errorf("marginMode required")
	}
	if s.order.PositionSide == "" {
		return nil, fmt.Errorf("positionSide required")
	}
	if s.order.Side == "" {
		return nil, fmt.Errorf("side required")
	}
	if s.order.Size == "" {
		return nil, fmt.Errorf("size required")
	}
	if s.order.TpTriggerPrice == "" && s.order.SlTriggerPrice == "" {
		return nil, fmt.Errorf("at least one of tpTriggerPrice or slTriggerPrice required")
	}
	if s.order.TpTriggerPrice != "" && s.order.TpOrderPrice == "" {
		return nil, fmt.Errorf("tpOrderPrice required when tpTriggerPrice is set")
	}
	if s.order.SlTriggerPrice != "" && s.order.SlOrderPrice == "" {
		return nil, fmt.Errorf("slOrderPrice required when slTriggerPrice is set")
	}

	endpoint := "/api/v1/trade/order-tpsl"
	bodyBytes, err := json.Marshal(s.order)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result PlaceTPSLOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PlaceAlgoOrderService service for placing algo order
type PlaceAlgoOrderService struct {
	c     *RestClient
	order PlaceAlgoOrderRequest
}

// NewPlaceAlgoOrderService creates a new service
func (c *RestClient) NewPlaceAlgoOrderService() *PlaceAlgoOrderService {
	return &PlaceAlgoOrderService{c: c}
}

// InstId sets instrument ID
func (s *PlaceAlgoOrderService) InstId(instId string) *PlaceAlgoOrderService {
	s.order.InstID = instId
	return s
}

// MarginMode sets margin mode
func (s *PlaceAlgoOrderService) MarginMode(marginMode string) *PlaceAlgoOrderService {
	s.order.MarginMode = marginMode
	return s
}

// PositionSide sets position side
func (s *PlaceAlgoOrderService) PositionSide(positionSide string) *PlaceAlgoOrderService {
	s.order.PositionSide = positionSide
	return s
}

// Side sets order side
func (s *PlaceAlgoOrderService) Side(side string) *PlaceAlgoOrderService {
	s.order.Side = side
	return s
}

// Size sets order size
func (s *PlaceAlgoOrderService) Size(size string) *PlaceAlgoOrderService {
	s.order.Size = size
	return s
}

// ClientOrderId sets client order ID
func (s *PlaceAlgoOrderService) ClientOrderId(clientOrderId string) *PlaceAlgoOrderService {
	s.order.ClientOrderId = clientOrderId
	return s
}

// OrderPrice sets order price
func (s *PlaceAlgoOrderService) OrderPrice(orderPrice string) *PlaceAlgoOrderService {
	s.order.OrderPrice = orderPrice
	return s
}

// OrderType sets order type
func (s *PlaceAlgoOrderService) OrderType(orderType string) *PlaceAlgoOrderService {
	s.order.OrderType = orderType
	return s
}

// TriggerPrice sets trigger price
func (s *PlaceAlgoOrderService) TriggerPrice(triggerPrice string) *PlaceAlgoOrderService {
	s.order.TriggerPrice = triggerPrice
	return s
}

// TriggerPriceType sets trigger price type
func (s *PlaceAlgoOrderService) TriggerPriceType(triggerPriceType string) *PlaceAlgoOrderService {
	s.order.TriggerPriceType = triggerPriceType
	return s
}

// ReduceOnly sets reduceOnly flag
func (s *PlaceAlgoOrderService) ReduceOnly(reduceOnly string) *PlaceAlgoOrderService {
	s.order.ReduceOnly = reduceOnly
	return s
}

// BrokerId sets broker ID
func (s *PlaceAlgoOrderService) BrokerId(brokerId string) *PlaceAlgoOrderService {
	s.order.BrokerId = brokerId
	return s
}

// AttachAlgoOrders sets attached TP/SL orders
func (s *PlaceAlgoOrderService) AttachAlgoOrders(orders []AttachAlgoOrder) *PlaceAlgoOrderService {
	s.order.AttachAlgoOrders = orders
	return s
}

// Do executes request
func (s *PlaceAlgoOrderService) Do(ctx context.Context) (*PlaceAlgoOrderResponse, error) {
	if s.order.InstID == "" {
		return nil, fmt.Errorf("instId required")
	}
	if s.order.MarginMode == "" {
		return nil, fmt.Errorf("marginMode required")
	}
	if s.order.PositionSide == "" {
		return nil, fmt.Errorf("positionSide required")
	}
	if s.order.Side == "" {
		return nil, fmt.Errorf("side required")
	}
	if s.order.Size == "" {
		return nil, fmt.Errorf("size required")
	}
	if s.order.OrderType == "" {
		return nil, fmt.Errorf("orderType required")
	}
	if s.order.TriggerPrice == "" {
		return nil, fmt.Errorf("triggerPrice required")
	}

	endpoint := "/api/v1/trade/order-algo"
	bodyBytes, err := json.Marshal(s.order)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result PlaceAlgoOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelOrderService service for canceling order
type CancelOrderService struct {
	c     *RestClient
	order CancelOrderRequest
}

// NewCancelOrderService creates a new service
func (c *RestClient) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c}
}

// InstId sets instrument ID
func (s *CancelOrderService) InstId(instId string) *CancelOrderService {
	s.order.InstID = instId
	return s
}

// OrderId sets order ID
func (s *CancelOrderService) OrderId(orderId string) *CancelOrderService {
	s.order.OrderId = orderId
	return s
}

// ClientOrderId sets client order ID
func (s *CancelOrderService) ClientOrderId(clientOrderId string) *CancelOrderService {
	s.order.ClientOrderId = clientOrderId
	return s
}

// Do executes request
func (s *CancelOrderService) Do(ctx context.Context) (*CancelOrderResponse, error) {
	if s.order.OrderId == "" && s.order.ClientOrderId == "" {
		return nil, fmt.Errorf("either orderId or clientOrderId required")
	}

	endpoint := "/api/v1/trade/cancel-order"
	bodyBytes, err := json.Marshal(s.order)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result CancelOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelBatchOrdersService service for canceling multiple orders
type CancelBatchOrdersService struct {
	c      *RestClient
	orders []CancelOrderRequest
}

// NewCancelBatchOrdersService creates a new service
func (c *RestClient) NewCancelBatchOrdersService() *CancelBatchOrdersService {
	return &CancelBatchOrdersService{c: c}
}

// Orders sets list of orders to cancel
func (s *CancelBatchOrdersService) Orders(orders []CancelOrderRequest) *CancelBatchOrdersService {
	s.orders = orders
	return s
}

// AddOrder adds order to list to cancel
func (s *CancelBatchOrdersService) AddOrder(order CancelOrderRequest) *CancelBatchOrdersService {
	s.orders = append(s.orders, order)
	return s
}

// Do executes request
func (s *CancelBatchOrdersService) Do(ctx context.Context) (*CancelBatchOrdersResponse, error) {
	if len(s.orders) == 0 {
		return nil, fmt.Errorf("at least one order required")
	}
	if len(s.orders) > 20 {
		return nil, fmt.Errorf("too many orders (max 20)")
	}

	// Check that all orders are for the same instrument
	firstInstId := s.orders[0].InstID
	for _, order := range s.orders[1:] {
		if order.InstID != firstInstId {
			return nil, fmt.Errorf("all orders must have the same instId")
		}
	}

	// Check that all orders have required parameters
	for i, order := range s.orders {
		if order.OrderId == "" && order.ClientOrderId == "" {
			return nil, fmt.Errorf("order %d: either orderId or clientOrderId required", i)
		}
	}

	endpoint := "/api/v1/trade/cancel-batch-orders"
	bodyBytes, err := json.Marshal(s.orders)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result CancelBatchOrdersResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelTPSLOrderService service for canceling TPSL orders
type CancelTPSLOrderService struct {
	c      *RestClient
	orders []CancelTPSLOrderRequest
}

// NewCancelTPSLOrderService creates a new service
func (c *RestClient) NewCancelTPSLOrderService() *CancelTPSLOrderService {
	return &CancelTPSLOrderService{c: c}
}

// Orders sets list of orders to cancel
func (s *CancelTPSLOrderService) Orders(orders []CancelTPSLOrderRequest) *CancelTPSLOrderService {
	s.orders = orders
	return s
}

// AddOrder adds order to list to cancel
func (s *CancelTPSLOrderService) AddOrder(order CancelTPSLOrderRequest) *CancelTPSLOrderService {
	s.orders = append(s.orders, order)
	return s
}

// Do executes request
func (s *CancelTPSLOrderService) Do(ctx context.Context) (*CancelTPSLOrderResponse, error) {
	if len(s.orders) == 0 {
		return nil, fmt.Errorf("at least one order required")
	}
	if len(s.orders) > 20 {
		return nil, fmt.Errorf("too many orders (max 20)")
	}

	// Check that all orders are for the same instrument
	firstInstId := s.orders[0].InstID
	for _, order := range s.orders[1:] {
		if order.InstID != firstInstId {
			return nil, fmt.Errorf("all orders must have the same instId")
		}
	}

	// Check that all orders have required parameters
	for i, order := range s.orders {
		if order.TpslId == "" && order.ClientOrderId == "" {
			return nil, fmt.Errorf("order %d: either tpslId or clientOrderId required", i)
		}
	}

	endpoint := "/api/v1/trade/cancel-tpsl"
	bodyBytes, err := json.Marshal(s.orders)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result CancelTPSLOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelAlgoOrderService service for canceling algo order
type CancelAlgoOrderService struct {
	c     *RestClient
	order CancelAlgoOrderRequest
}

// NewCancelAlgoOrderService creates a new service
func (c *RestClient) NewCancelAlgoOrderService() *CancelAlgoOrderService {
	return &CancelAlgoOrderService{c: c}
}

// InstId sets instrument ID
func (s *CancelAlgoOrderService) InstId(instId string) *CancelAlgoOrderService {
	s.order.InstID = instId
	return s
}

// AlgoId sets algo order ID
func (s *CancelAlgoOrderService) AlgoId(algoId string) *CancelAlgoOrderService {
	s.order.AlgoId = algoId
	return s
}

// ClientOrderId sets client order ID
func (s *CancelAlgoOrderService) ClientOrderId(clientOrderId string) *CancelAlgoOrderService {
	s.order.ClientOrderId = clientOrderId
	return s
}

// Do executes request
func (s *CancelAlgoOrderService) Do(ctx context.Context) (*CancelAlgoOrderResponse, error) {
	if s.order.AlgoId == "" && s.order.ClientOrderId == "" {
		return nil, fmt.Errorf("either algoId or clientOrderId required")
	}

	endpoint := "/api/v1/trade/cancel-algo"
	bodyBytes, err := json.Marshal(s.order)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result CancelAlgoOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPendingOrdersService service for getting active orders
type GetPendingOrdersService struct {
	c         *RestClient
	instId    string
	orderType string
	state     string
	after     string
	before    string
	limit     string
}

// NewGetPendingOrdersService creates a new service
func (c *RestClient) NewGetPendingOrdersService() *GetPendingOrdersService {
	return &GetPendingOrdersService{c: c}
}

// InstId sets instrument ID
func (s *GetPendingOrdersService) InstId(instId string) *GetPendingOrdersService {
	s.instId = instId
	return s
}

// OrderType sets order type
func (s *GetPendingOrdersService) OrderType(orderType string) *GetPendingOrdersService {
	s.orderType = orderType
	return s
}

// State sets order state
func (s *GetPendingOrdersService) State(state string) *GetPendingOrdersService {
	s.state = state
	return s
}

// After sets after parameter for pagination
func (s *GetPendingOrdersService) After(after string) *GetPendingOrdersService {
	s.after = after
	return s
}

// Before sets before parameter for pagination
func (s *GetPendingOrdersService) Before(before string) *GetPendingOrdersService {
	s.before = before
	return s
}

// Limit sets limit results
func (s *GetPendingOrdersService) Limit(limit string) *GetPendingOrdersService {
	s.limit = limit
	return s
}

// Do executes request
func (s *GetPendingOrdersService) Do(ctx context.Context) (*GetPendingOrdersResponse, error) {
	endpoint := "/api/v1/trade/orders-pending"
	params := make(map[string]string)

	if s.instId != "" {
		params["instId"] = s.instId
	}
	if s.orderType != "" {
		params["orderType"] = s.orderType
	}
	if s.state != "" {
		params["state"] = s.state
	}
	if s.after != "" {
		params["after"] = s.after
	}
	if s.before != "" {
		params["before"] = s.before
	}
	if s.limit != "" {
		params["limit"] = s.limit
	}

	req := s.c.httpClient.R().SetContext(ctx)
	if len(params) > 0 {
		req.SetQueryParams(params)
	}

	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetPendingOrdersResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPendingTPSLOrdersService service for getting active TPSL orders
type GetPendingTPSLOrdersService struct {
	c             *RestClient
	instId        string
	tpslId        string
	clientOrderId string
	after         string
	before        string
	limit         string
}

// NewGetPendingTPSLOrdersService creates a new service
func (c *RestClient) NewGetPendingTPSLOrdersService() *GetPendingTPSLOrdersService {
	return &GetPendingTPSLOrdersService{c: c}
}

// InstId sets instrument ID
func (s *GetPendingTPSLOrdersService) InstId(instId string) *GetPendingTPSLOrdersService {
	s.instId = instId
	return s
}

// TpslId sets TPSL order ID
func (s *GetPendingTPSLOrdersService) TpslId(tpslId string) *GetPendingTPSLOrdersService {
	s.tpslId = tpslId
	return s
}

// ClientOrderId sets client order ID
func (s *GetPendingTPSLOrdersService) ClientOrderId(clientOrderId string) *GetPendingTPSLOrdersService {
	s.clientOrderId = clientOrderId
	return s
}

// After sets after parameter for pagination
func (s *GetPendingTPSLOrdersService) After(after string) *GetPendingTPSLOrdersService {
	s.after = after
	return s
}

// Before sets before parameter for pagination
func (s *GetPendingTPSLOrdersService) Before(before string) *GetPendingTPSLOrdersService {
	s.before = before
	return s
}

// Limit sets limit results
func (s *GetPendingTPSLOrdersService) Limit(limit string) *GetPendingTPSLOrdersService {
	s.limit = limit
	return s
}

// Do executes request
func (s *GetPendingTPSLOrdersService) Do(ctx context.Context) (*GetPendingTPSLOrdersResponse, error) {
	endpoint := "/api/v1/trade/orders-tpsl-pending"
	params := make(map[string]string)

	if s.instId != "" {
		params["instId"] = s.instId
	}
	if s.tpslId != "" {
		params["tpslId"] = s.tpslId
	}
	if s.clientOrderId != "" {
		params["clientOrderId"] = s.clientOrderId
	}
	if s.after != "" {
		params["after"] = s.after
	}
	if s.before != "" {
		params["before"] = s.before
	}
	if s.limit != "" {
		params["limit"] = s.limit
	}

	req := s.c.httpClient.R().SetContext(ctx)
	if len(params) > 0 {
		req.SetQueryParams(params)
	}

	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetPendingTPSLOrdersResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPendingAlgoOrdersService service for getting active algo orders
type GetPendingAlgoOrdersService struct {
	c             *RestClient
	instId        string
	algoId        string
	clientOrderId string
	orderType     string
	after         string
	before        string
	limit         string
}

// NewGetPendingAlgoOrdersService creates a new service
func (c *RestClient) NewGetPendingAlgoOrdersService() *GetPendingAlgoOrdersService {
	return &GetPendingAlgoOrdersService{c: c}
}

// InstId sets instrument ID
func (s *GetPendingAlgoOrdersService) InstId(instId string) *GetPendingAlgoOrdersService {
	s.instId = instId
	return s
}

// AlgoId sets algo order ID
func (s *GetPendingAlgoOrdersService) AlgoId(algoId string) *GetPendingAlgoOrdersService {
	s.algoId = algoId
	return s
}

// ClientOrderId sets client order ID
func (s *GetPendingAlgoOrdersService) ClientOrderId(clientOrderId string) *GetPendingAlgoOrdersService {
	s.clientOrderId = clientOrderId
	return s
}

// OrderType sets order type
func (s *GetPendingAlgoOrdersService) OrderType(orderType string) *GetPendingAlgoOrdersService {
	s.orderType = orderType
	return s
}

// After sets after parameter for pagination
func (s *GetPendingAlgoOrdersService) After(after string) *GetPendingAlgoOrdersService {
	s.after = after
	return s
}

// Before sets before parameter for pagination
func (s *GetPendingAlgoOrdersService) Before(before string) *GetPendingAlgoOrdersService {
	s.before = before
	return s
}

// Limit sets limit results
func (s *GetPendingAlgoOrdersService) Limit(limit string) *GetPendingAlgoOrdersService {
	s.limit = limit
	return s
}

// Do executes request
func (s *GetPendingAlgoOrdersService) Do(ctx context.Context) (*GetPendingAlgoOrdersResponse, error) {
	if s.orderType == "" {
		return nil, fmt.Errorf("orderType required")
	}
	if s.after != "" && s.before != "" {
		return nil, fmt.Errorf("after and before cannot be used simultaneously")
	}

	endpoint := "/api/v1/trade/orders-algo-pending"
	params := make(map[string]string)

	if s.instId != "" {
		params["instId"] = s.instId
	}
	if s.algoId != "" {
		params["algoId"] = s.algoId
	}
	if s.clientOrderId != "" {
		params["clientOrderId"] = s.clientOrderId
	}
	if s.orderType != "" {
		params["orderType"] = s.orderType
	}
	if s.after != "" {
		params["after"] = s.after
	}
	if s.before != "" {
		params["before"] = s.before
	}
	if s.limit != "" {
		params["limit"] = s.limit
	}

	req := s.c.httpClient.R().SetContext(ctx)
	if len(params) > 0 {
		req.SetQueryParams(params)
	}

	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetPendingAlgoOrdersResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ClosePositionService service for closing position
type ClosePositionService struct {
	c     *RestClient
	order ClosePositionRequest
}

// NewClosePositionService creates a new service
func (c *RestClient) NewClosePositionService() *ClosePositionService {
	return &ClosePositionService{c: c}
}

// InstId sets instrument ID
func (s *ClosePositionService) InstId(instId string) *ClosePositionService {
	s.order.InstID = instId
	return s
}

// MarginMode sets margin mode
func (s *ClosePositionService) MarginMode(marginMode string) *ClosePositionService {
	s.order.MarginMode = marginMode
	return s
}

// PositionSide sets position side
func (s *ClosePositionService) PositionSide(positionSide string) *ClosePositionService {
	s.order.PositionSide = positionSide
	return s
}

// ClientOrderId sets client order ID
func (s *ClosePositionService) ClientOrderId(clientOrderId string) *ClosePositionService {
	s.order.ClientOrderId = clientOrderId
	return s
}

// BrokerId sets broker ID
func (s *ClosePositionService) BrokerId(brokerId string) *ClosePositionService {
	s.order.BrokerId = brokerId
	return s
}

// Do executes request
func (s *ClosePositionService) Do(ctx context.Context) (*ClosePositionResponse, error) {
	if s.order.InstID == "" {
		return nil, fmt.Errorf("instId required")
	}
	if s.order.MarginMode == "" {
		return nil, fmt.Errorf("marginMode required")
	}
	if s.order.PositionSide == "" {
		return nil, fmt.Errorf("positionSide required")
	}

	endpoint := "/api/v1/trade/close-position"
	bodyBytes, err := json.Marshal(s.order)
	if err != nil {
		return nil, err
	}

	r := s.c.httpClient.R().SetContext(ctx).SetBody(string(bodyBytes))
	s.c.addAuthHeaders(r, "POST", endpoint, string(bodyBytes))
	resp, err := r.Post(endpoint)
	if err != nil {
		return nil, err
	}

	var result ClosePositionResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetOrderHistoryService service for getting order history
type GetOrderHistoryService struct {
	c         *RestClient
	instId    string
	orderType string
	state     string
	after     string
	before    string
	begin     string
	end       string
	limit     string
}

// NewGetOrderHistoryService creates a new service
func (c *RestClient) NewGetOrderHistoryService() *GetOrderHistoryService {
	return &GetOrderHistoryService{c: c}
}

// InstId sets instrument ID
func (s *GetOrderHistoryService) InstId(instId string) *GetOrderHistoryService {
	s.instId = instId
	return s
}

// OrderType sets order type
func (s *GetOrderHistoryService) OrderType(orderType string) *GetOrderHistoryService {
	s.orderType = orderType
	return s
}

// State sets order state
func (s *GetOrderHistoryService) State(state string) *GetOrderHistoryService {
	s.state = state
	return s
}

// After sets after parameter for pagination
func (s *GetOrderHistoryService) After(after string) *GetOrderHistoryService {
	s.after = after
	return s
}

// Before sets before parameter for pagination
func (s *GetOrderHistoryService) Before(before string) *GetOrderHistoryService {
	s.before = before
	return s
}

// Begin sets start time for filtering
func (s *GetOrderHistoryService) Begin(begin string) *GetOrderHistoryService {
	s.begin = begin
	return s
}

// End sets end time for filtering
func (s *GetOrderHistoryService) End(end string) *GetOrderHistoryService {
	s.end = end
	return s
}

// Limit sets limit results
func (s *GetOrderHistoryService) Limit(limit string) *GetOrderHistoryService {
	s.limit = limit
	return s
}

// Do executes request
func (s *GetOrderHistoryService) Do(ctx context.Context) (*GetOrderHistoryResponse, error) {
	if s.after != "" && s.before != "" {
		return nil, fmt.Errorf("after and before cannot be used simultaneously")
	}

	endpoint := "/api/v1/trade/orders-history"
	params := make(map[string]string)

	if s.instId != "" {
		params["instId"] = s.instId
	}
	if s.orderType != "" {
		params["orderType"] = s.orderType
	}
	if s.state != "" {
		params["state"] = s.state
	}
	if s.after != "" {
		params["after"] = s.after
	}
	if s.before != "" {
		params["before"] = s.before
	}
	if s.begin != "" {
		params["begin"] = s.begin
	}
	if s.end != "" {
		params["end"] = s.end
	}
	if s.limit != "" {
		params["limit"] = s.limit
	}

	req := s.c.httpClient.R().SetContext(ctx)
	if len(params) > 0 {
		req.SetQueryParams(params)
	}

	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetOrderHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTPSLOrderHistoryService service for getting TPSL order history
type GetTPSLOrderHistoryService struct {
	c             *RestClient
	instId        string
	tpslId        string
	clientOrderId string
	state         string
	after         string
	before        string
	limit         string
}

// NewGetTPSLOrderHistoryService creates a new service
func (c *RestClient) NewGetTPSLOrderHistoryService() *GetTPSLOrderHistoryService {
	return &GetTPSLOrderHistoryService{c: c}
}

// InstId sets instrument ID
func (s *GetTPSLOrderHistoryService) InstId(instId string) *GetTPSLOrderHistoryService {
	s.instId = instId
	return s
}

// TpslId sets TPSL order ID
func (s *GetTPSLOrderHistoryService) TpslId(tpslId string) *GetTPSLOrderHistoryService {
	s.tpslId = tpslId
	return s
}

// ClientOrderId sets client order ID
func (s *GetTPSLOrderHistoryService) ClientOrderId(clientOrderId string) *GetTPSLOrderHistoryService {
	s.clientOrderId = clientOrderId
	return s
}

// State sets order state
func (s *GetTPSLOrderHistoryService) State(state string) *GetTPSLOrderHistoryService {
	s.state = state
	return s
}

// After sets after parameter for pagination
func (s *GetTPSLOrderHistoryService) After(after string) *GetTPSLOrderHistoryService {
	s.after = after
	return s
}

// Before sets before parameter for pagination
func (s *GetTPSLOrderHistoryService) Before(before string) *GetTPSLOrderHistoryService {
	s.before = before
	return s
}

// Limit sets limit results
func (s *GetTPSLOrderHistoryService) Limit(limit string) *GetTPSLOrderHistoryService {
	s.limit = limit
	return s
}

// Do executes request
func (s *GetTPSLOrderHistoryService) Do(ctx context.Context) (*GetTPSLOrderHistoryResponse, error) {
	if s.after != "" && s.before != "" {
		return nil, fmt.Errorf("after and before cannot be used simultaneously")
	}

	endpoint := "/api/v1/trade/orders-tpsl-history"
	params := make(map[string]string)

	if s.instId != "" {
		params["instId"] = s.instId
	}
	if s.tpslId != "" {
		params["tpslId"] = s.tpslId
	}
	if s.clientOrderId != "" {
		params["clientOrderId"] = s.clientOrderId
	}
	if s.state != "" {
		params["state"] = s.state
	}
	if s.after != "" {
		params["after"] = s.after
	}
	if s.before != "" {
		params["before"] = s.before
	}
	if s.limit != "" {
		params["limit"] = s.limit
	}

	req := s.c.httpClient.R().SetContext(ctx)
	if len(params) > 0 {
		req.SetQueryParams(params)
	}

	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetTPSLOrderHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAlgoOrderHistoryService service for getting algo order history
type GetAlgoOrderHistoryService struct {
	c             *RestClient
	instId        string
	algoId        string
	clientOrderId string
	state         string
	after         string
	before        string
	limit         string
	orderType     string
}

// NewGetAlgoOrderHistoryService creates a new service
func (c *RestClient) NewGetAlgoOrderHistoryService() *GetAlgoOrderHistoryService {
	return &GetAlgoOrderHistoryService{c: c}
}

// InstId sets instrument ID
func (s *GetAlgoOrderHistoryService) InstId(instId string) *GetAlgoOrderHistoryService {
	s.instId = instId
	return s
}

// AlgoId sets algo order ID
func (s *GetAlgoOrderHistoryService) AlgoId(algoId string) *GetAlgoOrderHistoryService {
	s.algoId = algoId
	return s
}

// ClientOrderId sets client order ID
func (s *GetAlgoOrderHistoryService) ClientOrderId(clientOrderId string) *GetAlgoOrderHistoryService {
	s.clientOrderId = clientOrderId
	return s
}

// State sets order state
func (s *GetAlgoOrderHistoryService) State(state string) *GetAlgoOrderHistoryService {
	s.state = state
	return s
}

// After sets after parameter for pagination
func (s *GetAlgoOrderHistoryService) After(after string) *GetAlgoOrderHistoryService {
	s.after = after
	return s
}

// Before sets before parameter for pagination
func (s *GetAlgoOrderHistoryService) Before(before string) *GetAlgoOrderHistoryService {
	s.before = before
	return s
}

// Limit sets limit results
func (s *GetAlgoOrderHistoryService) Limit(limit string) *GetAlgoOrderHistoryService {
	s.limit = limit
	return s
}

// OrderType sets algo order type
func (s *GetAlgoOrderHistoryService) OrderType(orderType string) *GetAlgoOrderHistoryService {
	s.orderType = orderType
	return s
}

// Do executes request
func (s *GetAlgoOrderHistoryService) Do(ctx context.Context) (*GetAlgoOrderHistoryResponse, error) {
	if s.orderType == "" {
		return nil, fmt.Errorf("orderType required")
	}
	if s.after != "" && s.before != "" {
		return nil, fmt.Errorf("after and before cannot be used simultaneously")
	}

	endpoint := "/api/v1/trade/orders-algo-history"
	params := make(map[string]string)

	if s.instId != "" {
		params["instId"] = s.instId
	}
	if s.algoId != "" {
		params["algoId"] = s.algoId
	}
	if s.clientOrderId != "" {
		params["clientOrderId"] = s.clientOrderId
	}
	if s.state != "" {
		params["state"] = s.state
	}
	if s.after != "" {
		params["after"] = s.after
	}
	if s.before != "" {
		params["before"] = s.before
	}
	if s.limit != "" {
		params["limit"] = s.limit
	}
	if s.orderType != "" {
		params["orderType"] = s.orderType
	}

	req := s.c.httpClient.R().SetContext(ctx)
	if len(params) > 0 {
		req.SetQueryParams(params)
	}

	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetAlgoOrderHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTradeHistoryService service for getting trade history
type GetTradeHistoryService struct {
	c       *RestClient
	instId  string
	orderId string
	after   string
	before  string
	begin   string
	end     string
	limit   string
}

// NewGetTradeHistoryService creates a new service
func (c *RestClient) NewGetTradeHistoryService() *GetTradeHistoryService {
	return &GetTradeHistoryService{c: c}
}

// InstId sets instrument ID
func (s *GetTradeHistoryService) InstId(instId string) *GetTradeHistoryService {
	s.instId = instId
	return s
}

// OrderId sets order ID
func (s *GetTradeHistoryService) OrderId(orderId string) *GetTradeHistoryService {
	s.orderId = orderId
	return s
}

// After sets after parameter for pagination
func (s *GetTradeHistoryService) After(after string) *GetTradeHistoryService {
	s.after = after
	return s
}

// Before sets before parameter for pagination
func (s *GetTradeHistoryService) Before(before string) *GetTradeHistoryService {
	s.before = before
	return s
}

// Begin sets start time for filtering
func (s *GetTradeHistoryService) Begin(begin string) *GetTradeHistoryService {
	s.begin = begin
	return s
}

// End sets end time for filtering
func (s *GetTradeHistoryService) End(end string) *GetTradeHistoryService {
	s.end = end
	return s
}

// Limit sets limit results
func (s *GetTradeHistoryService) Limit(limit string) *GetTradeHistoryService {
	s.limit = limit
	return s
}

// Do executes request
func (s *GetTradeHistoryService) Do(ctx context.Context) (*GetTradeHistoryResponse, error) {
	if s.after != "" && s.before != "" {
		return nil, fmt.Errorf("after and before cannot be used simultaneously")
	}

	endpoint := "/api/v1/trade/fills-history"
	params := make(map[string]string)

	if s.instId != "" {
		params["instId"] = s.instId
	}
	if s.orderId != "" {
		params["orderId"] = s.orderId
	}
	if s.after != "" {
		params["after"] = s.after
	}
	if s.before != "" {
		params["before"] = s.before
	}
	if s.begin != "" {
		params["begin"] = s.begin
	}
	if s.end != "" {
		params["end"] = s.end
	}
	if s.limit != "" {
		params["limit"] = s.limit
	}

	req := s.c.httpClient.R().SetContext(ctx)
	if len(params) > 0 {
		req.SetQueryParams(params)
	}

	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetTradeHistoryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetOrderPriceRangeService represents the service for getting order price range
type GetOrderPriceRangeService struct {
	c      *RestClient
	instId string
	side   string
}

// NewGetOrderPriceRangeService creates a new GetOrderPriceRangeService
func (c *RestClient) NewGetOrderPriceRangeService() *GetOrderPriceRangeService {
	return &GetOrderPriceRangeService{c: c}
}

// InstId sets the instId parameter
func (s *GetOrderPriceRangeService) InstId(instId string) *GetOrderPriceRangeService {
	s.instId = instId
	return s
}

// Side sets the side parameter
func (s *GetOrderPriceRangeService) Side(side string) *GetOrderPriceRangeService {
	s.side = side
	return s
}

// Do sends the request
func (s *GetOrderPriceRangeService) Do(ctx context.Context) (*GetOrderPriceRangeResponse, error) {
	if s.instId == "" {
		return nil, errors.New("instId required")
	}
	if s.side == "" {
		return nil, errors.New("side required")
	}

	endpoint := "/api/v1/trade/order/price-range"
	params := map[string]string{
		"instId": s.instId,
		"side":   s.side,
	}

	req := s.c.httpClient.R().SetContext(ctx).SetQueryParams(params)
	s.c.addAuthHeaders(req, "GET", endpoint, "")
	resp, err := req.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result GetOrderPriceRangeResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
