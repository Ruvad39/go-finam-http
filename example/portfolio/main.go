package main

import (
    "context"
    "log/slog"
    "os"
    "github.com/Ruvad39/go-finam-http"
)

func main(){
	ctx := context.Background()
	// создадим логер (для отладки)
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
		// 	"SecurityCode",     pos.SecurityCode,
		// 	"Market",           pos.Market,
		// 	"Balance",          pos.Balance,  
		// 	"CurrentPrice",     pos.CurrentPrice,
		// 	"AveragePrice",     pos.AveragePrice,
		// 	"UnrealizedProfit", pos.UnrealizedProfit,
		// )
	}

	// список валют счета
	slog.Info("portfolio.Currencies" , slog.Any("Currencies", portfolio.Currencies))

	// список денег
	slog.Info("portfolio.Money" , slog.Any("Money", portfolio.Money))

}