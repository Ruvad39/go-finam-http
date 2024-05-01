package finam

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"time"
)

// получить исторические свечи
// https://finamweb.github.io/trade-api-docs/rest-api/candles
// Обязательные параметры securityCode, securityBoard и timeFrame.
// Запросить можно как определенное количество свечей, так и за интервал.
// Для запроса количества свечей в запросе необходимо указать count и либо from (начиная с указанной даты) либо to (до указанной даты).
// Для запроса за интервал необходимо указать from и to.
// Запрос дневных/недельных свечей
// Максимальный интервал: 365 дней
// Максимальное кол-во запросов в минуту: 120
// Дата начала (окончания) в формате yyyy-MM-dd в часовом поясе UTC
// Запрос внутридневных свечей
// Максимальный интервал: 30 дней
// Максимальное кол-во запросов в минуту: 120
// Дата начала (окончания) в формате yyyy-MM-ddTHH:mm:ssZ в часовом поясе UTC
type CandlesService struct {
	c         *Client
	board     string
	symbol    string
	timeFrame TimeFrame
	count     int
	startTime *time.Time
	endTime   *time.Time
}

// Limit set limit
func (s *CandlesService) Count(param int) *CandlesService {
	s.count = param
	return s
}

// StartTime set startTime
func (s *CandlesService) StartTime(startTime time.Time) *CandlesService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *CandlesService) EndTime(endTime time.Time) *CandlesService {
	s.endTime = &endTime
	return s
}

// Do send request
func (s *CandlesService) Do(ctx context.Context) ([]Candle, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "public/api/v1/day-candles",
	}
	// если это внутредневная свеча = другой запрос
	if !s.timeFrame.IsDay() {
		r.endpoint = "public/api/v1/intraday-candles"
	}
	r.setParam("SecurityBoard", s.board)
	r.setParam("SecurityCode", s.symbol)
	r.setParam("TimeFrame", s.timeFrame)

	if s.count != 0 {
		r.setParam("Interval.Count", s.count)
	}
	// разный формат дат в зависимости от того дневная или нет
	// дневная yyyy-MM-dd в часовом поясе UTC
	layout := "2006-01-02"
	// внутредневная в формате yyyy-MM-ddTHH:mm:ssZ в часовом поясе UTC
	// 2006-01-02T15:04:05-0700

	if !s.timeFrame.IsDay() {
		layout = "2006-01-02T15:04:05-0700"
	}
	// и надо перевети в часовой пояс UTC
	if s.startTime != nil {
		t := *s.startTime
		r.setParam("Interval.From", t.UTC().Format(layout))
	}
	if s.endTime != nil {
		t := *s.endTime
		r.setParam("Interval.To", t.UTC().Format(layout))
	}

	type Data struct {
		Candles Candles `json:"candles"`
	}
	type responseData struct {
		Error ResponseError `json:"error"`
		Data  Data          `json:"data"`
	}
	var rd responseData

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return rd.Data.Candles, err
	}
	if err = json.Unmarshal(data, &rd); err != nil {
		return rd.Data.Candles, err
	}
	// если есть ошибка вернем ее
	if rd.Error.Code != "" {
		return rd.Data.Candles, fmt.Errorf(rd.Error.Message)
	}

	return rd.Data.Candles, nil
}

// структуры данных для свечей

// TimeFrame период свечей
type TimeFrame string

const (
	TimeFrame_M1  TimeFrame = "M1"
	TimeFrame_M5  TimeFrame = "M5"
	TimeFrame_M15 TimeFrame = "M15"
	TimeFrame_H1  TimeFrame = "H1"
	TimeFrame_D1  TimeFrame = "D1"
	TimeFrame_W1  TimeFrame = "W1"
)

// IsDay вернем признак дневных свечей
func (t TimeFrame) IsDay() bool {
	switch t {
	case TimeFrame_D1:
		return true
	case TimeFrame_W1:
		return true
	default:
		return false
	}
}

// Decimal Представляет десятичное число с плавающей запятой:
// Итоговое значение вычисляется по формуле: num * 10^(-scale). Где ^ оператор возведение в степень.
type Decimal struct {
	Num   int `json:"num"`   //  мантисса;
	Scale int `json:"scale"` //  экспонента по основанию 10.
}

type Candles []Candle //`json:"candles"`

// Candle структура свечи
type Candle struct {
	Date      string  `json:"date"`      // дневная свеча дата свечи в формате yyyy-MM-dd (в локальном времени биржи);
	Timestamp string  `json:"timestamp"` // внутридневная свеча дата и время свечи в формате yyyy-MM-ddTHH:mm:ssZ в поясе UTC;
	Open      Decimal `json:"open"`      // цена открытия (тип Decimal);
	Close     Decimal `json:"close"`     //  цена закрытия (тип Decimal);
	High      Decimal `json:"high"`      //  максимальная цена (тип Decimal);
	Low       Decimal `json:"low"`       //  минимальная цена (тип Decimal);
	Volume    int64   `json:"volume"`    //  объем торгов.
}

// расчитаем и вернем цену
func calcPrice(dec Decimal) float64 {
	return float64(dec.Num) * math.Pow(10, -float64(dec.Scale))
}

// GetClose расчитаем и вернем цену закрытия
func (k *Candle) GetClose() float64 {
	return calcPrice(k.Close)
}

func (k *Candle) GetOpen() float64 {
	return calcPrice(k.Open)
}

func (k *Candle) GetHigh() float64 {
	return calcPrice(k.High)
}

func (k *Candle) GetLow() float64 {
	return calcPrice(k.Low)
}

func (k *Candle) GetVolume() int64 {
	return k.Volume
}

// как определить в каком поле сидит дата-время?
// yyyy-MM-ddTHH:mm:ssZ
var layout = "2006-01-02T15:04:05Z"

func (k *Candle) GetDateTime() string {
	if k.Date != "" {
		return k.Date
	}

	// в часовом поясе UTC. перевести в Moscow ???
	//time.LoadLocation()
	t, err := time.Parse(layout, k.Timestamp)
	if err != nil {
		slog.Error("GetTime", "err", err.Error())
	}
	//t2 := t.In(Moscow)
	//return t2.String()
	return t.In(Moscow).String()
}

func (k *Candle) GetDateTimeToTime() time.Time {
	// если дневной тайм-фрейм = должны показать 00 часов
	// если внутредневной = время уже по Москве
	var t time.Time
	datetime := k.Timestamp
	layout := "2006-01-02T15:04:05Z"

	// дневная свеча
	if k.Date != "" {
		datetime = k.Date
		layout = "2006-01-02"
		loc, _ := time.LoadLocation("Europe/Moscow")
		t, err := time.ParseInLocation(layout, datetime, loc)
		if err != nil {
			slog.Error("GetTime", "err", err.Error())
		}
		return t
	}

	// внутредневная свеча
	t, err := time.Parse(layout, datetime)
	if err != nil {
		slog.Error("GetTime", "err", err.Error())
	}
	// в часовом поясе UTC. перевести в Moscow ???
	return t.In(Moscow)

}

func (k *Candle) String() string {
	str := fmt.Sprintf("DateTime:%v ,O:%v, H:%v, L:%v, C:%v, V:%v",
		k.GetDateTime(), k.GetOpen(), k.GetHigh(), k.GetLow(), k.GetClose(), k.GetVolume(),
	)
	return str
}

var Moscow = initMoscow()

func initMoscow() *time.Location {
	var loc, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.FixedZone("MSK", int(3*time.Hour/time.Second))
	}
	return loc
}
