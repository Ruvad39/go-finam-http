package finam

import (
    "fmt"
    "context"
    "os"
    "io/ioutil"
    "log/slog"
    "net/http"
    //"net/url"
    
)

const (
    libraryName    = "FINAM API REST GO"
    libraryVersion = "0.0.1"
    baseURL = "https://trade-api.finam.ru/"
)


type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

type Client struct {
    token       string 
    clientId    string 

    UserAgent   string    // если проставлен, пропишем User agent в http запросе
    httpClient  HTTPClient //*http.Client
    Logger      *slog.Logger

}


// создание клиента
func NewClient(token, clientId string, opts ...ClientOption) (*Client, error) {
    c := &Client{
        token:       token,
        clientId:    clientId,
        httpClient:  http.DefaultClient,
        Logger :     slog.New(slog.NewTextHandler(os.Stdout, nil)), //io.Discard
    }
    // обрабратаем входящие параметры
    for _, opt := range opts {
        opt(c)
    }

    return c, nil
}


// выполним запрос  (вернем http.Response)
func (client *Client) RequestHttp(ctx context.Context, httpMethod string, url string, body interface{})(*http.Response, error){

    req, err := http.NewRequestWithContext(ctx, httpMethod, url, nil) 
    if err != nil {
        client.Logger.Error("RequestHttp", "httpMethod", httpMethod, "url", url, "err", err.Error())
        return nil, err
    }

    // добавляем заголовки
    if client.UserAgent != "" {
        req.Header.Set("User-Agent", client.UserAgent)
    }    
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    if client.token != ""{
        req.Header.Add("X-Api-Key", client.token)
    }
    
    resp, err := client.httpClient.Do(req)
    if err != nil {
        client.Logger.Error("RequestHttp", "httpMethod", httpMethod, "url", url, "err", err.Error())
        return nil, err
    }

    client.Logger.Debug("RequestHttp", "httpMethod", httpMethod, "url", url, "StatusCode", resp.StatusCode)


    return resp, err
}

// выполним запрос  (вернем []byte)
func (client *Client) GetHttp(ctx context.Context, httpMethod string, url string, body interface{})([]byte, error){

    resp, err := client.RequestHttp(ctx, httpMethod, url, body )

    if err != nil {
        client.Logger.Error("RequestHttp", "httpMethod", httpMethod, "url", url, "err", err.Error())
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        //client.Logger.Error("RequestHttp", slog.Any("resp",resp))
        return nil, fmt.Errorf(resp.Status)

    //     return nil, fmt.Errorf("responce StatusCode %s", resp.StatusCode)
    }

    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)

}


// вернем текущую версию
func (c *Client) Version() string{
    return libraryVersion
}




type ClientOption func(c *Client)

// WithLogger задает логгер 
// По умолчанию логирование включено на ошибки
func WithLogger(logger *slog.Logger) ClientOption {
    return func(opts *Client) {
        opts.Logger = logger
    }
}

func WithClientId(id string) ClientOption {
    return func(opts *Client) {
        opts.clientId = id
    }
}

// установим свой HttpClient
// по умолчанию стоит http.DefaultClient
func WithGttpClient(client HTTPClient) ClientOption {
    return func(opts *Client) {
        opts.httpClient = client
    }
}

