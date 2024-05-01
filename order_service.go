package finam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetOrderService получить список ордеров
type GetOrderService struct {
	c               *Client
	includeMatched  bool // Вернуть исполненные заявки;
	includeCanceled bool // Вернуть отмененные заявки;
	includeActive   bool // Вернуть активные заявки.
}

// IncludeMatched Вернуть исполненные заявки
func (s *GetOrderService) IncludeMatched(param bool) *GetOrderService {
	s.includeMatched = param
	return s
}

// IncludeCanceled Вернуть исполненные заявки
func (s *GetOrderService) IncludeCanceled(param bool) *GetOrderService {
	s.includeCanceled = param
	return s
}

// IncludeActive Вернуть активные заявки
func (s *GetOrderService) IncludeActive(param bool) *GetOrderService {
	s.includeActive = param
	return s
}

func (s *GetOrderService) Do(ctx context.Context) ([]Order, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "public/api/v1/orders/",
	}
	// параметры
	r.setParam("ClientId", s.c.clientId)
	if s.includeMatched {
		r.setParam("IncludeMatched", s.includeMatched)
	}
	if s.includeCanceled {
		r.setParam("IncludeCanceled", s.includeCanceled)
	}
	if s.includeActive {
		r.setParam("IncludeActive", s.includeActive)
	}

	type Data struct {
		Orders []Order `json:"orders"`
	}
	type responseData struct {
		Error ResponseError `json:"error"`
		Data  Data          `json:"data"`
	}
	var rd responseData

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return rd.Data.Orders, err
	}

	err = json.Unmarshal(data, &rd)
	if err != nil {
		return rd.Data.Orders, err
	}
	// если есть ошибка
	if rd.Error.Code != "" {
		return rd.Data.Orders, fmt.Errorf(rd.Error.Message)
	}

	return rd.Data.Orders, nil
}

// CancelOrderService удалить заявку по ее ID
type CancelOrderService struct {
	c             *Client
	transactionId int64
}

// Do sends the request.
func (s CancelOrderService) Do(ctx context.Context) error {
	r := &request{
		method:   http.MethodDelete,
		endpoint: "public/api/v1/orders/",
	}
	r.setParam("ClientId", s.c.clientId)
	r.setParam("TransactionId", s.transactionId)

	type responseData struct {
		Error ResponseError `json:"error"`
	}
	var rd responseData

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &rd)
	if err != nil {
		return err
	}
	// если есть ошибка
	if rd.Error.Code != "" {
		return fmt.Errorf(rd.Error.Message)
	}

	return nil
}

// CreateOrderService создать новую заявку (ордер)
type CreateOrderService struct {
	c     *Client
	order OrderRequest
}

// установим направление ордера
func (s *CreateOrderService) Side(side SideType) *CreateOrderService {
	s.order.BuySell = side
	return s
}

// Quantity установим объем заявки в лотах;
func (s *CreateOrderService) Quantity(quantity int32) *CreateOrderService {
	s.order.Quantity = quantity
	return s
}

// Price установим цену
// Для рыночной заявки указать значение null (или не передавать это поле).
// Для условной заявки необходимо указать цену исполнения;
func (s *CreateOrderService) Price(price float64) *CreateOrderService {
	s.order.Price = &price
	s.order.Condition.Price = price
	// TODO проверить правильность выставление Condition.Type
	if price == 0 {
		s.order.Price = nil
		s.order.Condition.Type = "Bid"
	}
	if price != 0 && s.order.BuySell == SideBuy {
		s.order.Condition.Type = "LastDown"
	}
	if price != 0 && s.order.BuySell == SideSell {
		s.order.Condition.Type = "LastUp"
	}
	//Bid - лучшая цена покупки;
	//BidOrLast- лучшая цена покупки или сделка по заданной цене и выше;
	//Ask - лучшая цена продажи;
	//AskOrLast - лучшая цена продажи или сделка по заданной цене и ниже;
	//Time - время выставления заявки на Биржу (параметр time должен быть установлен);
	//CovDown - обеспеченность ниже заданной;
	//CovUp - обеспеченность выше заданной;
	//LastUp - сделка на рынке по заданной цене или выше;
	//LastDown- сделка на рынке по заданной цене или ниже.
	return s
}

// TimeInForce установим условие по времени действия заявки
// TillEndSession - заявка действует до конца сессии;
// TillCancelled - заявка действует, пока не будет отменена;
// ExactTime - заявка действует до указанного времени. Параметр time должен быть задан (где его указать?)
func (s *CreateOrderService) TimeInForce(timeInForce TimeInForceType) *CreateOrderService {
	s.order.ValidBefore.Type = timeInForce.String()
	return s
}

