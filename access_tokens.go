package finam

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/url"
	"path"
	"fmt"
)

// https://trade-api.finam.ru/public/api/v1/access-tokens/check
// проверка токена
func (client *Client) AccessTokens(ctx context.Context) ( ok bool, err error){
	endPoint := "public/api/v1/access-tokens/check"
	ok = false

	url, err := url.Parse(baseURL)
	if err != nil {
		return false, err
	}
	url.Path = path.Join(url.Path, endPoint)

	resp, err := client.GetHttp(ctx,"GET", url.String(), nil)
	if err != nil {
		return false, err 
	}
	//client.Logger.Debug("AccessTokens", "resp", resp)

	type responseData struct {
		Error  ResponseError  `json:"error"`
	}
	var rd responseData
	if err = json.Unmarshal(resp, &rd); err != nil {
		client.Logger.Error("AccessTokens Ошибка при разборе ответа JSON", "err", err.Error())
		return false, err
	}

	// если нет ошибки = вернем ok
	if rd.Error.Code == ""{
		return true, nil
	}

	client.Logger.Debug("AccessTokens", slog.Any("rd", rd))	
	return false , fmt.Errorf(rd.Error.Message)


}
