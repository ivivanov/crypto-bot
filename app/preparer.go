package app

import (
	"fmt"
	"log"

	bs "github.com/ivivanov/crypto-bot/bitstamp"
	bsre "github.com/ivivanov/crypto-bot/bitstamp/response"

	"github.com/ivivanov/crypto-bot/helper"
)

type OrderCreator interface {
	PostSellLimitOrder(currencyPair, clientOrderID string, amount, price float64) (*bsre.LimitOrder, error)
	PostBuyLimitOrder(currencyPair, clientOrderID string, amount, price float64) (*bsre.LimitOrder, error)
}

type Preparer struct {
	account            string
	pair               string
	limitOrdersCreator OrderCreator
}

func NewPreparer(
	account string,
	apiKey string,
	apiSecret string,
	customerID string,
	pair string,
	verbose bool,
) (*Preparer, error) {
	apiConn, err := bs.NewAuthConn(apiKey, apiSecret, customerID, verbose)
	if err != nil {
		return nil, err
	}

	return &Preparer{
		account:            account,
		pair:               pair,
		limitOrdersCreator: apiConn,
	}, nil
}

func (p *Preparer) OpenBuyOrders(bank, price, grid, orders float64) error {
	currPrice := price
	amount := helper.Round2dec(bank / orders)

	for i := 0; i < int(orders); i++ {
		currPrice = helper.Round5dec(currPrice - currPrice*grid/100)
		clientOrderID := helper.GetClientOrderID(p.account, currPrice)

		resp, err := p.limitOrdersCreator.PostBuyLimitOrder(p.pair, clientOrderID, amount, currPrice)
		if err != nil {
			return err
		}

		if resp.IsError() {
			log.Print(resp.Reason)
			return fmt.Errorf("failed to create order")
		}

		log.Printf("Order-created-> %v: %v @ %v", "buy", resp.Amount, resp.Price)
	}

	return nil
}

func (p *Preparer) OpenSellOrders(bank, price, grid, orders float64) error {
	currPrice := price
	amount := helper.Round2dec(bank / orders)

	for i := 0; i < int(orders); i++ {
		currPrice = helper.Round5dec(currPrice + currPrice*grid/100)
		clientOrderID := helper.GetClientOrderID(p.account, currPrice)

		resp, err := p.limitOrdersCreator.PostSellLimitOrder(p.pair, clientOrderID, amount, currPrice)
		if err != nil {
			return err
		}

		if resp.IsError() {
			log.Print(resp.Reason)
			return fmt.Errorf("failed to create order")
		}

		log.Printf("Order-created-> %v: %v @ %v", "sell", resp.Amount, resp.Price)
	}

	return nil
}

func (p *Preparer) OpenBuySellOrders(bank, price, grid, orders float64) error {
	if err := p.OpenBuyOrders(bank/2, price, grid, orders/2); err != nil {
		return err
	}

	if err := p.OpenSellOrders(bank/2, price, grid, orders/2); err != nil {
		return err
	}

	return nil
}