// NewCreateOrderService init creating order service
// какие обязательные параметры?
func (c *Client) NewCreateOrderService(board, symbol string, sideType SideType, lot int32) *CreateOrderService {
	// условие по времени действия заявки
	// по умолчанию поставим = заявка действует до конца сессии;
	validBefore := &OrderValidBefore{
		Type: "TillEndSession",
	}
	// Свойства выставления заявок.
	condition := &OrderCondition{
		Type:  "Bid",
		Price: 1,
	}

	return &CreateOrderService{
		c: c,
		order: OrderRequest{
			ClientId:      c.clientId,
			SecurityCode:  symbol,
			SecurityBoard: board,
			BuySell:       sideType,
			Quantity:      lot,
			Property:      "PutInQueue", // неисполненная часть заявки помещается в очередь заявок биржи;
			UseCredit:     true,
			ValidBefore:   validBefore,
			Condition:     condition,
		},
	}
}

func (s *CreateOrderService) Do(ctx context.Context) (int64, error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "public/api/v1/orders/",
	}
	// в request.body надо записать s.order
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&s.order)
	r.body = buf

	type Data struct {
		SecurityCode  string `json:"securityCode"`
		TransactionId int64  `json:"transactionId"`
	}

	type responseData struct {
		Error ResponseError `json:"error"`
		Data  Data          `json:"data"`
	}
	var rd responseData

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(data, &rd)
	if err != nil {
		return 0, err
	}
	// если есть ошибка
	if rd.Error.Code != "" {
		return 0, fmt.Errorf(rd.Error.Message)
	}

	return rd.Data.TransactionId, nil
}

// Order структура ордера
type Order struct {
	OrderNo       int64            `json:"orderNo"`       // уникальный идентификатор заявки на бирже. Задается после того, как заявка будет принята биржей (см. поле status);
	TransactionId int64            `json:"transactionId"` // внутренний идентификатор заявки в системе TRANSAQ (для чужой заявки значение всегда равно 0);
	ClientId      string           `json:"clientId"`      // торговый код клиента;
	SecurityCode  string           `json:"securityCode"`  // код инструмента;
	SecurityBoard string           `json:"securityBoard"` // основной режим торгов инструмента;
	Market        string           `json:"market"`        // рынок инструмента. Тип Market.
	Status        string           `json:"status"`        // текущий статус заявки. Тип OrderStatus;
	BuySell       string           `json:"buySell"`       // тип BuySell ( SideType);
	CreatedAt     string           `json:"createdAt"`     // время регистрации заявки на бирже (UTC);
	Price         float64          `json:"price"`         // цена исполнения условной заявки. Для рыночной заявки значение всегда равно 0;
	Quantity      int              `json:"quantity"`      // объем заявки в лотах;
	Balance       int              `json:"balance"`       // неисполненный остаток, в лотах. Изначально равен quantity, но по мере исполнения заявки (совершения сделок) будет уменьшаться на объем сделки. Значение 0 будет соответствовать полностью исполненной заявке (см. поле status);
	Message       string           `json:"message"`       // содержит сообщение об ошибке, возникшей при обработке заявки. Заявка может быть отклонена по разным причинам сервером TRANSAQ или биржей с выставлением поля status;
	Currency      string           `json:"currency"`      // код валюты цены
	AcceptedAt    string           `json:"acceptedAt"`    // время регистрации заявки на сервере TRANSAQ (UTC);
	Condition     *OrderCondition  `json:"condition"`     // может быть null/ свойства выставления заявок. Тип OrderCondition;
	ValidBefore   OrderValidBefore `json:"validBefore"`   // условие по времени действия заявки. Тип OrderValidBefore;

}

// OrderRequest Запрос на создание заявки.
type OrderRequest struct {
	ClientId      string            `json:"clientId,omitempty"`      // Идентификатор торгового счёта.
	SecurityBoard string            `json:"securityBoard,omitempty"` // Trading Board. Режим торгов.
	SecurityCode  string            `json:"securityCode,omitempty"`  // Security Code. Тикер инструмента.
	BuySell       SideType          `json:"buySell,omitempty"`       // Направление сделки.
	Quantity      int32             `json:"quantity,omitempty"`      // Объем заявки в лотах;
	UseCredit     bool              `json:"useCredit,omitempty"`     // Использовать кредит. Недоступно для срочного рынка.
	Price         *float64          `json:"price"`                   // Цена заявки. Используйте "null", чтобы выставить рыночную заявку.
	Property      string            `json:"property,omitempty"`      // Свойства исполнения частично исполненных заявок.
	Condition     *OrderCondition   `json:"condition,omitempty"`     // Свойства выставления заявок.
	ValidBefore   *OrderValidBefore `json:"validBefore,omitempty"`   // Условие по времени действия заявки.
}
