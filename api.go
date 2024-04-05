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
	GetCandles(ctx context.Context, board, symbol string, timeFrame TimeFrame, from, to string) ([]Candle, error)

	// TODO
	// получить список заявок
	// создать новую заявку
	// отменить заявку
	// получить список стоп-заявок
	// создать новую стоп-заявку
	// отменить стоп-заявку

}
