package main

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	bsresponse "github.com/ivivanov/crypto-socks/bitstamp/response"
	"github.com/ivivanov/crypto-socks/response"
)

func TestBreakEvenSellPrice(t *testing.T) {
	trader := Trader{
		app: &App{
			fee: 0.02,
		},
	}
	buyPrice := 0.99954
	expSellPrice := 0.99974
	actualSellPrice := trader.CalculateBreakevenPrice(buyPrice, true)

	if buyPrice > actualSellPrice {
		t.Error("must: buy price < sell price")
	}

	if expSellPrice != actualSellPrice {
		t.Errorf("exp: %v actual: %v", expSellPrice, actualSellPrice)
	}
}

func TestBreakEvenBuyPrice(t *testing.T) {
	trader := Trader{
		app: &App{
			fee: 0.02,
		},
	}
	sellPrice := 0.99954
	expBuyPrice := 0.99934
	actualBuyPrice := trader.CalculateBreakevenPrice(sellPrice, false)

	if sellPrice < actualBuyPrice {
		t.Error("must: sell price > buy price")
	}

	if expBuyPrice != actualBuyPrice {
		t.Errorf("exp: %v actual: %v", expBuyPrice, actualBuyPrice)
	}
}

func TestPostCounterTrade(t *testing.T) {
	trader := Trader{
		app: &App{
			fee:           0.02,
			ordersCreator: &OrdersCreatorMock{},
		},
	}

	myTradeJson := `
	{
		"data": {
			"id": 256568732,
			"amount": "15.00000",
			"price": "0.99982",
			"microtimestamp": "1669538134992000",
			"fee": "0.003",
			"order_id": 1559359938404352,
			"side": "buy"
		},
		"channel": "private-my_trades_usdtusd-1106202",
		"event": "trade"
	}
	`

	myTrade := &response.MyTrade{}
	err := json.Unmarshal([]byte(myTradeJson), myTrade)
	if err != nil {
		t.Fatal(err)
	}

	resp, _ := trader.PostCounterTrade(myTrade)

	log.Printf("%#v", resp)
	sellResp, ok := resp.(*bsresponse.SellLimitOrder)
	if !ok {
		t.Fatal(err)
	}

	expAmount := myTrade.Amount() + myTrade.Fee()
	if sellResp.Amount != expAmount {
		t.Errorf("exp: %v, actual: %v", expAmount, sellResp.Amount)
	}
}

type OrdersCreatorMock struct {
}

func (ocm *OrdersCreatorMock) PostSellLimitOrder(currencyPair string, amount float64, price float64) (*bsresponse.SellLimitOrder, error) {
	return &bsresponse.SellLimitOrder{
		Price:    price,
		Amount:   amount,
		Type:     1,
		ID:       0,
		DateTime: bsresponse.DateTime(time.Now()),
	}, nil
}
func (ocm *OrdersCreatorMock) PostBuyLimitOrder(currencyPair string, amount float64, price float64) (*bsresponse.BuyLimitOrder, error) {
	return nil, nil
}
