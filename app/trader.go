package app

import (
	bsresponse "github.com/ivivanov/crypto-bot/bitstamp/response"
	"github.com/ivivanov/crypto-bot/helper"

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
	bot    *Bot
	tradeC <-chan *response.MyTrade
}

func NewTrader(bot *Bot, tradeC <-chan *response.MyTrade) (*Trader, error) {
	return &Trader{
		bot:    bot,
		tradeC: tradeC,
	}, nil
}

// Must start in new routine
func (t *Trader) Start() {
	for {
		trade := <-t.tradeC

		switch trade.Data.Side {
		case "buy":
			_, err := t.PostSellCounterTrade(trade)
			helper.HandleFatalError(err)
		case "sell":
			_, err := t.PostBuyCounterTrade(trade)
			helper.HandleFatalError(err)
		}
	}
}

// Returns round number up to 5 decimal precision
// start with buy USDT
func (t *Trader) CalculateSellPrice(buyPrice, amount float64) float64 {
	fee := (amount * t.bot.maker) / 100
	totalUsdPaid := buyPrice*amount + fee
	usdBaseSell := totalUsdPaid + fee
	usdBaseSellProfit := usdBaseSell + (usdBaseSell*t.bot.profit)/100
	sellPrice := usdBaseSellProfit / amount

	return helper.Round5dec(sellPrice)
}

// Returns round number up to 5 decimal precision
// start with sell USDT
func (t *Trader) CalculateBuyPrice(sellPrice, amount float64) float64 {
	fee := (amount * t.bot.maker) / 100
	totalUsdGain := sellPrice*amount - fee
	usdBaseBuy := totalUsdGain - fee
	usdBaseBuyProfit := usdBaseBuy - (usdBaseBuy*t.bot.profit)/100
	buyPrice := usdBaseBuyProfit / amount

	return helper.Round5dec(buyPrice)
}

func (t *Trader) PostSellCounterTrade(trade *response.MyTrade) (*bsresponse.SellLimitOrder, error) {
	price, err := helper.GetPriceFrom(trade.Data.ClientOrderID)
	if err != nil {
		return nil, err
	}

	sellPrice := t.CalculateSellPrice(price, trade.Amount())
	resp, err := t.bot.limitOrdersCreator.PostSellLimitOrder(t.bot.pair, helper.GetClientOrderID(t.bot.account, sellPrice), trade.Amount(), sellPrice)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *Trader) PostBuyCounterTrade(trade *response.MyTrade) (*bsresponse.BuyLimitOrder, error) {
	price, err := helper.GetPriceFrom(trade.Data.ClientOrderID)
	if err != nil {
		return nil, err
	}

	buyPrice := t.CalculateBuyPrice(price, trade.Amount())
	resp, err := t.bot.limitOrdersCreator.PostBuyLimitOrder(t.bot.pair, helper.GetClientOrderID(t.bot.account, buyPrice), trade.Amount(), buyPrice)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// func (t *Trader) PostBuySmaTrade(trade *response.MyTrade) (*bsresponse.BuyLimitOrder, error) {

// }
