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
	client.Debug = true
	//slog.Info(client.Version())
	//return

	// через создание service
	// вернуть активные заявки стоит по умолчанию
	//orders, err := client.NewGetOrderService().IncludeActive(true).Do(ctx)
	// orders, err := client.NewGetOrderService().Do(ctx)

	// через вызов метода
	// WithIncludeMatched(true)  // вернуть исполненные заявки
	// WithIncludeCanceled(true) // вернуть отмененные заявки;
	// WithIncludeActive(true)   // вернуть активные заявки.
	orders, err := client.GetOrders(ctx, finam.WithIncludeActive(true))
	if err != nil {
		slog.Info("main.GetOrders", "err", err.Error())
		return
	}

	slog.Info("GetOrders", "кол-во", len(orders))
	//список заявок
	for n, order := range orders {
		slog.Info("GetOders", "row", n, slog.Any("order", order))
	}

	// удалим заявку
	// нужно послать TransactionId ордера
	tId := int64(22077388)
	// через создание service
	err = client.NewCancelOrderService(tId).Do(ctx)
	// через вызов метода
	//err = client.DeleteOrder(ctx, tId)
	if err != nil {
		slog.Error("main.CancelOrders", "err", err.Error())
		//return
	} else {
		slog.Info("CancelOrder OK")
	}

	// новая заявка ( board, symbol, sideType, lot, price)
	// если хотим купить\продать по рынку: ставим price = 0 или ставим цену ниже\выше текущей

	// купить по рынку
	//tId, err = client.BuyMarket(ctx, "FUT", "SiM4", 1)
	//tId, err = client.NewCreateOrderService("FUT", "SiM4", finam.SideBuy, 1).Do(ctx)

	// продать по рынку
	tId, err = client.SellMarket(ctx, "FUT", "SiM4", 1)
	//tId, err = client.NewCreateOrderService("FUT", "SiM4", finam.SideSell, 1).Do(ctx)

	// лимитная заявка на продажу
	//tId, err = client.SellLimit(ctx, "FUT", "SiM4", 1, 92556)
	//tId, err = client.NewCreateOrderService("FUT", "SiM4", finam.SideSell, 1).Price(92556).Do(ctx)

	// лимитная заявка на покупку
	//tId, err = client.BuyLimit(ctx, "FUT", "SiM4", 1, 92584)
	//tId, err = client.NewCreateOrderService("FUT", "SiM4", finam.SideBuy, 1).Price(92584).Do(ctx)

	if err != nil {
		slog.Error("main.SendOrder", "err", err.Error())
		return
	}
	slog.Info("SendOrder", "TransactionId", tId)

}
