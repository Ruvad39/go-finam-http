/*
	идея логики клиента принадлежит https://github.com/adshao
	и его проекту https://github.com/adshao/go-binance
*/

package finam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	libraryName    = "FINAM-REST API GO"
	libraryVersion = "0.0.2"
	baseAPIMainURL = "https://trade-api.finam.ru/"
	headerKey      = "X-Api-Key"
)

// UseDevelop использовать тестовый или боевой сервер
//var UseDevelop = false

// getAPIEndpoint return the base endpoint of the Rest API according the UseDevelop flag
func getAPIEndpoint() string {
	//if UseDevelop {
	//	return baseAPITestnetURL
	//}
	return baseAPIMainURL
}

// NewClient создание нового клиента
func NewClient(token, clientId string) *Client {
	return &Client{
		token:      token,
		clientId:   clientId,
		BaseURL:    getAPIEndpoint(),
		UserAgent:  "Finam/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "go-finam ", log.LstdFlags),
	}
}

type doFunc func(req *http.Request) (*http.Response, error)

// Client define API client
type Client struct {
	token      string
	clientId   string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger
	TimeOffset int64
	do         doFunc
}

func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}
	err = r.validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint)

	queryString := r.query.Encode()
	//body := &bytes.Buffer{}
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}
	if r.body != nil {
		header.Set("Content-Type", "application/json")
		c.debug("r.body: %s", r.body)
		//body = r.body
	}

	bodyString := r.form.Encode()
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.body = bytes.NewBufferString(bodyString)
		//body = bytes.NewBufferString(bodyString)
		c.debug("bodyString: %s", bodyString)
	}
	//headerKey := "X-Api-Key"
	if c.token != "" {
		//header.Set("X-Api-Key", c.token)
		header.Set(headerKey, c.token)
	}

	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	//c.debug("full url: %s, body: %s", fullURL, bodyString)
	c.debug("full url: %s", fullURL)

	r.fullURL = fullURL
	r.header = header
	//r.body = body
	return nil
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	c.debug("request: %#v", req)
	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return []byte{}, err
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	c.debug("response: %#v", res)
	c.debug("response body: %s", string(data))
	c.debug("response status code: %d", res.StatusCode)

	// из финама приходит другая структура ошибки
	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(APIError)
		e := json.Unmarshal(data, apiErr)
		if e != nil {
			c.debug("failed to unmarshal json: %s", e)
		}
		return nil, apiErr
	}
	return data, nil
}

// структура ошибки
type ResponseError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type APIError struct {
	ResponseError ResponseError `json:"error"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("<APIError> code=%s, msg=%s, data=%s", e.ResponseError.Code, e.ResponseError.Message, e.ResponseError.Data)
}

func IsAPIError(e error) bool {
	_, ok := e.(*APIError)
	return ok
}

// (debug) вернем текущую версию
func (c *Client) Version() string {
	return libraryName + " v." + libraryVersion
}
