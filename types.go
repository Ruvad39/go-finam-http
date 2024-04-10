package finam

import (
	"fmt"
	"log/slog"
	"math"
	"time"
)

// структура ощибки
type ResponseError struct {
	Code      string  `json:"code"`    
	Message   string  `json:"message"`    	
	Data      string  `json:"data"`    	
}

// позиции портфеля 
type Position struct {
	SecurityCode      string  `json:"securityCode"`     // код инструмента
	Market            string  `json:"market"`    	    //  рынок инструмента. Тип Market;
	Balance           float64 `json:"balance"`          // текущая позиция;
	CurrentPrice      float64 `json:"currentPrice"`  	 // текущая цена в валюте инструмента;
	Equity            float64 `json:"equity"`  	         //  текущая оценка инструмента;
	AveragePrice      float64 `json:"averagePrice"`  	 // средняя цена;
	Currency          string  `json:"currency"`  	     // код валюты риска;
	AccumulatedProfit float64 `json:"accumulatedProfit"` // прибыль/убыток по входящим; 		
	TodayProfit       float64 `json:"todayProfit"`  	 // прибыль/убыток по сделкам;	
	UnrealizedProfit  float64 `json:"unrealizedProfit"`  //  нереализованная прибыль/убыток;		
	Profit 			  float64 `json:"profit"`  			 //  прибыль/убыток;
	MaxBuy            float64 `json:"maxBuy"`  			 //  максимально возможное количество лотов на покупку/продажу (вычисляется, если указать флаг includeMaxBuySell в true, иначе значение будет равно 0);	
	MaxSell           float64 `json:"maxSell"`  					
	PriceCurrency     string  `json:"priceCurrency"`  		 // priceCurrency				
	AverageRate       float64 `json:"averageRate"`  		 // код валюты балансовой цены;			
	AveragePriceCurrency string `json:"averagePriceCurrency"`// кросс-курс валюты балансовой цены к валюте риска. 							

}

// валюта портфеля
type Сurrency struct {
  Name             string   `json:"name"`             // код валюты;
  Equity           float64  `json:"equity"`           // оценка позиции;
  Balance          float64  `json:"balance"`          // текущая позиция;  
  CrossRate        float64  `json:"crossRate"`        // курс валюты;
  UnrealizedProfit float64  `json:"unrealizedProfit"` // нереализованная прибыль/убыток.

}

// денежные позиции
type Money struct {
  Market           string   `json:"market"`     // рынок. Тип Market;
  Currency         string  `json:"currency"`    //  код валюты; 
  Balance          float64  `json:"balance"`    // текущая позиция. 

}

// структура Портфеля
type Portfolio struct {
  ClientId    string   `json:"clientId"`     // торговый код клиента; 
  Equity      float64  `json:"equity"`       // текущая оценка портфеля;
  Balance     float64  `json:"balance"`      // входящая оценка стоимости портфеля; 
  Positions   []Position `json:"positions"`   // позиции портфеля 
  Currencies  []Сurrency `json:"currencies"` //  валюта портфеля
  Money       []Money    `json:"money"`      // денежные позиции
  	// content
}

// список инструментов
type Securities []Security //`json:"securities"`

