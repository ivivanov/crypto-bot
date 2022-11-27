package main

import (
	"log"
	"math"
	"time"

	"github.com/ivivanov/crypto-socks/response"
)

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

		log.Print("Put to sleep trader for 30 sec")
		time.Sleep(30 * time.Second)
	}
}

// Returns round number up to 5 decimal precision
func (t *Trader) CalculateBreakevenPrice(price float64, isBuyPrice bool) float64 {
	precision := 100_000

	if isBuyPrice {
		sellPrice := price + price*t.app.fee/100
		return math.Round(sellPrice*float64(precision)) / float64(precision)
	}

	buyPrice := price - price*t.app.fee/100
	return math.Round(buyPrice*float64(precision)) / float64(precision)
}

func (t *Trader) PostCounterTrade(trade *response.MyTrade) (interface{}, error) {
	newAmount := trade.Amount() + trade.Fee()

	var resp interface{}
	var err error

	switch trade.Data.Side {
	case "buy":
		sellPrice := t.CalculateBreakevenPrice(trade.Price(), true)
		log.Print(t.app.pair, newAmount, sellPrice, sellPrice)
		resp, err = t.app.ordersCreator.PostSellLimitOrder(t.app.pair, newAmount, sellPrice)
	case "sell":
		buyPrice := t.CalculateBreakevenPrice(trade.Price(), false)
		log.Print(t.app.pair, newAmount, buyPrice, buyPrice)
		resp, err = t.app.ordersCreator.PostBuyLimitOrder(t.app.pair, newAmount, buyPrice)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}
