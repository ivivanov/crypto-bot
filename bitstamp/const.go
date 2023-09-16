package bitstamp

import "errors"

var (
	ErrNotFound     = errors.New("Not Found")
	ErrAuthRequired = errors.New("Authentication Required")
)

// todo remove
var CurrencyPairs = []string{
	"btcusd",
	"btceur",
	"eurusd",
	"xrpusd",
	"xrpeur",
	"xrpbtc",
	"ltcusd",
	"ltceur",
	"ltcbtc",
	"ethusd",
	"etheur",
	"ethbtc",
}
