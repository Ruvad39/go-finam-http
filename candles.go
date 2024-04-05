package finam

import (
  "context"
  "encoding/json"
  "net/url"
  "path"
  "fmt"
)

// Для получения списка свечей необходимо выполнить GET запрос на /api/v1/day-candles/ или /api/v1/intraday-candles.
// Необходимо указать securityCode, securityBoard и timeFrame.
// https://finamweb.github.io/trade-api-docs/rest-api/candles   

//https://trade-api.finam.ru/public/api/v1/day-candles?SecurityBoard=TQBR&SecurityCode=SBER&TimeFrame=D1&Interval.From=2023-04-01&Interval.To=2024-03-31


func (client *Client) GetCandles(ctx context.Context, 
                                board string, 
                                symbol string,
                                timeFrame TimeFrame,
                                from string,
                                to string,

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
  q.Set("Interval.From", from)
  if to != ""{
    q.Set("Interval.To", to)
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
