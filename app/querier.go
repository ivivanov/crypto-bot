package app

import (
	"encoding/json"
	"log"

	bs "github.com/ivivanov/crypto-bot/bitstamp"
	"github.com/ivivanov/crypto-bot/bitstamp/response"
)

type PrivateGetter interface {
	PostAccountBalances(currencyPair string) (*[]response.AccountBalances, error)
}

type PublicGetter interface {
	GetOHLC(currencyPair string, step, limit int) (interface{}, error)
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
