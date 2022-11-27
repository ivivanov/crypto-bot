package response

import "strconv"

type BaseResponse struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
}

type Heartbeat struct {
	Data struct {
		Status string `json:"status"`
	} `json:"data"`
}

type MyTrade struct {
	Data struct {
		ID             int    `json:"id"`
		OrderID        int    `json:"order_id"`
		ClientOrderID  string `json:"client_order_id"`
		Amount         string `json:"amount"`
		Price          string `json:"price"`
		Fee            string `json:"fee"`
		Side           string `json:"side"`
		Microtimestamp string `json:"microtimestamp"`
	} `json:"data"`
}

func (mt *MyTrade) Price() float64 {
	price, _ := strconv.ParseFloat(mt.Data.Price, 64)

	return price
}

func (mt *MyTrade) Amount() float64 {
	amount, _ := strconv.ParseFloat(mt.Data.Amount, 64)

	return amount
}

func (mt *MyTrade) Fee() float64 {
	fee, _ := strconv.ParseFloat(mt.Data.Fee, 64)

	return fee
}

type MyOrder struct {
	Data struct {
		ID             int     `json:"id"`
		IdStr          string  `json:"id_str"`
		ClientOrderID  string  `json:"client_order_id"`
		Amount         float64 `json:"amount"`
		AmountStr      string  `json:"amount_str"`
		Price          float64 `json:"price"`
		PriceStr       string  `json:"price_str"`
		OrderType      int8    `json:"order_type"`
		Microtimestamp string  `json:"microtimestamp"`
	} `json:"data"`
}

// Order type (0 - buy, 1 - sell).
func (mo *MyOrder) GetOrderType() string {
	if mo.Data.OrderType == 0 {
		return "buy"
	}

	return "sell"
}
