package app

import (
	"log"

	bs "github.com/ivivanov/crypto-bot/bitstamp"
	"github.com/ivivanov/crypto-bot/bitstamp/response"
)

type OrderCanceller interface {
	PostCancelAllOrders(currencyPair string) (*response.CancelAllOrders, error)
}

type Canceler struct {
	pair           string
	orderCanceller OrderCanceller
}

func NewCanceler(
	apiKey string,
	apiSecret string,
	customerID string,
	pair string,
	verbose bool,
) (*Canceler, error) {
	apiConn, err := bs.NewAuthConn(apiKey, apiSecret, customerID, verbose)
	if err != nil {
		return nil, err
	}

	return &Canceler{
		pair:           pair,
		orderCanceller: apiConn,
	}, nil
}

func (c *Canceler) CancelAll() error {
	resp, err := c.orderCanceller.PostCancelAllOrders(c.pair)
	if err != nil {
		return err
	}

	log.Printf("%+v", resp)

	return nil
}
