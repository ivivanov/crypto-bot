package app

import (
	"encoding/json"
	"log"

	bs "github.com/ivivanov/crypto-bot/bitstamp"
	"github.com/ivivanov/crypto-bot/bitstamp/response"
)

type BalanceGetter interface {
	PostAccountBalances(currencyPair string) (*[]response.AccountBalances, error)
}

type Balancer struct {
	balanceGetter BalanceGetter
}

func NewBalancer(
	apiKey string,
	apiSecret string,
	customerID string,
	verbose bool,
) (*Balancer, error) {
	apiConn, err := bs.NewAuthConn(apiKey, apiSecret, customerID, verbose)
	if err != nil {
		return nil, err
	}

	return &Balancer{
		balanceGetter: apiConn,
	}, nil
}

func (b *Balancer) BalanceAll(currencyPair string) error {
	resp, err := b.balanceGetter.PostAccountBalances(currencyPair)
	if err != nil {
		return err
	}

	r, _ := json.MarshalIndent(resp, "", "	")
	log.Printf("%s", string(r))

	return nil
}
