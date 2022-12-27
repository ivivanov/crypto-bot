package traders

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ivivanov/crypto-bot/app"
	bsre "github.com/ivivanov/crypto-bot/bitstamp/response"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/ivivanov/crypto-bot/response"
	"golang.org/x/exp/slices"
)

//go:embed ohlc_test_data.json
var ohlcData string

const (
	OhclLen int = 18
)

type ApiSMAMock struct {
	openOrders []bsre.OpenOrder
	ohlc       []bsre.OHLC
	orderC     chan struct{}
}

func newApiSMAMock(orderC chan struct{}) *ApiSMAMock {
	ohlc := []bsre.OHLC{}

	err := json.Unmarshal([]byte(ohlcData), &ohlc)
	if err != nil {
		panic(err)
	}

	return &ApiSMAMock{
		openOrders: []bsre.OpenOrder{},
		ohlc:       ohlc,
		orderC:     orderC,
	}
}

func (a *ApiSMAMock) PostSellLimitOrder(currencyPair, clientOrderID string, amount, price float64) (*bsre.LimitOrder, error) {
	a.openOrders = append(a.openOrders, bsre.OpenOrder{
		Amount: amount,
		Price:  price,
		Type:   Sell,
	})

	a.orderC <- struct{}{}

	return nil, nil
}
func (a *ApiSMAMock) PostBuyLimitOrder(currencyPair, clientOrderID string, amount, price float64) (*bsre.LimitOrder, error) {
	a.openOrders = append(a.openOrders, bsre.OpenOrder{
		Amount: amount,
		Price:  price,
		Type:   Buy,
	})

	a.orderC <- struct{}{}

	return nil, nil
}

func (a *ApiSMAMock) PostOpenOrders(currencyPair string) (*[]bsre.OpenOrder, error) {
	return &a.openOrders, nil
}

func (a *ApiSMAMock) GetOHLC(currencyPair string, step, limit int) (*[]bsre.OHLC, error) {
	sub := a.ohlc[OhclLen-limit:]

	return &sub, nil
}

func (a *ApiSMAMock) PostCancelOrder(id int64) (*bsre.CancelOrder, error) {
	a.openOrders = slices.Delete(a.openOrders, 0, 1)

	return nil, nil
}

func Test_SMA_Polling(t *testing.T) {
	smaC := make(chan float64)

	config := &app.BotCtx{
		Account: "account",
		Pair:    "pair",
	}

	trader := &SMATrader{
		Ctx:       config,
		ApiSMA:    newApiSMAMock(nil),
		Timeframe: 1,
		OHLCLimit: 18,
		Length:    8,
		Offset:    0.0001,
		Profit:    0.0002,
		Amount:    100,
	}

	go trader.SMAPolling(smaC)

	sma := <-smaC
	exp := 1.0000700000000007
	if sma != exp {
		t.Fatalf("exp %v, act: %v", exp, sma)
	}

	sma = <-smaC
	exp = 1.0000700000000007
	if sma != exp {
		t.Fatalf("exp %v, act: %v", exp, sma)
	}
}

func Test_ReOrder_OnNewSMA(t *testing.T) {
	tradeC := make(chan *response.MyTrade)
	orderC := make(chan struct{})
	api := newApiSMAMock(orderC)

	config := &app.BotCtx{
		Account: "account",
		Pair:    "pair",
	}

	trader := &SMATrader{
		Ctx:       config,
		ApiSMA:    api,
		Timeframe: 1,
		OHLCLimit: 18,
		Length:    8,
		Offset:    0.0001,
		Profit:    0.0002,
		Amount:    100,
	}

	go trader.Start(tradeC)

	// 1st tick
	<-orderC

	orderLen := len(api.openOrders)
	if orderLen != 1 {
		t.Fatalf("exp %v, act: %v", 1, orderLen)
	}

	order := api.openOrders[0]
	expPrice := helper.Round5dec(0.99997) //(sma - trader.Offset, 1.0000700000000007 - 0.0001 = 0.99997000)
	if order.Price != expPrice {
		t.Fatalf("exp %v, act: %v", expPrice, order.Price)
	}

	// change SMA
	api.ohlc = append(api.ohlc, bsre.OHLC{Close: 1.00016})

	// 2nd tick
	<-orderC

	orderLen = len(api.openOrders)
	if orderLen != 1 {
		t.Fatalf("exp %v, act: %v", 1, orderLen)
	}

	order = api.openOrders[0]
	expPrice = helper.Round5dec(0.99999125) //(sma - trader.Offset, 1.0000912500000005 - 0.0001 = 0.99999125)
	if order.Price != expPrice {
		t.Fatalf("exp %v, act: %v", expPrice, order.Price)
	}
}

