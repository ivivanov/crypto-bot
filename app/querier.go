package app

import (
	bs "github.com/ivivanov/crypto-bot/bitstamp"
	"github.com/ivivanov/crypto-bot/bitstamp/response"
	"github.com/ivivanov/crypto-bot/helper"
)

type PrivateGetter interface {
	PostAccountBalances(currencyPair string) (*[]response.AccountBalances, error)
	PostOpenOrders(currencyPair string) (*[]response.OpenOrder, error)
}

type PublicGetter interface {
	GetOHLC(currencyPair string, step, limit int) (*[]response.OHLC, error)
}

type Querier struct {
	privateGetter PrivateGetter
	publicGetter  PublicGetter
}

func NewQuerier(
	apiKey string,
	apiSecret string,
	customerID string,
	verbose bool,
) (*Querier, error) {
	apiConn, err := bs.NewAuthConn(apiKey, apiSecret, customerID, verbose)
	if err != nil {
		return nil, err
	}

	return &Querier{
		privateGetter: apiConn,
		publicGetter:  apiConn,
	}, nil
}

func (b *Querier) BalanceAll(currencyPair string) error {
	resp, err := b.privateGetter.PostAccountBalances(currencyPair)
	if err != nil {
		return err
	}

	helper.PrintIdent(resp)

	return nil
}

func (b *Querier) OHLC(pair string, step, limit int) (*[]response.OHLC, error) {
	resp, err := b.publicGetter.GetOHLC(pair, step, limit)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *Querier) SMA(pair string, step, limit, period int) ([]float64, error) {
	ohlc, err := b.publicGetter.GetOHLC(pair, step, limit)
	if err != nil {
		return nil, err
	}

	result := smaFrom(period, *ohlc, func(v response.OHLC) float64 { return v.Close })

	return result, nil
}

func smaFrom(length int, history []response.OHLC, getVal func(v response.OHLC) float64) []float64 {
	result := make([]float64, len(history))
	sum := float64(0)

	for i, ohlc := range history {
		count := i + 1
		sum += getVal(ohlc)

		if i >= length {
			sum -= getVal(history[i-length])
			count = length
		}

		result[i] = sum / float64(count)
	}

	return result
}
