package main

import (
    "context"
    "log/slog"
    "os"
    "time"

    "github.com/Ruvad39/go-finam-http"
)

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

	// дневная. за период
	// TimeFrame_M1 TimeFrame_M5 TimeFrame_M15 TimeFrame_H1 TimeFrame_D1 TimeFrame_W1
	tf := finam.TimeFrame_D1
	from := "2024-04-01" // дата с 
	to := "2024-04-05"   // дата по
	count := 0           // кол-во свечей


	candles, err := client.GetCandles(ctx, board, symbol, tf, from, to, count)
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

	// внутредневная. последняя свеча м1
	// дата доожна быть в формате yyyy-MM-ddTHH:mm:ssZ
	var layout = "2006-01-02T15:04:05Z"
	// время должны быть в в часовом поясе UTC.  нужно перевести
	time_now := time.Now().UTC().Format(layout)
	slog.Info("текущее время", "time_now", time_now)

	tf = finam.TimeFrame_M1
	from = ""
	to   = time_now
	count = 1

	candles, err = client.GetCandles(ctx, board, symbol, tf, from, to, count)
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