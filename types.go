package finam

// условие по времени действия заявки
// TillEndSession - заявка действует до конца сессии;
// TillCancelled - заявка действует, пока не будет отменена;
// ExactTime - заявка действует до указанного времени. Параметр time должен быть задан.
type TimeInForceType string

const (
	TimeInForceGTC       TimeInForceType = "TillCancelled"  //"GTC" // Ордер будет находится в очереди до тех пор, пока не будет снят
	TimeInForceDAY       TimeInForceType = "TillEndSession" // "DAY" // Ордер будет действовать только в течение текущего торгового дня (до конца сессии)
	TimeInForceExactTime TimeInForceType = "ExactTime"      // заявка действует до указанного времени. Параметр time должен быть задан.
)

func (t TimeInForceType) String() string {
	return string(t)
}

type OrderValidBefore struct {
	Type string  `json:"type"`
	Time *string `json:"time,omitempty"` // может быть null
}

// тип SideType (BuySell)  Определяет тип операции: покупка или продажа.
type SideType string

const (
	SideBuy  SideType = "Buy"  // покупка,
	SideSell SideType = "Sell" // продажа
)

func (side SideType) String() string {
	return string(side)
}

// Тип OrderStatus Статус заявки.
type OrderStatus string

const (
	OrderStatusNone      OrderStatus = "None"      // принята сервером TRANSAQ, и заявке присвоен transactionId;
	OrderStatusActive    OrderStatus = "Active"    // принята биржей, и заявке присвоен orderNo;
	OrderStatusMatched   OrderStatus = "Matched"   // полностью исполнилась (выполнилась);
	OrderStatusCancelled OrderStatus = "Cancelled" // была отменена (снята) пользователем или биржей.
)

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
type OrderCondition struct {
	Type  string  `json:"type"`
	Price float64 `json:"price"`
	Time  string  `json:"time,omitempty"`
}

// property - свойства исполнения частично исполненных заявок. Принимает следующие значения:
// PutInQueue    - неисполненная часть заявки помещается в очередь заявок биржи;
// CancelBalance - неисполненная часть заявки снимается с торгов;
// ImmOrCancel   - сделки совершаются только в том случае, если заявка может быть удовлетворена полностью и сразу при выставлении.
