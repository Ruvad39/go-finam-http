package finam

// параметры для запросов

type Option func(o *Options)

// для запроса портфеля
// includeCurrencies - запросить информацию по валютам портфеля;
// includeMoney - запросить информацию по денежным позициям портфеля;
// includePositions - запросить информацию по позициям портфеля;
// includeMaxBuySell - запросить информацию о максимальном доступном объеме на покупку/продажу.

type Options struct {
    IncludeCurrencies    bool // запросить информацию по валютам портфеля;
    IncludeMoney         bool // запросить информацию по денежным позициям портфеля;
    IncludePositions     bool // запросить информацию по позициям портфеля;
    IncludeMaxBuySell    bool // запросить информацию о максимальном доступном объеме на покупку/продажу.

    IncludeMatched       bool // вернуть исполненные заявки;
    IncludeCanceled      bool // вернуть отмененные заявки;
    IncludeActive        bool // вернуть активные заявки.
}

func NewOptions() (*Options) {
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