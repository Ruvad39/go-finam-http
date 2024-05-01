package finam

import "time"

// параметры для запросов

type Option func(o *Options)

type Options struct {
	// для запроса портфеля
	IncludeCurrencies bool // запросить информацию по валютам портфеля;
	IncludeMoney      bool // запросить информацию по денежным позициям портфеля;
	IncludePositions  bool // запросить информацию по позициям портфеля;
	IncludeMaxBuySell bool // запросить информацию о максимальном доступном объеме на покупку/продажу.
	// для запроса ордеров
	IncludeMatched  bool // вернуть исполненные заявки;
	IncludeCanceled bool // вернуть отмененные заявки;
	IncludeActive   bool // вернуть активные заявки.
	// для запроса свечей
	Count     int
	StartTime *time.Time
	EndTime   *time.Time
}

func NewOptions() *Options {
	p := &Options{}
	return p
}

// (true) запросить информацию по позициям портфеля;
func WithIncludePositions(param bool) Option {
	return func(opts *Options) {
		opts.IncludePositions = param
	}
}

// (true) запросить информацию по валютам портфеля;
func WithIncludeCurrencies(param bool) Option {
	return func(opts *Options) {
		opts.IncludeCurrencies = param
	}
}

// (true) запросить информацию по денежным позициям портфеля;
func WithIncludeMoney(param bool) Option {
	return func(opts *Options) {
		opts.IncludeMoney = param
	}
}

// (true) запросить информацию о максимальном доступном объеме на покупку/продажу.
func WithIncludeMaxBuySell(param bool) Option {
	return func(opts *Options) {
		opts.IncludeMaxBuySell = param
	}
}

// (true)  вернуть исполненные заявки;
func WithIncludeMatched(param bool) Option {
	return func(opts *Options) {
		opts.IncludeMatched = param
	}
}

// (true) вернуть отмененные заявки;
func WithIncludeCanceled(param bool) Option {
	return func(opts *Options) {
		opts.IncludeCanceled = param
	}
}

// (true) IncludeActive  вернуть активные заявки.
func WithIncludeActive(param bool) Option {
	return func(opts *Options) {
		opts.IncludeActive = param
	}
}

// Limit Для запроса количества свечей
func WithCount(param int) Option {
	return func(opts *Options) {
		opts.Count = param
	}
}

// startTime  Для запроса  свечей: начальная дата
func WithStartTime(param time.Time) Option {
	return func(opts *Options) {
		opts.StartTime = &param
	}
}

// endTime  Для запроса  свечей: конечная дата
func WithEndTime(param time.Time) Option {
	return func(opts *Options) {
		opts.EndTime = &param
	}
}
