package traders

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ivivanov/crypto-bot/app"
	bsre "github.com/ivivanov/crypto-bot/bitstamp/response"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/ivivanov/crypto-bot/response"
)

type OrderCreatorMock struct {
}

func (ocm *OrderCreatorMock) PostSellLimitOrder(currencyPair, clientOrderID string, amount, price float64) (*bsre.LimitOrder, error) {
	return &bsre.LimitOrder{
		Price:         price,
		Amount:        amount,
		Type:          1,
		ID:            0,
		DateTime:      bsre.DateTime(time.Now()),
		ClientOrderID: clientOrderID,
	}, nil
}

func (ocm *OrderCreatorMock) PostBuyLimitOrder(currencyPair, clientOrderID string, amount, price float64) (*bsre.LimitOrder, error) {
	return &bsre.LimitOrder{
		Price:         price,
		Amount:        amount,
		Type:          0,
		ID:            0,
		DateTime:      bsre.DateTime(time.Now()),
		ClientOrderID: clientOrderID,
	}, nil
}

func TestCalculateSellPrice(t *testing.T) {
	for _, tc := range []struct {
		title        string
		profit       float64
		buyPrice     float64
		amount       float64
		expSellPrice float64
		makerFee     float64
	}{
		{
			title:        "Breakeven sell price",
			profit:       0.0,
			buyPrice:     0.99944,
			amount:       50.00000,
			expSellPrice: helper.Round5dec(0.99984),
			makerFee:     0.02,
		},
		{
			title:        "Breakeven sell price bitstamp maker fee",
			profit:       0.0,
			buyPrice:     1.0,
			amount:       50.00000,
			expSellPrice: helper.Round5dec(1.00040),
			makerFee:     0.02,
		},
		{
			title:        "Breakeven sell price kraken maker fee",
			profit:       0.0,
			buyPrice:     1.0,
			amount:       50.00000,
			expSellPrice: helper.Round5dec(1.00032),
			makerFee:     0.016,
		},
		{
			title:        "Sell price with profit",
			profit:       0.01,
			buyPrice:     0.99982,
			amount:       15.00000,
			expSellPrice: helper.Round5dec(1.000320002),
			makerFee:     0.02,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			trader := GridTrader{
				MakerFee: tc.makerFee,
				Profit:   tc.profit,
			}

			actualSellPrice := trader.CalculateSellPrice(tc.buyPrice, tc.amount)

			if tc.buyPrice > actualSellPrice {
				t.Error("must: buy price < sell price")
			}

			if tc.expSellPrice != actualSellPrice {
				t.Errorf("exp: %v actual: %v", tc.expSellPrice, actualSellPrice)
			}
		})
	}
}

func TestCalculateBuyPrice(t *testing.T) {
	for _, tc := range []struct {
		title       string
		profit      float64
		sellPrice   float64
		amount      float64
		expBuyPrice float64
	}{
		{
			title:       "Breakeven buy price",
			profit:      0.0,
			sellPrice:   0.99938,
			amount:      50.00000,
			expBuyPrice: helper.Round5dec(0.99898),
		},
		{
			title:       "Buy price with profit",
			profit:      0.01,
			sellPrice:   0.99982,
			amount:      15.00000,
			expBuyPrice: helper.Round5dec(0.99932),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			trader := GridTrader{
				Profit:   tc.profit,
				MakerFee: 0.02,
			}

			actualBuyPrice := trader.CalculateBuyPrice(tc.sellPrice, tc.amount)

			if tc.sellPrice < actualBuyPrice {
				t.Error("must: sell price > buy price")
			}

			if tc.expBuyPrice != actualBuyPrice {
				t.Errorf("exp: %v actual: %v", tc.expBuyPrice, actualBuyPrice)
			}
		})
	}
}

func TestPostSellCounterTrade(t *testing.T) {
	for _, tc := range []struct {
		title     string
		tradeJson string
		expAmount float64
		expPrice  float64
	}{
		{
			title: "Regular trade. Maker fee applied",
			tradeJson: `
			{
				"data": {
					"amount": "15.00000",
					"price": "0.99982",
					"fee": "0.003",
					"side": "buy",
					"client_order_id": "account_0.99982_1298498081"
				}
			}
			`,
			expAmount: 15.00000,
			expPrice:  1.00022,
		},
		{
			title: "Sell order hit directly bid. Taker fee applied",
			tradeJson: `
			{
				"data": {
					"amount": "15.00000",
					"price": "1.97000",
					"fee": "0.0045",
					"side": "buy",
					"client_order_id": "account_0.99982_1298498081"
				}
			}
			`,
			expAmount: 15.00000,
			expPrice:  1.00022,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			trader := GridTrader{
				Ctx:          &app.BotCtx{},
				MakerFee:     0.02,
				OrderCreator: &OrderCreatorMock{},
			}

			myBuyTrade := &response.MyTrade{}
			err := json.Unmarshal([]byte(tc.tradeJson), myBuyTrade)
			if err != nil {
				t.Fatal(err)
			}

			sell, err := trader.PostSellCounterTrade(myBuyTrade)
			if err != nil {
				t.Fatal(err)
			}

			if sell.Amount != tc.expAmount {
				t.Errorf("exp: %v, actual: %v", tc.expAmount, sell.Amount)
			}

			if sell.Price != tc.expPrice {
				t.Errorf("exp: %v, actual: %v", tc.expPrice, sell.Price)
			}
		})
	}
}

func TestPostBuyCounterTrade(t *testing.T) {
	for _, tc := range []struct {
		title     string
		tradeJson string
		expAmount float64
		expPrice  float64
	}{
		{
			title: "Regular trade. Maker fee applied",
			tradeJson: `
			{
				"data": {
					"amount": "50.00000",
					"price": "0.99938",
					"fee": "0.01",
					"side": "sell",
					"client_order_id": "account_0.99938_1298498081"
				}
			}
			`,
			expAmount: 50.00000,
			expPrice:  0.99898,
		},
		{
			title: "Buy order hit directly ask. Taker fee applied",
			tradeJson: `
			{
				"data": {
					"amount": "50.00000",
					"price": "2.99938",
					"fee": "0.01000",
					"side": "sell",
					"client_order_id": "account_0.99954_1298498081"
				}
			}
			`,
			expAmount: 50.00000,
			expPrice:  0.99914,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			trader := GridTrader{
				Ctx:          &app.BotCtx{},
				MakerFee:     0.02,
				OrderCreator: &OrderCreatorMock{},
			}

			mySellTrade := &response.MyTrade{}
			err := json.Unmarshal([]byte(tc.tradeJson), mySellTrade)
			if err != nil {
				t.Fatal(err)
			}

			buy, err := trader.PostBuyCounterTrade(mySellTrade)
			if err != nil {
				t.Fatal(err)
			}

			if buy.Amount != tc.expAmount {
				t.Errorf("exp: %v, actual: %v", tc.expAmount, buy.Amount)
			}

			if buy.Price != tc.expPrice {
				t.Errorf("exp: %v, actual: %v", tc.expPrice, buy.Price)
			}
		})
	}
}
