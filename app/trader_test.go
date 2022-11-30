package app

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	bsresponse "github.com/ivivanov/crypto-bot/bitstamp/response"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/ivivanov/crypto-bot/response"
)

func TestBreakEvenSellPrice(t *testing.T) {
	trader := Trader{
		bot: &Bot{
			profit: 0.0,
		},
	}

	buyPrice := 0.99944
	amount := 50.00000
	fee := 0.01499
	expSellPrice := 1.0000396
	expSellPriceRounded := helper.Round5dec(expSellPrice)
	actualSellPrice := trader.CalculatePrice(buyPrice, amount, fee, true)

	if buyPrice > actualSellPrice {
		t.Error("must: buy price < sell price")
	}

	if expSellPriceRounded != actualSellPrice {
		t.Errorf("exp: %v actual: %v", expSellPriceRounded, actualSellPrice)
	}
}

func TestBreakEvenBuyPrice(t *testing.T) {
	trader := Trader{
		bot: &Bot{
			profit: 0.0,
		},
	}

	sellPrice := 0.99938
	amount := 50.00000
	fee := 0.01499
	expBuyPrice := 0.9987804
	expBuyPriceRounded := helper.Round5dec(expBuyPrice)
	actualBuyPrice := trader.CalculatePrice(sellPrice, amount, fee, false)

	if sellPrice < actualBuyPrice {
		t.Error("must: sell price > buy price")
	}

	if expBuyPriceRounded != actualBuyPrice {
		t.Errorf("exp: %v actual: %v", expBuyPriceRounded, actualBuyPrice)
	}
}

func TestSellPriceWithProfit(t *testing.T) {
	trader := Trader{
		bot: &Bot{
			profit: 0.01,
		},
	}

	buyPrice := 0.99982
	amount := 15.00000
	fee := 0.00300
	expSellPrice := 1.000320002
	expSellPriceRounded := helper.Round5dec(expSellPrice)
	actualSellPrice := trader.CalculatePrice(buyPrice, amount, fee, true)

	if buyPrice > actualSellPrice {
		t.Error("must: buy price < sell price")
	}

	if expSellPriceRounded != actualSellPrice {
		t.Errorf("exp: %v actual: %v", expSellPriceRounded, actualSellPrice)
	}
}

func TestBuyPriceWithProfit(t *testing.T) {
	trader := Trader{
		bot: &Bot{
			profit: 0.01,
		},
	}

	sellPrice := 0.99938
	amount := 50.00000
	fee := 0.01499
	expBuyPrice := 0.99868
	expBuyPriceRounded := helper.Round5dec(expBuyPrice)
	actualBuyPrice := trader.CalculatePrice(sellPrice, amount, fee, false)

	if sellPrice < actualBuyPrice {
		t.Error("must: sell price > buy price")
	}

	if expBuyPriceRounded != actualBuyPrice {
		t.Errorf("exp: %v actual: %v", expBuyPriceRounded, actualBuyPrice)
	}
}

func TestPostCounterTrade(t *testing.T) {
	trader := Trader{
		bot: &Bot{
			limitOrdersCreator: &OrdersCreatorMock{},
		},
	}

	myTradeJson := `
	{
		"data": {
			"amount": "15.00000",
			"price": "0.99982",
			"fee": "0.003",
			"side": "buy"
		}
	}
	`

	myTrade := &response.MyTrade{}
	err := json.Unmarshal([]byte(myTradeJson), myTrade)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := trader.PostCounterTrade(myTrade)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", resp)
	sellResp, ok := resp.(*bsresponse.SellLimitOrder)
	if !ok {
		t.Fatal(err)
	}

	expAmount := 15.00000
	if sellResp.Amount != expAmount {
		t.Errorf("exp: %v, actual: %v", expAmount, sellResp.Amount)
	}

	expPrice := 1.00022
	if sellResp.Price != expPrice {
		t.Errorf("exp: %v, actual: %v", expPrice, sellResp.Price)
	}
}

type OrdersCreatorMock struct {
}

func (ocm *OrdersCreatorMock) PostSellLimitOrder(currencyPair, clientOrderID string, amount float64, price float64) (*bsresponse.SellLimitOrder, error) {
	return &bsresponse.SellLimitOrder{
		Price:         price,
		Amount:        amount,
		Type:          1,
		ID:            0,
		DateTime:      bsresponse.DateTime(time.Now()),
		ClientOrderID: clientOrderID,
	}, nil
}

func (ocm *OrdersCreatorMock) PostBuyLimitOrder(currencyPair, clientOrderID string, amount float64, price float64) (*bsresponse.BuyLimitOrder, error) {
	return nil, nil
}
