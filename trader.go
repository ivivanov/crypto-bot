package main

import (
	"log"
	"math"

	"github.com/ivivanov/crypto-bot/response"
)

// https://docs.google.com/spreadsheets/d/1OqDC3QNYZoLhXN1BwxTOxbgOBTo7aTbWj87HgYxSLfQ/edit?usp=sharing
// The idea is to buy/sell @ breakeven price with 0.02% fee.
//
// 1) bought USDT:
// buy USDT		buy price	USD before fee	fee			total USD paid
// 50.00000		0.99944		49.97200		0.01499		49.98699
//
// create new sell order:
// sell USDT	sell price	USD base		fee			exp USD gain	profit
// 50.00000		1.0000396	50.00198		0.01499		49.98699		0.0%
//
// total USD paid = fee + USD before fee
// USD base = total USD paid + fee
// USD base with profit = USD base + USD base * profit
// sell price = USD base / sell USDT
//
// 2) sold USDT:
// sell USDT	sell price	USD before fee	fee			total USD gain
// 50.00000		0.99938		49.96900		0.01499		49.95401
//
// create new buy order:
// buy USDT		buy price	USD base		fee			exp USD paid	profit
// 50.00000		0.9987804	49.93902		0.01499		49.95401		0.0%
//
// total USD gain = -fee + USD before fee
// USD base = total USD gain - fee
// USD base with profit = USD base - USD base * profit
// buy price = USD base / buy USDT

type Trader struct {
	app    *App
	tradeC <-chan *response.MyTrade
}

func NewTrader(app *App, tradeC <-chan *response.MyTrade) *Trader {
	return &Trader{
		app:    app,
		tradeC: tradeC,
	}
}

// Must start in new routine
func (t *Trader) Start() {
	for {
		trade := <-t.tradeC

		resp, err := t.PostCounterTrade(trade)
		log.Printf("Counter trade: %#v", resp)

		if err != nil {
			log.Fatal(err)
		}

		// log.Print("Put to sleep trader for 30 sec")
		// time.Sleep(30 * time.Second)
	}
}

// Returns round number up to 5 decimal precision
func (t *Trader) CalculatePrice(price, amount, fee float64, isBuyPrice bool) float64 {
	if isBuyPrice { // start with buy USDT
		totalUsdPaid := price*amount + fee
		usdBaseSell := totalUsdPaid + fee
		usdBaseSellProfit := usdBaseSell + (usdBaseSell*t.app.profit)/100
		sellPrice := usdBaseSellProfit / amount

		return round5dec(sellPrice)
	}

	// else -> start with sell USDT
	totalUsdGain := price*amount - fee
	usdBaseBuy := totalUsdGain - fee
	usdBaseBuyProfit := usdBaseBuy - (usdBaseBuy*t.app.profit)/100
	buyPrice := usdBaseBuyProfit / amount

	return round5dec(buyPrice)
}

func (t *Trader) PostCounterTrade(trade *response.MyTrade) (interface{}, error) {
	var resp interface{}
	var err error

	switch trade.Data.Side {
	case "buy":
		sellPrice := t.CalculatePrice(trade.Price(), trade.Amount(), trade.Fee(), true)
		resp, err = t.app.ordersCreator.PostSellLimitOrder(t.app.pair, trade.Amount(), sellPrice)
	case "sell":
		buyPrice := t.CalculatePrice(trade.Price(), trade.Amount(), trade.Fee(), false)
		resp, err = t.app.ordersCreator.PostBuyLimitOrder(t.app.pair, trade.Amount(), buyPrice)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func round5dec(num float64) float64 {
	precision := 100_000
	return math.Round(num*float64(precision)) / float64(precision)
}
