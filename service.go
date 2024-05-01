package finam

// методы создания нужных сервисов по api

// NewGetPortfolioService init portfоlio service
func (c *Client) NewGetPortfolioService() *GetPortfolioService {
	return &GetPortfolioService{c: c}
}

// NewCandlesService init candles service
// обязательные параметры: board, symbol, timeFrame
func (c *Client) NewCandlesService(board, symbol string, timeFrame TimeFrame) *CandlesService {
	return &CandlesService{
		c:         c,
		board:     board,
		symbol:    symbol,
		timeFrame: timeFrame,
	}
}

// NewSecurityService init security service
func (c *Client) NewSecurityService() *SecurityService {
	return &SecurityService{c: c}
}

// NewGetOrderService init GetOrder Service
func (c *Client) NewGetOrderService() *GetOrderService {
	return &GetOrderService{c: c,
		includeActive: true, // сразу проставим "Вернуть активные заявки"
	}
}

// NewCancelOrderService init CancelOrder Service
// обязательные параметры: id
func (c *Client) NewCancelOrderService(transactionId int64) *CancelOrderService {
	return &CancelOrderService{
		c:             c,
		transactionId: transactionId,
	}
}
