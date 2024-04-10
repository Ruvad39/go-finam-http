package finam

import (
	"context"
	"encoding/json"
	"fmt"
	//"log/slog"
	"net/url"
	"path"
	"strconv"
)

// получить список ордеров
// clientId - торговый код клиента (обязательный);
// includeMatched - вернуть исполненные заявки;
// includeCanceled - вернуть отмененные заявки;
// includeActive - вернуть активные заявки.
func (client *Client) GetOrders(ctx context.Context, opts ...Option) ( []Order, error){
  endPoint := "public//api/v1/orders/"
  metod    := "GET"

  type Data struct{
    Orders []Order `json:"orders"`
  }
  type responseData struct {
    Error  ResponseError  `json:"error"`
    Data   Data          `json:"data"`
  }
  var rd responseData

  p := &Options{
    IncludeActive: true , 
  }
  // обработаем входящие параметры
  for _, opt := range opts {
      opt(p)
  }

  url, err := url.Parse(baseURL)
  if err != nil {
    return rd.Data.Orders, err
  }
  url.Path = path.Join(url.Path, endPoint)

  // создаем параметры
  q := url.Query()
  q.Set("ClientId", client.clientId)

  // includeMatched - вернуть исполненные заявки;
  // includeCanceled - вернуть отмененные заявки;
  // includeActive - вернуть активные заявки.

  if p.IncludeMatched{
    q.Set("IncludeMatched", "true")
  }
  if p.IncludeCanceled{
    q.Set("IncludeCanceled", "true")
  }
  if p.IncludeActive{
    q.Set("IncludeActive", "true")
  }

  // добавляем к URL параметры
  url.RawQuery = q.Encode()

  resp, err := client.GetHttp(ctx, metod, url.String(), nil)
  if err != nil {
    return rd.Data.Orders, err 
  }

  if err = json.Unmarshal(resp, &rd); err != nil {
    client.Logger.Error("GetOrders Ошибка при разборе ответа JSON", "err", err.Error())
    return rd.Data.Orders, err
  }

  // если есть ошибка 
  if rd.Error.Code != ""{
    return rd.Data.Orders, fmt.Errorf(rd.Error.Message)
  }

  return rd.Data.Orders, nil
}

// удаление заявки
// clientId - торговый код клиента (обязательный)
// transactionId int64 (обязательный)

func (client *Client) DeleteOrder(ctx context.Context, transactionId int64) error{
  endPoint := "public//api/v1/orders/"
  metod    := "DELETE"

  url, err := url.Parse(baseURL)
  url.Path = path.Join(url.Path, endPoint)

  // создаем параметры
  q := url.Query()
  q.Set("ClientId", client.clientId)
  q.Set("TransactionId", strconv.FormatInt(transactionId, 10) ) //strconv.Itoa(transactionId))

  // добавляем к URL параметры
  url.RawQuery = q.Encode()

  resp, err := client.GetHttp(ctx, metod, url.String(), nil)

  type responseData struct {
    Error  ResponseError  `json:"error"`
  }
  var rd responseData
  if err = json.Unmarshal(resp, &rd); err != nil {
    client.Logger.Error("DeleteOrder Ошибка при разборе ответа JSON", "err", err.Error())
    return  err
  }

  // если нет ошибки = вернем ok
  //client.Logger.Debug("DeleteOrder", slog.Any("rd.Error",rd.Error))
  if rd.Error.Code != ""{
    return fmt.Errorf("Ошибка %s %s",rd.Error.Code, rd.Error.Message)
  }


  return nil
}

