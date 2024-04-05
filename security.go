package finam

import (
  "context"
  "encoding/json"
  "net/url"
  "path"
  "fmt"
)


// запросить список инструментов
// Максимальное Количество запросов в минуту = 1
// https://finamweb.github.io/trade-api-docs/rest-api/securities    

func (client *Client) GetSecurity(ctx context.Context, board string, seccode string) ( Securities, error){
  endPoint := "public/api/v1/securities/"

  type Data struct{
    Sec Securities `json:"securities"`
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

  if board != ""{
    q.Set("Board", board)
  }
  if seccode != ""{
    q.Set("Seccode", seccode)
  }

  // добавляем к URL параметры
  url.RawQuery = q.Encode()

  resp, err := client.GetHttp(ctx,"GET", url.String(), nil)
  if err != nil {
    return rd.Data.Sec, err 
  }
  //client.Logger.Debug("GetSecurity", "resp", resp)

  if err = json.Unmarshal(resp, &rd); err != nil {
    client.Logger.Error("GetSecurity Ошибка при разборе ответа JSON", "err", err.Error())
    return rd.Data.Sec, err
  }

  // если есть ошибка 
  if rd.Error.Code != ""{
    return rd.Data.Sec, fmt.Errorf(rd.Error.Message)
  }

  return rd.Data.Sec, nil
}
