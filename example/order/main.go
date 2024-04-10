package main

import (
    "context"
    "log/slog"
    "os"

    "github.com/Ruvad39/go-finam-http"
    "github.com/joho/godotenv"
)

func main(){
	ctx := context.Background()

    // загрузим значения для переменных окружения .env 
    if err := godotenv.Load(); err != nil {
        slog.Error("No .env file found")
    }

	// Level: slog.LevelDebug,
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})		
	logger_ := slog.New(handler)

	// создание клиента
	//token := ""
	//clientId := ""
	// значения из переменных окружения
	token, _    := os.LookupEnv("FINAM_TOKEN")
	clientId, _ := os.LookupEnv("FINAM_CLIENTID")

	client, err := finam.NewClient(token, clientId, finam.WithLogger(logger_))

	if err != nil {
		slog.Error("main", slog.Any("ошибка создания finam.client", err))
	}
	
	// запрос списка активных ордеров
	// WithIncludeMatched(true)  // вернуть исполненные заявки
	// WithIncludeCanceled(true) // вернуть отмененные заявки;
	// WithIncludeActive(true)   // вернуть активные заявки.
	orders, err := client.GetOrders(ctx, finam.WithIncludeActive(true), )
	if err != nil{
		slog.Info("main.GetOrders", "err", err.Error())
	}

	slog.Info("GetOrders", "кол-во", len(orders))
	//список заявок
	for n, order := range orders {
		slog.Info("GetOders", "row", n, slog.Any("order", order))
	}


	// удалим заявку. нужно послать TransactionId ордера
	// err = client.DeleteOrder(ctx, 29491098 )
	// if err != nil{
	// 	slog.Error("main.DeleteOrders", "err", err.Error())
	// } else {
	// 	slog.Info("DeleteOrder OK")
	// }
	

	// новая заявка ( board, symbol, sideType, lot, price)
	// если хотим купить\продать по рунку: ставим price = 0 или ставим цену ниже\выше текущей
	// newOrder := client.NewOrder("TQBR","SIBN", finam.SideBuy, 1 , 772)
	// t_id, err := client.SendOrder(ctx, newOrder)
	// if err !=nil{
	// 	slog.Error("main.SendOrder", "err", err.Error())
	// }
	// slog.Info("SendOrder", "TransactionId", t_id)

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