// создадим и вернем частично заполненную структуру для запроса нового ордера
func (client *Client) NewOrder(board, symbol string, sideType SideType, lot int32 ,price float64 ) NewOrderRequest{
  // условие по времени действия заявки
  // TillEndSession - заявка действует до конца сессии;
  // TillCancelled - заявка действует, пока не будет отменена;
  // ExactTime - заявка действует до указанного времени. Параметр time должен быть задан.
  validBefore := &OrderValidBefore{
      Type: "TillEndSession",
  }
  // если хотим зайти по рынку = 
  // Для рыночной заявки указать значение null (или не передавать это поле)
  var order_price *float64
  if price !=0 {
    order_price = &price
  }
  // LastUp   // сделка на рынке по заданной цене или выше;
  // LastDown // сделка на рынке по заданной цене или ниже.
  condition_type := "LastUp" 
  if price == 0 {
    condition_type = "Bid" 
  }
  if sideType == SideBuy{
    condition_type = "LastDown"
    if price == 0 {
      condition_type = "Bid" //"Ask" 
    }

  }
  // если хотим по рынку. price ставим = 0
  // НО в condition_price не допускается нудевая цена
  condition_price := price
  if price == 0 {
    condition_price = 1 
  }
  condition := &OrderCondition{
      Type:  condition_type,
      Price: condition_price,
  }

  // Property свойства исполнения частично исполненных заявок. Принимает следующие значения:
  // PutInQueue - неисполненная часть заявки помещается в очередь заявок биржи;
  // CancelBalance - неисполненная часть заявки снимается с торгов;
  // ImmOrCancel - сделки совершаются только в том случае, если заявка может быть удовлетворена полностью и сразу при выставлении.

  newOrder := NewOrderRequest{
    ClientId:      client.clientId,
    SecurityBoard: board,
    SecurityCode:  symbol,
    BuySell:       sideType.String(),
    Price:         order_price,
    Quantity:      lot,
    UseCredit:     true,
    Property:      "PutInQueue",
    Condition:     condition,
    ValidBefore:   validBefore,
  }

  return newOrder
}



// послать новую заявку
// проверка на корректность цен на совести пользователя
func (client *Client) SendOrder(ctx context.Context, order NewOrderRequest) (int64, error){
  endPoint := "public//api/v1/orders/"
  metod    := "POST"

  url, err := url.Parse(baseURL)
  if err != nil {
    return 0, err
  }
  url.Path = path.Join(url.Path, endPoint)

  type Data struct{
    SecurityCode string `json:"securityCode"`
    TransactionId int64 `json:"transactionId"`
  } 

  type responseData struct {
    Error  ResponseError  `json:"error"`
    Data   Data           `json:"data"`
  }
  var rd responseData

  resp, err := client.GetHttp(ctx, metod, url.String(), order)
  if err != nil {
    return 0, err 
  }
  if err = json.Unmarshal(resp, &rd); err != nil {
    client.Logger.Error("SendOrder Ошибка при разборе ответа JSON", "err", err.Error())
    return 0, err
  }

  // если есть ошибка 
  if rd.Error.Code != ""{
    //return 0, fmt.Errorf(rd.Error.Message, rd.Error.Data)
    return 0, fmt.Errorf("Ошибка %s %s",rd.Error.Code, rd.Error.Message)
  }  


  return rd.Data.TransactionId, nil
}

// купить по рынку
func (client *Client) BuyMarket(ctx context.Context, board, symbol string, lot int32 ) (int64 , error){
  newOrder  := client.NewOrder(board, symbol, SideBuy, lot , 0)
  t_id, err := client.SendOrder(ctx, newOrder)
  if err !=nil{
    client.Logger.Error("BuyMarket", "err", err.Error())
    return 0, err
  }  

  return t_id, nil
}

// выставить лимитную заявку на покупку
func (client *Client) BuyLimit(ctx context.Context, board, symbol string, lot int32, price float64 ) (int64 , error){
  newOrder  := client.NewOrder(board, symbol, SideBuy, lot , price)
  //newOrder.Condition
  t_id, err := client.SendOrder(ctx, newOrder)
  if err !=nil{
    client.Logger.Error("BuyLimit", "err", err.Error())
    return 0, err
  }  

  return t_id, nil
}

// продать по рынку
func (client *Client) SellMarket(ctx context.Context, board, symbol string, lot int32 ) (int64 , error){
  newOrder  := client.NewOrder(board, symbol, SideSell, lot , 0)
  t_id, err := client.SendOrder(ctx, newOrder)
  if err !=nil{
    client.Logger.Error("SellMarket", "err", err.Error())
    return 0, err
  }  

  return t_id, nil
}

// выставить лимитную заявку на продажу
func (client *Client) SellLimit(ctx context.Context, board, symbol string, lot int32, price float64 ) (int64 , error){
  newOrder  := client.NewOrder(board, symbol, SideSell, lot , price)
  t_id, err := client.SendOrder(ctx, newOrder)
  if err !=nil{
    client.Logger.Error("SellLimit", "err", err.Error())
    return 0, err
  }  

  return t_id, nil
}