// Инструмент
type Security struct {
 	Code           string   `json:"code"`     // код инструмента;
	Board           string  `json:"board"`    // основной режим торгов инструмента;
	Market          string  `json:"market"`   // рынок инструмента. Тип Market;
	ShortName       string  `json:"shortName"`    // название инструмента;
	Ticker          string  `json:"ticker"`    //  тикер инструмента на биржевой площадке листинга;
	Decimals        int     `json:"decimals"`    // количество знаков в дробной части цены;
	LotSize         int     `json:"lotSize"`    //  размер лота;
	MinStep         float32 `json:"minStep"`    // минимальный шаг цены;
	Currency        string  `json:"currency"`    // код валюты номинала цены;
	Properties      int     `json:"properties"`    // параметры инструмента. Значение представлено в виде битовой маски:
	TimeZoneName    string  `json:"timeZoneName"` // имя таймзоны;
	BpCost          float64 `json:"bpCost"`    //  стоимость пункта цены одного инструмента (не лота), без учета НКД;
	AccruedInterest float64 `json:"accruedInterest"`    // текущий НКД;
	PriceSign       string  `json:"priceSign"`    // допустимая цена инструмента. Принимает следующие значения:
	LotDivider      int     `json:"lotDivider"`    // коэффициент дробления ценной бумаги в одном стандартном лоте.
}
/*
priceSign - допустимая цена инструмента. Принимает следующие значения:
Positive - положительная,
NonNegative - неотрицательная,
Any - любая.

properties - параметры инструмента. Значение представлено в виде битовой маски:

0 - нет параметров;
1 - инструмент торгуется на бирже;
2 - инструмент допущен к торгам у брокера - существенно для НЕ ГЛАВНЫХ трейдеров. Главным доступны все инструменты, торгуемые на биржах;
4 - рыночные заявки (без ограничения по цене) разрешены;
8 - признак маржинальности бумаги;
16 - опцион Call;
32 - опцион Put;
48 - фьючерс Call | Put;
64 - разрешен для резидентов;
128 - разрешен для нерезидентов.
*/

// Представляет десятичное число с плавающей запятой:
// Итоговое значение вычисляется по формуле: num * 10^(-scale). Где ^ оператор возведение в степень.
type Decimal struct {
	Num int `json:"num"`     //  мантисса;
	Scale int `json:"scale"` //  экспонента по основанию 10.
}

// список свечей
type Candles []Candle //`json:"candles"`

// структура свечи 
// `json:"symbol" yaml:"symbol"`
type Candle struct {
		//Time      string  `json:"date" json:"timestamp"` 
		Date      string  `json:"date"`      // дневная свеча дата свечи в формате yyyy-MM-dd (в локальном времени биржи);
		Timestamp string  `json:"timestamp"` // внутридневная свеча дата и время свечи в формате yyyy-MM-ddTHH:mm:ssZ в поясе UTC;
		Open      Decimal `json:"open"`      // цена открытия (тип Decimal);
		Close     Decimal `json:"close"`     //  цена закрытия (тип Decimal);
		High      Decimal `json:"high"`      //  максимальная цена (тип Decimal);
  	Low       Decimal `json:"low"`       //  минимальная цена (тип Decimal);
  	Volume    int64   `json:"volume"`    //  объем торгов.
}

// расчитаем и вернем цену
func calcPrice(dec Decimal) float64{
	return float64(dec.Num) * math.Pow(10, -float64(dec.Scale))
}

// расчитаем и вернем цену закрытия
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
	if k.Date !=""{
		return k.Date
	}
	
	// в часовом поясе UTC. перевести в Moscow ???
	//time.LoadLocation()
	time, err := time.Parse(layout, k.Timestamp)
	if err != nil {
				slog.Error("GetTime", "err", err.Error())
			}
	t2 := time.In(Moscow)
	return t2.String()
	//return k.Timestamp
}

func (k Candle) String() string {
    str := fmt.Sprintf("DateTime:%v ,O:%v, H:%v, L:%v, C:%v, V:%v", 
    	k.GetDateTime(), k.GetOpen(), k.GetHigh(), k.GetLow(), k.GetClose(), k.GetVolume(),
    )
    return str
}

