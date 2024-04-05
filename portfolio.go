package finam

import (
  "context"
  "encoding/json"
  "net/url"
  "path"
  "fmt"
)

// Посмотреть портфель
// https://finamweb.github.io/trade-api-docs/rest-api/portfolios      

// clientId - торговый код клиента (обязательный);
// includeCurrencies - запросить информацию по валютам портфеля;
// includeMoney - запросить информацию по денежным позициям портфеля;
// includePositions - запросить информацию по позициям портфеля;
// includeMaxBuySell - запросить информацию о максимальном доступном объеме на покупку/продажу.

func (client *Client) GetPortfolio(ctx context.Context, opts ...Option) ( Portfolio, error){
  endPoint := "public/api/v1/portfolio"

  type responseData struct {
    Error  ResponseError  `json:"error"`
    Data   Portfolio      `json:"data"`
  }
  var rd responseData

  p := &Options{
    IncludePositions: true,
  }
  // обработаем входящие параметры
  for _, opt := range opts {
      opt(p)
  }  

  url, err := url.Parse(baseURL)
  if err != nil {
    return rd.Data, err
  }
  url.Path = path.Join(url.Path, endPoint)

  // создаем параметры
  q := url.Query()

  q.Set("ClientId", client.clientId)

  if p.IncludeCurrencies{
    q.Set("Content.includeCurrencies", "true")
  }
  if p.IncludePositions{
    q.Set("Content.IncludePositions", "true")
  }
  if p.IncludeMoney{
    q.Set("Content.IncludeMoney", "true")
  }
  if p.IncludeMaxBuySell{
    q.Set("Content.IncludeMaxBuySell", "true")
  }

  // добавляем к URL параметры
  url.RawQuery = q.Encode()

  resp, err := client.GetHttp(ctx,"GET", url.String(), nil)
  if err != nil {
    return rd.Data, err 
  }
  //client.Logger.Debug("GetPortfolio", "resp", resp)

  if err = json.Unmarshal(resp, &rd); err != nil {
    client.Logger.Error("GetPortfolio Ошибка при разборе ответа JSON", "err", err.Error())
    return rd.Data, err
  }

  // если есть ошибка 
  if rd.Error.Code != ""{
    return rd.Data, fmt.Errorf(rd.Error.Message)
  }

  return rd.Data, nil
}

