package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/Ruvad39/go-finam-http"
)

func main() {
	ctx := context.Background()
	// создание клиента
	token := ""
	clientId := ""

	client := finam.NewClient(token, clientId)
	client.Debug = false

	// запрос свечей
	board := "TQBR"  // FUT TQBR
	symbol := "SBER" // "SiM4"

	// дневные свечи
	tf := finam.TimeFrame_D1
	from, _ := time.Parse("2006-01-02", "2024-04-01")
	to, _ := time.Parse("2006-01-02", "2024-04-30")

	// через создание сервиса
	//candles, err := client.NewCandlesService(board, symbol, tf).StartTime(from).EndTime(to).Do(ctx)

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
	//-------------------------------------
	// внутредневные:
	// выберем последние 5 минутных свечей
	tf = finam.TimeFrame_M1
	time_to := time.Now()
	count := 5

	// через создание сервиса
	candles, err = client.NewCandlesService(board, symbol, tf).EndTime(time_to).Count(count).Do(ctx)

	// через вызов метода
	//candles, err = client.GetCandles(ctx, board, symbol, tf, finam.WithEndTime(time_to), finam.WithCount(count))

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

}
