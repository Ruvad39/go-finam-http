package finam

import "context"

// какой api реализован
type IFinamClient interface {

	// AccessTokens проверка подлинности токена
	AccessTokens(ctx context.Context) (ok bool, err error)

	// GetPortfolio получить данные по портфелю
	GetPortfolio(ctx context.Context, opts ...Option) (Portfolio, error)

	// GetSecurity получить список инструментов (Максимальное Количество запросов в минуту = 1 )
	GetSecurity(ctx context.Context, board, symbol string) (Securities, error)

	// GetCandles получить свечи
	GetCandles(ctx context.Context, board, symbol string, timeFrame TimeFrame, opts ...Option) ([]Candle, error)

	// получить список заявок
	GetOrders(ctx context.Context, opts ...Option) ([]Order, error)

	// отменить заявку
	DeleteOrder(ctx context.Context, transactionId int64) error

	//// создать новую заявку
	//SendOrder(ctx context.Context, order NewOrderRequest) (int64, error)
	//// купить по рынку
	//BuyMarket(ctx context.Context, board, symbol string, lot int32 ) (int64 , error)
	//// выставить лимитную заявку на покупку
	//BuyLimit(ctx context.Context, board, symbol string, lot int32, price float64 ) (int64 , error)
	//// продать по рынку
	//SellMarket(ctx context.Context, board, symbol string, lot int32 ) (int64 , error)
	//// выставить лимитную заявку на продажу
	//SellLimit(ctx context.Context, board, symbol string, lot int32, price float64 ) (int64 , error)

	// TODO
	// получить список стоп-заявок
	// создать новую стоп-заявку
	// отменить стоп-заявку

}

// все методы api это обертки над сервисами

// GetPortfolio получить данные по портфелю
func (c *Client) GetPortfolio(ctx context.Context, opts ...Option) (Portfolio, error) {
	p := &Options{
		IncludePositions: true,
	}
	for _, opt := range opts {
		opt(p)
	}
	s := c.NewGetPortfolioService().
		IncludePositions(p.IncludePositions).
		IncludeMoney(p.IncludeMoney).
		IncludeCurrencies(p.IncludeCurrencies).
		IncludeMaxBuySell(p.IncludeMaxBuySell)

	return s.Do(ctx)
}

// GetSecurity получить список инструментов (Максимальное Количество запросов в минуту = 1 )
func (c *Client) GetSecurity(ctx context.Context, board, symbol string) (Securities, error) {
	return c.NewSecurityService().Board(board).Symbol(symbol).Do(ctx)
}

// получить список заявок
func (c *Client) GetOrders(ctx context.Context, opts ...Option) ([]Order, error) {
	p := &Options{}
	for _, opt := range opts {
		opt(p)
	}
	s := c.NewGetOrderService().
		IncludeActive(p.IncludeActive).
		IncludeCanceled(p.IncludeActive).
		IncludeCanceled(p.IncludeCanceled)

	return s.Do(ctx)
}

// DeleteOrder удаление заявки
// clientId - торговый код клиента (обязательный)
// transactionId int64 (обязательный)
func (c *Client) DeleteOrder(ctx context.Context, transactionId int64) error {
	return c.NewCancelOrderService(transactionId).Do(ctx)
}

// GetCandles получить свечи
// Для запроса количества свечей в запросе необходимо указать count и либо from (начиная с указанной даты) либо to (до указанной даты).
// Для запроса за интервал необходимо указать from и to.
func (c *Client) GetCandles(ctx context.Context, board, symbol string, timeFrame TimeFrame, opts ...Option) ([]Candle, error) {
	p := &Options{}
	for _, opt := range opts {
		opt(p)
	}

	s := c.NewCandlesService(board, symbol, timeFrame).Count(p.Count)
	if p.EndTime != nil {
		s.endTime = p.EndTime
	}
	if p.StartTime != nil {
		s.startTime = p.StartTime
	}
	return s.Do(ctx)

}
