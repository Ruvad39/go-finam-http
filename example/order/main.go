package main

import (
	"context"
	"github.com/Ruvad39/go-finam-http"
	"log/slog"
)

func main() {
	ctx := context.Background()

	// создание клиента
	token := "
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
	tId := int64(15528)
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
	// newOrder := client.NewOrder("TQBR","SIBN", finam.SideBuy, 1 , 772)
	// t_id, err := client.SendOrder(ctx, newOrder)

	//tId, err := client.NewCreateOrderService("TQBR", "SBER", finam.SideBuy, 1).Do(ctx)
	//tId, err = client.NewCreateOrderService("TQBR", "SBER", finam.SideBuy, 1).Price(307.84).Do(ctx)
	//if err != nil {
	//	slog.Error("main.SendOrder", "err", err.Error())
	//	return
	//}
	//slog.Info("SendOrder", "TransactionId", tId)

	// купить по рынку
	//t_id, err := client.BuyMarket(ctx, "TQBR","SBER", 2)

	// выставить лимитную заявку на покупку
	//t_id, err := client.BuyLimit(ctx, "TQBR","SBER", 1, 306.02)

	// продать по рынку
	//t_id, err := client.SellMarket(ctx, "TQBR","SBER", 3)

	// выставить лимитную заявку на продажу
	//t_id, err := client.SellLimit(ctx, "TQBR","SBER", 2, 306.06)

	// if err !=nil{
	// 	slog.Error("main.BuySell", "err", err.Error())
	// }
	// slog.Info("BuySell", "TransactionId", t_id)

}
