package finam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// https://trade-api.finam.ru/public/api/v1/access-tokens/check
// проверка токена
func (c *Client) AccessTokens(ctx context.Context) (ok bool, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "public/api/v1/access-tokens/check",
	}
	ok = false

	data, err := c.callAPI(ctx, r)
	if err != nil {
		return false, err
	}

	type responseData struct {
		Error ResponseError `json:"error"`
	}
	var rd responseData

	if err = json.Unmarshal(data, &rd); err != nil {
		return false, err
	}

	// если нет ошибки = вернем ok
	if rd.Error.Code == "" {
		return true, nil
	}

	c.debug("AccessTokens: %s ", rd)
	return false, fmt.Errorf(rd.Error.Message)

}