func Test_PostSellOrder_OnTrade(t *testing.T) {
	tradeC := make(chan *response.MyTrade)
	orderC := make(chan struct{})
	api := newApiSMAMock(orderC)

	config := &app.BotCtx{
		Account: "account",
		Pair:    "pair",
	}

	initialAmount := 1000.0
	trader := &SMATrader{
		Ctx:       config,
		ApiSMA:    api,
		Timeframe: 1,
		OHLCLimit: 18,
		Length:    8,
		Offset:    0.0001,
		Profit:    0.0002,
		Amount:    initialAmount,
	}

	go trader.Start(tradeC)

	// 1st tick - Initial order
	<-orderC

	orderLen := len(api.openOrders)
	if orderLen != 1 {
		t.Fatalf("exp %v, act: %v", 1, orderLen)
	}

	order := api.openOrders[0]
	expBuyPrice := helper.Round5dec(0.99997) //(sma - trader.Offset, 1.0000700000000007 - 0.0001 = 0.99997000)
	if order.Price != expBuyPrice {
		t.Fatalf("exp %v, act: %v", expBuyPrice, order.Price)
	}

	if order.Amount != initialAmount {
		t.Fatalf("exp %v, act: %v", initialAmount, order.Amount)
	}

	// make trade
	expAmount := 101.0
	tradeC <- createMyTrade("buy", expBuyPrice, expAmount)

	// 2nd tick
	<-orderC

	orderLen = len(api.openOrders)
	if orderLen != 2 {
		t.Fatalf("exp %v, act: %v", 2, orderLen)
	}

	order = api.openOrders[1]

	if order.Type != Sell {
		t.Fatalf("Must be sell order")
	}

	if order.Amount != expAmount {
		t.Fatalf("exp %v, act: %v", expAmount, order.Amount)
	}

	expSellPrice := helper.Round5dec(1.00017) //(buyPrice + trader.Profit, 0.99997 + 0.0002 = 1.00017)
	if order.Price != expSellPrice {
		t.Fatalf("exp %v, act: %v", expSellPrice, order.Price)
	}
}

func Test_PostBuyOrder_OnTrade(t *testing.T) {
	tradeC := make(chan *response.MyTrade)
	orderC := make(chan struct{})
	api := newApiSMAMock(orderC)

	config := &app.BotCtx{
		Account: "account",
		Pair:    "pair",
	}

	initialAmount := 1000.0
	trader := &SMATrader{
		Ctx:       config,
		ApiSMA:    api,
		Timeframe: 1,
		OHLCLimit: 18,
		Length:    8,
		Offset:    0.0001,
		Profit:    0.0002,
		Amount:    initialAmount,
	}

	go trader.Start(tradeC)

	// 1st tick - Initial order
	<-orderC

	// make trade
	expAmount := 101.0
	expSellPrice := helper.Round5dec(1.00017) //(buyPrice + trader.Profit, 0.99997 + 0.0002 = 1.00017)
	tradeC <- createMyTrade("sell", expSellPrice, expAmount)

	// 2nd tick
	<-orderC

	orderLen := len(api.openOrders)
	if orderLen != 2 {
		t.Fatalf("exp %v, act: %v", 2, orderLen)
	}

	order := api.openOrders[1]

	if order.Type != Buy {
		t.Fatalf("Must be buy order")
	}

	if order.Amount != expAmount {
		t.Fatalf("exp %v, act: %v", expAmount, order.Amount)
	}

	expBuyPrice := helper.Round5dec(0.99997) //(sma - trader.Offset, 1.0000700000000007 - 0.0001 = 0.99997000)
	if order.Price != expBuyPrice {
		t.Fatalf("exp %v, act: %v", expBuyPrice, order.Price)
	}
}

func createMyTrade(side string, price, amount float64) *response.MyTrade {
	return &response.MyTrade{
		Data: struct {
			ID             int    `json:"id"`
			OrderID        int    `json:"order_id"`
			ClientOrderID  string `json:"client_order_id"`
			Amount         string `json:"amount"`
			Price          string `json:"price"`
			Fee            string `json:"fee"`
			Side           string `json:"side"`
			Microtimestamp string `json:"microtimestamp"`
		}{
			Side:   side,
			Price:  fmt.Sprint(price),
			Amount: fmt.Sprint(amount),
		},
	}
}
