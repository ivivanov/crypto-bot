package main

import (
	"testing"
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
			// apiConn: ,
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

// func TestPostNewTrade(t *testing.T) {
// 	trader := Trader{
// 		app: &App{
// 			fee:     0.02,
// 			apiConn: &ConnMock{},
// 		},
// 	}

// }
