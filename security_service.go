package finam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SecurityService запросить список инструментов
// Максимальное Количество запросов в минуту = 1
// https://finamweb.github.io/trade-api-docs/rest-api/securities
type SecurityService struct {
	c      *Client
	board  string
	symbol string
}

// Board set board
func (s *SecurityService) Board(board string) *SecurityService {
	s.board = board
	return s
}

// Symbol set symbol
func (s *SecurityService) Symbol(symbol string) *SecurityService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *SecurityService) Do(ctx context.Context) ([]Security, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "public/api/v1/securities/",
	}
	if s.board != "" {
		r.setParam("Board", s.board)
	}
	if s.symbol != "" {
		r.setParam("Seccode", s.symbol)
	}

	type Data struct {
		Sec Securities `json:"securities"`
	}
	type responseData struct {
		Error ResponseError `json:"error"`
		Data  Data          `json:"data"`
	}
	var rd responseData

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return rd.Data.Sec, err
	}
	if err = json.Unmarshal(data, &rd); err != nil {
		return rd.Data.Sec, err
	}
	// если есть ошибка вернем ее
	if rd.Error.Code != "" {
		return rd.Data.Sec, fmt.Errorf(rd.Error.Message)
	}

	return rd.Data.Sec, nil

}

// список инструментов
type Securities []Security //`json:"securities"`

// Инструмент
type Security struct {
	Code            string  `json:"code"`            // код инструмента;
	Board           string  `json:"board"`           // основной режим торгов инструмента;
	Market          string  `json:"market"`          // рынок инструмента. Тип Market;
	ShortName       string  `json:"shortName"`       // название инструмента;
	Ticker          string  `json:"ticker"`          //  тикер инструмента на биржевой площадке листинга;
	Decimals        int     `json:"decimals"`        // количество знаков в дробной части цены;
	LotSize         int     `json:"lotSize"`         //  размер лота;
	MinStep         float32 `json:"minStep"`         // минимальный шаг цены;
	Currency        string  `json:"currency"`        // код валюты номинала цены;
	Properties      int     `json:"properties"`      // параметры инструмента. Значение представлено в виде битовой маски:
	TimeZoneName    string  `json:"timeZoneName"`    // имя таймзоны;
	BpCost          float64 `json:"bpCost"`          //  стоимость пункта цены одного инструмента (не лота), без учета НКД;
	AccruedInterest float64 `json:"accruedInterest"` // текущий НКД;
	PriceSign       string  `json:"priceSign"`       // допустимая цена инструмента. Принимает следующие значения:
	LotDivider      int     `json:"lotDivider"`      // коэффициент дробления ценной бумаги в одном стандартном лоте.
}
