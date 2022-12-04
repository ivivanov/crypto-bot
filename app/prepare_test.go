package app

import (
	"fmt"
	"testing"

	bsresponse "github.com/ivivanov/crypto-bot/bitstamp/response"
	"github.com/ivivanov/crypto-bot/helper"
)

type PrepareOrdersCreatorMock struct {
	i               int
	expAccount      string
	expCurrencyPair string
	expAmount       float64
	expPrices       map[int]float64
}

func (ocm *PrepareOrdersCreatorMock) PostSellLimitOrder(currencyPair, clientOrderID string, amount float64, price float64) (*bsresponse.SellLimitOrder, error) {
	return nil, nil
}

func (ocm *PrepareOrdersCreatorMock) PostBuyLimitOrder(currencyPair, clientOrderID string, amount float64, price float64) (*bsresponse.BuyLimitOrder, error) {
	if ocm.expPrices[ocm.i] != price {
		return nil, fmt.Errorf("exp: %v, act: %v", ocm.expPrices[ocm.i], price)
	}

	if ocm.expAmount != amount {
		return nil, fmt.Errorf("exp: %v, act: %v", ocm.expAmount, amount)
	}

	if ocm.expCurrencyPair != currencyPair {
		return nil, fmt.Errorf("exp: %v, act: %v", ocm.expCurrencyPair, currencyPair)
	}

	account, err := helper.GetAccountFrom(clientOrderID)
	if err != nil {
		return nil, err
	}

	if ocm.expAccount != account {
		return nil, fmt.Errorf("exp: %v, act: %v", ocm.expAccount, account)
	}

	ocm.i++

	return &bsresponse.BuyLimitOrder{
		Price:         price,
		Amount:        amount,
		ClientOrderID: clientOrderID,
	}, nil
}

func TestOpenBuyOrders(t *testing.T) {
	account := "test"
	pair := "testpair"
	price := 1.0
	bank := 500.0
	grid := 0.04
	orders := 6.0

	prepare := Preparer{
		bot: &Bot{
			pair:    pair,
			account: account,
			limitOrdersCreator: &PrepareOrdersCreatorMock{
				i:               0,
				expAccount:      account,
				expCurrencyPair: pair,
				expAmount:       83.33,
				expPrices: map[int]float64{
					0: 0.99960,
					1: 0.99920,
					2: 0.99880,
					3: 0.99840,
					4: 0.99800,
					5: 0.99760,
				},
			},
		},
	}

	err := prepare.OpenBuyOrders(bank, price, grid, orders)
	if err != nil {
		t.Fatal(err)
	}
}