// TODO сейчес не правильно идет перевод дневной свечи. она уже в 
// дата свечи в формате yyyy-MM-dd (в локальном времени биржи);
// timestamp - дата и время свечи в формате yyyy-MM-ddTHH:mm:ssZ в поясе UTC;
func (k *Candle) GetDateTimeToTime() time.Time {
	var datetime  = k.Timestamp
	var layout    = "2006-01-02T15:04:05Z"

	if k.Date !=""{
		datetime = k.Date
		layout = "2006-01-02"
	}
	
	// в часовом поясе UTC. перевести в Moscow ???
	time, err := time.Parse(layout, datetime)
	if err != nil {
				slog.Error("GetTime", "err", err.Error())
			}
	//t2 := time.In(Moscow)
	return time.In(Moscow)
	//return k.Timestamp
}

// M1 - 1 минута
// M5 - 5 минут
// M15 - 15 минут
// H1 - 1 час
// Доступные таймфреймы для дневных свечей:
// D1 - 1 день
// W1 - 1 неделя
// TimeFrame период свечей
type TimeFrame string

// пока такой список: при необходимости добавлю
var TimeFrame_M1  = TimeFrame("M1")
var TimeFrame_M5  = TimeFrame("M5")
var TimeFrame_M15  = TimeFrame("M15")
var TimeFrame_H1  = TimeFrame("H1")
var TimeFrame_D1  = TimeFrame("D1")
var TimeFrame_W1  = TimeFrame("W1")

func (t TimeFrame) String() string {
    return string(t)
}

// вернем признак дневных свечей
// иначе внутридневных
func (t TimeFrame) IsDay() bool {
		switch t {
		case TimeFrame_D1:
				return true
		case TimeFrame_W1:
				return true
		default:
				return false
		}
    //return false
}


var Moscow = initMoscow()
// для обработки даты с квика
func initMoscow() *time.Location {
    var loc, err = time.LoadLocation("Europe/Moscow")
    if err != nil {
        loc = time.FixedZone("MSK", int(3*time.Hour/time.Second))
    }
    return loc
}


// условие по времени действия заявки
// TillEndSession - заявка действует до конца сессии;
// TillCancelled - заявка действует, пока не будет отменена;
// ExactTime - заявка действует до указанного времени. Параметр time должен быть задан.
type OrderValidBefore struct{
	Type string 	`json:"type"`
	Time *string     `json:"time,omitempty"` // может быть null
}
	

// тип BuySell
// Определяет тип операции: покупка или продажа. Принимает следующие значения:
// Buy - покупка,
// Sell- продажа.
type SideType string

const (
	SideBuy  = SideType("Buy")
	SideSell = SideType("Sell")
)
func (side SideType) String() string {
	return string(side)
}

// Тип OrderStatus Статус заявки. Принимает следующие значения:
// None - заявка принята сервером TRANSAQ, и заявке присвоен transactionId;
// Active - заявка принята биржей, и заявке присвоен orderNo;
// Matched - заявка полностью исполнилась (выполнилась);
// Cancelled - заявка была отменена (снята) пользователем или биржей.
type OrderStatus string

const (
	OrderStatusNone  = OrderStatus("None")
	OrderStatusActive = OrderStatus("Active")
	OrderStatusMatched = OrderStatus("Matched")
	OrderStatusCancelled = OrderStatus("Cancelled")
)
func (status OrderStatus) String() string {
	return string(status)
}


// Свойства выставления заявок. 
// Тип условия определяет значение поля type, которое принимает следующие значения:
// Bid - лучшая цена покупки;
// BidOrLast- лучшая цена покупки или сделка по заданной цене и выше;
// Ask - лучшая цена продажи;
// AskOrLast - лучшая цена продажи или сделка по заданной цене и ниже;
// Time - время выставления заявки на Биржу (параметр time должен быть установлен);
// CovDown - обеспеченность ниже заданной;
// CovUp - обеспеченность выше заданной;
// LastUp - сделка на рынке по заданной цене или выше;
// LastDown- сделка на рынке по заданной цене или ниже.
type OrderCondition struct{
	Type  string  `json:"type"`
	//Price float64 `json:"price,omitempty"`   
	Price float64 `json:"price"`   
	Time  string  `json:"time,omitempty"` 
}

