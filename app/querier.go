package app

import (
	"encoding/json"
	"log"
	"time"

	bs "github.com/ivivanov/crypto-bot/bitstamp"
	"github.com/ivivanov/crypto-bot/bitstamp/response"
)

type PrivateGetter interface {
	PostAccountBalances(currencyPair string) (*[]response.AccountBalances, error)
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

	r, _ := json.MarshalIndent(resp, "", "	")
	log.Printf("%s", string(r))

	return nil
}

func (b *Querier) OHLC(pair string, step, limit int) error {
	resp, err := b.publicGetter.GetOHLC(pair, step, limit)
	if err != nil {
		return err
	}

	r, _ := json.MarshalIndent(resp, "", "	")
	log.Printf("%s", string(r))

	return nil
}

func (b *Querier) SMA(pair string, step, limit, period int) error {
	ohlc, err := b.publicGetter.GetOHLC(pair, step, limit)
	if err != nil {
		return err
	}

	result := smaFrom(period, *ohlc, func(v response.OHLC) float64 { return v.Close })

	count := len(result)
	hourNow := time.Now().Hour()

	log.Printf("%v:00 %v, %v:00 %v, %v:00 %v", hourNow-2, result[count-3], hourNow-1, result[count-2], hourNow, result[count-1])

	return nil
}

func smaFrom(period int, history []response.OHLC, getVal func(v response.OHLC) float64) []float64 {
	result := make([]float64, len(history))
	sum := float64(0)

	for i, ohlc := range history {
		count := i + 1
		sum += getVal(ohlc)

		if i >= period {
			sum -= getVal(history[i-period])
			count = period
		}

		result[i] = sum / float64(count)
	}

	return result
}
