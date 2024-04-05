package main

import (
    "context"
    "log/slog"
    "os"
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

	// проверка токена
	ok, err := client.AccessTokens(ctx)
	if err != nil{
		slog.Info("main.AccessTokens", "ошибка проверки токена:", err.Error())
		return 
	}
	slog.Info("main.AccessTokens", "ok", ok)

	
	// запрос списка инсртументов
	board := "TQBR" // FUT TQBR
	symbol := "" // "SiM4"
	Sec, err := client.GetSecurity(ctx, board, symbol)
	if err != nil{
		slog.Info("main.GetSecurity", "err", err.Error())
	}

	//slog.Info("GetSecurity", slog.Any("Sec", Sec))
	//slog.Info("GetSecurity", "кол-во", len(Sec))
	// список инструментов
	 for n, sec := range Sec {
	 	slog.Info("securities",
	 		"n_row",     n, 
			"Code",      sec.Code,
			"Market",    sec.Market,
			"board",     sec.Board,  
			"ShortName", sec.ShortName,
			"LotSize",   sec.LotSize,
			"Decimals",  sec.Decimals,
			"MinStep",   sec.MinStep,
		)
	}



}