# Библиотека, которая позволяет работать с функционалом [Finam Trade API через REST IP](https://finamweb.github.io/trade-api-docs/category/rest-api)  брокера [Финам](https://www.finam.ru/) из GO



## Установка


```bash
go get github.com/Ruvad39/go-finam-http
```
## какой api реализован на текущий момент
```go
// проверка подлинности токена
AccessTokens(ctx context.Context) (ok bool, err error)
// Посмотреть портфель
GetPortfolio(ctx context.Context, opts ...Option) (Portfolio, error)
// список инструментов (Максимальное Количество запросов в минуту = 1 )
GetSecurity(ctx context.Context, board string, seccode string) ( Securities, error)
// получить свечи
GetCandles(ctx context.Context, board, symbol string, timeFrame TimeFrame, from, to string) ([]Candle, error)
```
## Примеры

### Пример получения данных о портфеле

```go
func main(){
    ctx := context.Background()
    // создание клиента
    token := "token"
    clientId := "client_id"

    client, err := finam.NewClient(token, clientId)
    if err != nil {
        slog.Error("main", slog.Any("ошибка создания finam.client", err))
    }

    // проверка токена
    ok, err := client.AccessTokens(ctx)
    if err != nil{
        slog.Info("main.AccessTokens", "ошибка проверки токена:", err.Error())
        return 
    }
    slog.Info("main.AccessTokens", "ok", ok)

    // запрос состояния счета
    // IncludePositions по умолчанию = true
    portfolio, err := client.GetPortfolio(ctx,
                            finam.WithIncludePositions(true), 
                            finam.WithIncludeCurrencies(true), 
                            finam.WithIncludeMoney(true),
                            finam.WithIncludeMaxBuySell(true),
                        )
    if err != nil{
        slog.Info("main.GetPortfolio", "err", err.Error())
    }

    // баланс счета
    slog.Info("Balance", "Equity", portfolio.Equity, "Balance", portfolio.Balance)

    // список позиций
    for _, pos := range portfolio.Positions {
        slog.Info("position", slog.Any("pos",pos))
        // slog.Info("position",
        //  "SecurityCode",     pos.SecurityCode,
        //  "Market",           pos.Market,
        //  "Balance",          pos.Balance,  
        //  "CurrentPrice",     pos.CurrentPrice,
        //  "AveragePrice",     pos.AveragePrice,
        //  "UnrealizedProfit", pos.UnrealizedProfit,
        // )
    }

    // список валют счета
    slog.Info("portfolio.Currencies" , slog.Any("Currencies", portfolio.Currencies))

    // список денег
    slog.Info("portfolio.Money" , slog.Any("Money", portfolio.Money))

}
```

### Пример получения свечей

```go
func main(){
    ctx := context.Background()

    // Level: slog.LevelDebug,
    handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelDebug,
        })      
    logger_ := slog.New(handler)

    // создание клиента
    token := ""
    clientId := ""

    client, err := finam.NewClient(token, clientId, finam.WithLogger(logger_))
    if err != nil {
        slog.Error("main", slog.Any("ошибка создания finam.client", err))
    }
    
    // запрос свечей
    board := "TQBR" // FUT TQBR
    symbol := "SBER" // "SiM4"
    // TimeFrame_M1 TimeFrame_M5 TimeFrame_M15 TimeFrame_H1 TimeFrame_D1 TimeFrame_W1
    tf := finam.TimeFrame_D1
    from := "2024-03-25"
    to := "2024-04-05"

    // внутредневная
    //tf = finam.TimeFrame_M1
    //from = "2024-03-29T20:06:11Z"

    candles, err := client.GetCandles(ctx, board, symbol, tf, from, to)
    if err != nil{
        slog.Info("main.GetCandles", "err", err.Error())
    }

    slog.Info("GetCandles", "кол-во", len(candles))
    //список свечей
    for n, candle := range candles {
        slog.Info("securities",
            "n_row",     n, 
            "datetime", candle.GetDateTimeToTime().String(),
            "candles",   candle.String(),
        )
    }
}
```

### другие примеры смотрите [тут](/example)