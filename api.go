package finam
import (
	"context"
)

// какой api реализован
type IFinamClient interface {
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
	SendOrder(ctx context.Context, order NewOrderRequest) (int64, error)
	// купить по рынку
	BuyMarket(ctx context.Context, board, symbol string, lot int32 ) (int64 , error)
	// выставить лимитную заявку на покупку
	BuyLimit(ctx context.Context, board, symbol string, lot int32, price float64 ) (int64 , error)	
	// продать по рынку
	SellMarket(ctx context.Context, board, symbol string, lot int32 ) (int64 , error)
	// выставить лимитную заявку на продажу
	SellLimit(ctx context.Context, board, symbol string, lot int32, price float64 ) (int64 , error)	

	// TODO
	// получить список стоп-заявок
	// создать новую стоп-заявку
	// отменить стоп-заявку

}
