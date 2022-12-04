package app

import (
	"fmt"
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
		currPrice = helper.Round5dec(currPrice - currPrice*grid/100)
		clientOrderID := helper.GetClientOrderID(p.bot.account, currPrice)

		resp, err := p.bot.limitOrdersCreator.PostBuyLimitOrder(p.bot.pair, clientOrderID, amount, currPrice)
		if err != nil {
			return err
		}

		if resp.IsError() {
			log.Print(resp.Reason)
			return fmt.Errorf("failed to create order")
		}

		log.Printf("Order-created-> %v: %v @ %v", "buy", resp.Amount, resp.Price)
	}

	return nil
}
