package response

import (
	"encoding/json"
	"strconv"
	"testing"
)

func TestMyTrade(t *testing.T) {
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

	myTrade := &MyTrade{}
	err := json.Unmarshal([]byte(myTradeJson), myTrade)
	if err != nil {
		t.Error(err)
	}

	if myTrade.Data.ID != 256568732 {
		t.Error("id not matching")
	}

	var expAmount float64 = 15.00000
	var actAmount float64 = myTrade.Amount()
	if expAmount != actAmount {
		t.Errorf("exp: %v, actual: %v", expAmount, actAmount)
	}
}

func TestMyOrder(t *testing.T) {
	myOrderJson := `
	{
		"data": {
			"id": 1559367229386752,
			"id_str": "1559367229386752",
			"order_type": 0,
			"datetime": "1669539865",
			"microtimestamp": "1669539864644000",
			"amount": 15,
			"amount_str": "15.00000",
			"price": 0.99931,
			"price_str": "0.99931"
		},
		"channel": "private-my_orders_usdtusd-1106202",
		"event": "order_created"
	}
	`

	myOrder := &MyOrder{}
	err := json.Unmarshal([]byte(myOrderJson), myOrder)
	if err != nil {
		t.Error(err)
	}

	if myOrder.Data.ID != 1559367229386752 {
		t.Error("id not matching")
	}

	var expPrice float64 = 0.99931
	actPrice, _ := strconv.ParseFloat(myOrder.Data.PriceStr, 64)
	if expPrice != actPrice || expPrice != myOrder.Data.Price {
		t.Errorf("exp: %v, actual: %v", expPrice, myOrder.Data.Price)
	}
}
