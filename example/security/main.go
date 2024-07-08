package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ruvad39/go-finam-http"
	"log/slog"
	"os"
)

func main() {

	ctx := context.Background()

	// создание клиента
	token := ""
	clientId := ""

	client := finam.NewClient(token, clientId)
	//client.Debug = true

	// запрос списка инструментов
	// через создание service
	//Sec, err := client.NewSecurityService().Board("TQBR").Symbol("SBER").Do(ctx)

	// через вызов метода
	Sec, err := client.GetSecurity(ctx, "TQBR", "")
	if err != nil {
		slog.Info("main.GetSecurity", "err", err.Error())
		return
	}

	slog.Info("GetSecurity", "кол-во", len(Sec))
	// список инструментов
	//for n, sec := range Sec {
	//	slog.Info("securities",
	//		"row", n,
	//		"Code", sec.Code,
	//		"Market", sec.Market,
	//		"board", sec.Board,
	//		"ShortName", sec.ShortName,
	//		"LotSize", sec.LotSize,
	//		"Decimals", sec.Decimals,
	//		"MinStep", sec.MinStep,
	//	)
	//}
	// запишу в свой файл
	type security struct {
		Symbol    string  `json:"symbol"`     // Код инструмента
		ShortName string  `json:"short_name"` // Краткое наименование
		Board     string  `json:"board"`      // Код класса (Board)
		LotSize   int     `json:"lot_size"`   // Размер лота
		MinStep   float32 `json:"minStep"`    // минимальный шаг цены;
		Decimals  int     `json:"decimals"`   // количество знаков в дробной части цены;
	}
	len := len(Sec)
	Securites_ := make([]security, 0, len)

	for n, sec := range Sec {
		if sec.Board == "TQBR" || sec.Board == "FUT" {
			slog.Info("securities",
				"row", n,
				"Code", sec.Code,
				"Market", sec.Market,
				"board", sec.Board,
				"ShortName", sec.ShortName,
				"LotSize", sec.LotSize,
				"Decimals", sec.Decimals,
				"MinStep", sec.MinStep,
			)

			s := security{}
			s.Symbol = sec.Code
			s.Board = sec.Board
			s.ShortName = sec.ShortName
			s.LotSize = sec.LotSize
			s.Decimals = sec.Decimals
			s.MinStep = sec.MinStep

			Securites_ = append(Securites_, s)
		}
	}
	// запишем в файл
	jsonData, err := json.Marshal(Securites_)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = os.WriteFile("security.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

}
