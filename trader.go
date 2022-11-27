package main

import (
	"log"
	"math"
	"time"

	"github.com/ivivanov/crypto-socks/response"
)

type Trader struct {
	app    *App
	tradeC <-chan response.MyTrade
}

func NewTrader(app *App, tradeC <-chan response.MyTrade) *Trader {
	return &Trader{
		app:    app,
		tradeC: tradeC,
	}
}

// Must start in new routine
func (t *Trader) Start() {
	for {
		trade := <-t.tradeC
		t.PostNewTrade(trade)
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

func (t *Trader) PostNewTrade(trade response.MyTrade) {
	newAmount := trade.Amount() + trade.Fee()
	
	switch trade.Data.Side {
	case "buy":
		sellPrice := t.CalculateBreakevenPrice(trade.Price(), true)
		log.Print(t.app.pair, newAmount, sellPrice, sellPrice)
		resp, err := t.app.apiConn.PostSellLimitOrder(t.app.pair, newAmount, sellPrice, sellPrice, false)
		if err != nil {
			log.Fatal(err)
		}

		log.Print(resp)
	case "sell":
		buyPrice := t.CalculateBreakevenPrice(trade.Price(), false)
		log.Print(t.app.pair, newAmount, buyPrice, buyPrice)

		resp, err := t.app.apiConn.PostBuyLimitOrder(t.app.pair, newAmount, buyPrice, buyPrice, false)
		if err != nil {
			log.Fatal(err)
		}

		log.Print(resp)
	}
}
