package finam

import (
  "context"
  "encoding/json"
  "net/url"
  "path"
  "fmt"
  "strconv"
)

// Необходимо указать securityCode, securityBoard и timeFrame.
// https://finamweb.github.io/trade-api-docs/rest-api/candles   

// Запрос дневных/недельных свечей
// Максимальный интервал: 365 дней
// Максимальное кол-во запросов в минуту: 120
// Дата начала (окончания) в формате yyyy-MM-dd в часовом поясе UTC
// Запрос внутридневных свечей
// Максимальный интервал: 30 дней
// Максимальное кол-во запросов в минуту: 120
// Дата начала (окончания) в формате yyyy-MM-ddTHH:mm:ssZ в часовом поясе UTC
func (client *Client) GetCandles(ctx context.Context, board string, symbol string, timeFrame TimeFrame,
                                from string,
                                to string,
                                count int,

) ( []Candle, error){
  endPoint := "public/api/v1/day-candles"
  if !timeFrame.IsDay(){
    endPoint = "public/api/v1/intraday-candles"
  }


  type Data struct{
    Sec Candles `json:"candles"`
  } 

  type responseData struct {
    Error  ResponseError  `json:"error"`
    Data   Data           `json:"data"`
  }

  var rd responseData

  url, err := url.Parse(baseURL)
  if err != nil {
    return rd.Data.Sec, err
  }
  url.Path = path.Join(url.Path, endPoint)

  // создаем параметры
  q := url.Query()

  q.Set("SecurityBoard", board)
  q.Set("SecurityCode", symbol)
  q.Set("TimeFrame", timeFrame.String())

  // TODO проверка дат

  if from !=""{
    q.Set("Interval.From", from)
  }
    if to != ""{
    q.Set("Interval.To", to)
  }
  if count != 0{
    q.Set("Interval.Count", strconv.Itoa(count))
  }
  
  // добавляем к URL параметры
  url.RawQuery = q.Encode()

  resp, err := client.GetHttp(ctx,"GET", url.String(), nil)
  if err != nil {
    return rd.Data.Sec, err 
  }
  //client.Logger.Debug("GetCandles", "resp", resp)

  if err = json.Unmarshal(resp, &rd); err != nil {
    client.Logger.Error("GetCandles Ошибка при разборе ответа JSON", "err", err.Error())
    return rd.Data.Sec, err
  }

  // если есть ошибка 
  if rd.Error.Code != ""{
    return rd.Data.Sec, fmt.Errorf(rd.Error.Message)
  }

  return rd.Data.Sec, nil
}
