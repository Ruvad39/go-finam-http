package main

import (
	"context"
	"github.com/Ruvad39/go-finam-http"
	"log/slog"
)

func main() {
	ctx := context.Background()

	// создание клиента
	token := ""
	clientId := ""

	client := finam.NewClient(token, clientId)
	//client.Debug = true

	// проверка токена
	ok, err := client.AccessTokens(ctx)
	if err != nil {
		slog.Info("main.AccessTokens", "ошибка проверки токена:", err.Error())
		return
	}
	slog.Info("main.AccessTokens", "ok", ok)

	// запрос состояния счета
	// первый способ: создание service
	//portfolio, err := client.NewGetPortfolioService().
	//	IncludePositions(true).
	//	IncludeCurrencies(true).
	//	IncludeMoney(true).
	//	IncludeMaxBuySell(true).Do(ctx)

	// второй способ: вызов метода
	portfolio, err := client.GetPortfolio(ctx,
		finam.WithIncludePositions(true),
		finam.WithIncludeCurrencies(true),
		finam.WithIncludeMoney(true),
		finam.WithIncludeMaxBuySell(true))

	if err != nil {
		slog.Info("main.GetPortfolio", "err", err.Error())
		return
	}

	// баланс счета
	//slog.Info("Balance", "Equity", portfolio.Equity, "Balance", portfolio.Balance)
	slog.Info("GetPortfolio", "portfolio", portfolio)
	// список позиций
	for _, pos := range portfolio.Positions {
		slog.Info("position", slog.Any("pos", pos))
	}
	// список валют счета
	slog.Info("portfolio.Currencies", slog.Any("Currencies", portfolio.Currencies))
	// список денег
	slog.Info("portfolio.Money", slog.Any("Money", portfolio.Money))

}
