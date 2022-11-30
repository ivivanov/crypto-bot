package app

import (
	"log"

	"github.com/ivivanov/crypto-bot/helper"
)

type Preparer struct {
	bot *Bot
}

func NewPreparer(bot *Bot) (*Preparer, error) {
	return &Preparer{
		bot: bot,
	}, nil
}

func (p *Preparer) OpenBuyOrders(bank, price, grid, orders float64) error {
	currPrice := price
	amount := helper.Round2dec(bank / orders)

	for i := 0; i < int(orders); i++ {
		clientOrderID := helper.GetRnClientOrderID(p.bot.account)
		currPrice = helper.Round5dec(currPrice - currPrice*grid/100)

		// log.Println(i, id, p.bot.pair, amount, currPrice)
		resp, err := p.bot.ordersCreator.PostBuyLimitOrder(p.bot.pair, clientOrderID, amount, currPrice)
		if err != nil {
			return err
		}

		log.Printf("%#v", resp)
	}

	return nil
}