// property - свойства исполнения частично исполненных заявок. Принимает следующие значения:
// PutInQueue    - неисполненная часть заявки помещается в очередь заявок биржи;
// CancelBalance - неисполненная часть заявки снимается с торгов;
// ImmOrCancel   - сделки совершаются только в том случае, если заявка может быть удовлетворена полностью и сразу при выставлении.

// order
type Order struct {
 	OrderNo       int64 	`json:"orderNo"`      //  уникальный идентификатор заявки на бирже. Задается после того, как заявка будет принята биржей (см. поле status); 
 	TransactionId int64 	`json:"transactionId"` // внутренний идентификатор заявки в системе TRANSAQ (для чужой заявки значение всегда равно 0);
 	ClientId      string 	`json:"clientId"`     //  торговый код клиента;
 	SecurityCode  string 	`json:"securityCode"` //код инструмента;
 	SecurityBoard string 	`json:"securityBoard"` //основной режим торгов инструмента;
	Market        string 	`json:"market"`                // рынок инструмента. Тип Market.	
 	Status        string 	`json:"status"`       // текущий статус заявки. Тип OrderStatus;
 	BuySell       string 	`json:"buySell"`      //  тип BuySell;
 	CreatedAt     string 	`json:"createdAt"`    // время регистрации заявки на бирже (UTC);
	Price         float64 	`json:"price"`        // цена исполнения условной заявки. Для рыночной заявки значение всегда равно 0;
	Quantity      int 	`json:"quantity"` //  объем заявки в лотах;
	Balance       int 	`json:"balance"` //  неисполненный остаток, в лотах. Изначально равен quantity, но по мере исполнения заявки (совершения сделок) будет уменьшаться на объем сделки. Значение 0 будет соответствовать полностью исполненной заявке (см. поле status);
	Message       string 	`json:"message"` //  содержит сообщение об ошибке, возникшей при обработке заявки. Заявка может быть отклонена по разным причинам сервером TRANSAQ или биржей с выставлением поля status;
	Currency      string 	`json:"currency"`
	AcceptedAt    string 	`json:"acceptedAt"`    // время регистрации заявки на сервере TRANSAQ (UTC);
	Condition     *OrderCondition   `json:"condition"`  // может быть null/ свойства выставления заявок. Тип OrderCondition;
	ValidBefore   OrderValidBefore `json:"validBefore"` // условие по времени действия заявки. Тип OrderValidBefore;

}





// New Order Request.
// Запрос на создание заявки.
type NewOrderRequest struct {
	// Trade Account ID. Идентификатор торгового счёта.
	ClientId string `json:"clientId,omitempty"`
	// Trading Board. Режим торгов.
	SecurityBoard string `json:"securityBoard,omitempty"`
	// Security Code. Тикер инструмента.
	SecurityCode string `"json:"securityCode,omitempty"`
	// Направление сделки.
	//BuySell SideType `json:"buySell,omitempty"`
	BuySell string `json:"buySell,omitempty"`
	// Количество лотов инструмента для заявки.
	Quantity int32 `json:"quantity,omitempty"`
	// Использовать кредит. Недоступно для срочного рынка.
	UseCredit bool `json:"useCredit,omitempty"`
	// Order price. Use "null" to place Market Order.
	// Цена заявки. Используйте "null", чтобы выставить рыночную заявку.
	//Price *float64 `json:"price,omitempty"`
	Price *float64 `json:"price"`
	// Unfilled order execution property.
	// Свойства исполнения частично исполненных заявок.
	//Property OrderProperty `json:"property,omitempty"`
	Property string `json:"property,omitempty"`
	// Свойства выставления заявок.
	Condition *OrderCondition `json:"condition,omitempty"`
	// Order lifetime condition.
	// Условие по времени действия заявки.
	ValidBefore *OrderValidBefore `protobuf:"json:"validBefore,omitempty"`
}
