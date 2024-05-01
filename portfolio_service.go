package finam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetPortfolioService получим информацию по портфелю
// https://finamweb.github.io/trade-api-docs/rest-api/portfolios
// clientId - торговый код клиента (обязательный);
// includeCurrencies - запросить информацию по валютам портфеля;
// includeMoney - запросить информацию по денежным позициям портфеля;
// includePositions - запросить информацию по позициям портфеля;
// includeMaxBuySell - запросить информацию о максимальном доступном объеме на покупку/продажу.
type GetPortfolioService struct {
	c                 *Client
	includeCurrencies bool // запросить информацию по валютам портфеля;
	includeMoney      bool // запросить информацию по денежным позициям портфеля;
	includePositions  bool // запросить информацию по позициям портфеля;
	includeMaxBuySell bool // запросить информацию о максимальном доступном объеме на покупку/продажу.
}

// IncludeCurrencies запросить информацию по валютам портфеля
func (s *GetPortfolioService) IncludeCurrencies(param bool) *GetPortfolioService {
	s.includeCurrencies = param
	return s
}

// IncludeMoney запросить информацию по денежным позициям портфеля
func (s *GetPortfolioService) IncludeMoney(param bool) *GetPortfolioService {
	s.includeMoney = param
	return s
}

// IncludePositions запросить информацию по позициям портфеля
func (s *GetPortfolioService) IncludePositions(param bool) *GetPortfolioService {
	s.includePositions = param
	return s
}

// IncludeMaxBuySell запросить информацию о максимальном доступном объеме на покупку/продажу
func (s *GetPortfolioService) IncludeMaxBuySell(param bool) *GetPortfolioService {
	s.includeMaxBuySell = param
	return s
}

// https://trade-api.finam.ru/public/api/v1/portfolio?Content.IncludeCurrencies=true&Content.IncludeMoney=true&Content.IncludePositions=true&Content.IncludeMaxBuySell=true
// Do sends the request.
func (s *GetPortfolioService) Do(ctx context.Context) (Portfolio, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "public/api/v1/portfolio",
	}
	// параметры
	r.setParam("ClientId", s.c.clientId)
	if s.includeCurrencies {
		r.setParam("Content.IncludeCurrencies", s.includeCurrencies)
	}
	if s.includeMoney {
		r.setParam("Content.IncludeMoney", s.includeMoney)
	}
	if s.includePositions {
		r.setParam("Content.IncludePositions", s.includePositions)
	}
	if s.includeMaxBuySell {
		r.setParam("Content.IncludeMaxBuySell", s.includeMaxBuySell)
	}

	type responseData struct {
		Error ResponseError `json:"error"`
		Data  Portfolio     `json:"data"`
	}
	var rd responseData

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return rd.Data, err
	}

	err = json.Unmarshal(data, &rd)
	if err != nil {
		return rd.Data, err
	}
	// если есть ошибка
	if rd.Error.Code != "" {
		return rd.Data, fmt.Errorf(rd.Error.Message)
	}
	return rd.Data, nil
}

// позиции портфеля
type Position struct {
	SecurityCode         string  `json:"securityCode"`      // код инструмента
	Market               string  `json:"market"`            //  рынок инструмента. Тип Market;
	Balance              float64 `json:"balance"`           // текущая позиция;
	CurrentPrice         float64 `json:"currentPrice"`      // текущая цена в валюте инструмента;
	Equity               float64 `json:"equity"`            //  текущая оценка инструмента;
	AveragePrice         float64 `json:"averagePrice"`      // средняя цена;
	Currency             string  `json:"currency"`          // код валюты риска;
	AccumulatedProfit    float64 `json:"accumulatedProfit"` // прибыль/убыток по входящим;
	TodayProfit          float64 `json:"todayProfit"`       // прибыль/убыток по сделкам;
	UnrealizedProfit     float64 `json:"unrealizedProfit"`  //  нереализованная прибыль/убыток;
	Profit               float64 `json:"profit"`            //  прибыль/убыток;
	MaxBuy               float64 `json:"maxBuy"`            //  максимально возможное количество лотов на покупку/продажу (вычисляется, если указать флаг includeMaxBuySell в true, иначе значение будет равно 0);
	MaxSell              float64 `json:"maxSell"`
	PriceCurrency        string  `json:"priceCurrency"`        // priceCurrency
	AverageRate          float64 `json:"averageRate"`          // код валюты балансовой цены;
	AveragePriceCurrency string  `json:"averagePriceCurrency"` // кросс-курс валюты балансовой цены к валюте риска.
}

// валюта портфеля
type Сurrency struct {
	Name             string  `json:"name"`             // код валюты;
	Equity           float64 `json:"equity"`           // оценка позиции;
	Balance          float64 `json:"balance"`          // текущая позиция;
	CrossRate        float64 `json:"crossRate"`        // курс валюты;
	UnrealizedProfit float64 `json:"unrealizedProfit"` // нереализованная прибыль/убыток.
}

// денежные позиции
type Money struct {
	Market   string  `json:"market"`   // рынок. Тип Market;
	Currency string  `json:"currency"` //  код валюты;
	Balance  float64 `json:"balance"`  // текущая позиция.
}

// структура Портфеля
type Portfolio struct {
	ClientId   string     `json:"clientId"`   // торговый код клиента;
	Equity     float64    `json:"equity"`     // текущая оценка портфеля;
	Balance    float64    `json:"balance"`    // входящая оценка стоимости портфеля;
	Positions  []Position `json:"positions"`  // позиции портфеля
	Currencies []Сurrency `json:"currencies"` //  валюта портфеля
	Money      []Money    `json:"money"`      // денежные позиции
}
