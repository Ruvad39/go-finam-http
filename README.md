# Библиотека, которая позволяет работать с функционалом [Finam Trade API через REST IP](https://finamweb.github.io/trade-api-docs/category/rest-api)  брокера [Финам](https://www.finam.ru/) из GO



## Установка


```bash
go get github.com/Ruvad39/go-finam-http
```

## api реализован на текущий момент:
```go
// проверка подлинности токена
AccessTokens(ctx context.Context) (ok bool, err error)
// Посмотреть портфель
GetPortfolio(ctx context.Context, opts ...Option) (Portfolio, error)
// список инструментов (Максимальное Количество запросов в минуту = 1 )
GetSecurity(ctx context.Context, board string, seccode string) ( Securities, error)
// получить свечи
GetCandles(ctx context.Context, board, symbol string, timeFrame TimeFrame, from, to string, count int) ([]Candle, error)
// получить список заявок
GetOrders(ctx context.Context, opts ...Option) ( []Order, error)
// отменить заявку
DeleteOrder(ctx context.Context, transactionId int64) error
// создать новую заявку
// купить по рынку
BuyMarket(ctx context.Context, board, symbol string, lot int32 ) (int64 , error)
// продать по рынку
SellMarket(ctx context.Context, board, symbol string, lot int32 ) (int64 , error)
// выставить лимитную заявку на покупку
BuyLimit(ctx context.Context, board, symbol string, lot int32, price float64 ) (int64 , error)
// выставить лимитную заявку на продажу
SellLimit(ctx context.Context, board, symbol string, lot int32, price float64 ) (int64 , error)

// TODO
// стоп-завки

```

## Примеры

### создание клиента
```go
token := "token"
clientId := "client_id"

client, err := finam.NewClient(token, clientId)
if err != nil {
    slog.Error("main", slog.Any("ошибка создания finam.client", err))
}

ctx := context.Background()
// проверка токена
ok, err := client.AccessTokens(ctx)
if err != nil{
slog.Info("main.AccessTokens", "ошибка проверки токена:", err.Error())
return
}
slog.Info("main.AccessTokens", "ok", ok)
````

### Пример получения данных о портфеле

```go
// запрос через создание сервиса
portfolio, err := client.NewGetPortfolioService().
    IncludeCurrencies(true).
    IncludePositions(true).
    IncludeMoney(true).
    IncludeMaxBuySell(true).Do(ctx)

// запрос через вызов метода
//portfolio, err := client.GetPortfolio(ctx,
//    finam.WithIncludePositions(true),
//    finam.WithIncludeCurrencies(true),
//    finam.WithIncludeMoney(true),
//    finam.WithIncludeMaxBuySell(true))

if err != nil {
    slog.Info("main.GetPortfolio", "err", err.Error())
return
}

// баланс счета
slog.Info("Balance", "Equity", portfolio.Equity, "Balance", portfolio.Balance)

// список позиций
for _, pos := range portfolio.Positions {
    slog.Info("position", slog.Any("pos", pos))
}
// список валют счета
slog.Info("portfolio.Currencies", slog.Any("Currencies", portfolio.Currencies))
// список денег
slog.Info("portfolio.Money", slog.Any("Money", portfolio.Money))

```

### Пример получения свечей

```go
board := "TQBR"  // FUT TQBR
symbol := "SBER" // "SiM4"
// дневные свечи
tf := finam.TimeFrame_D1
from, _ := time.Parse("2006-01-02", "2024-03-01")
to, _ := time.Parse("2006-01-02", "2024-04-20")

// через создание сервиса
// candles, err := client.NewCandlesService(board, symbol, tf).StartTime(from).EndTime(to).Do(ctx)

// через вызов метода
candles, err := client.GetCandles(ctx, board, symbol, tf, finam.WithStartTime(from), finam.WithEndTime(to))
if err != nil {
    slog.Info("main.candles", "err", err.Error())
    return
}
slog.Info("GetCandles", "кол-во", len(candles))
//список свечей
for n, candle := range candles {
    slog.Info("candles",
        "row", n,
        "datetime", candle.GetDateTimeToTime().String(),
        "candle", candle.String(),
    )
}

// внутредневные:
// выберем последние 5 минутных свечей
tf = finam.TimeFrame_M1
time_to := time.Now()
count := 5

// через создание сервиса
// candles, err = client.NewCandlesService(board, symbol, tf).EndTime(time_to).Count(count).Do(ctx)

// через вызов метода
candles, err = client.GetCandles(ctx, board, symbol, tf, finam.WithEndTime(time_to), finam.WithCount(count))

if err != nil {
    slog.Info("main.GetCandles", "err", err.Error())
}
slog.Info("GetCandles", "кол-во", len(candles))
//список свечей
for n, candle := range candles {
    slog.Info("candle",
        "row", n,
        "datetime", candle.GetDateTimeToTime().String(),
        "candle", candle.String(),
    )
}

```

### Другие примеры смотрите [тут](/example)

---
Многие идеи по организации структуры клиента принадлежат [adshao](https://github.com/adshao)  
из его проекта [go-binance](https://github.com/adshao/go-binance)
за что ему большое спасибо