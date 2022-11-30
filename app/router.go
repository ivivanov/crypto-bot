package app

import (
	"encoding/json"
	"log"

	"github.com/ivivanov/crypto-bot/helper"
	"github.com/ivivanov/crypto-bot/response"
)

type Router struct {
	bot    *Bot
	tradeC chan<- *response.MyTrade
}

func NewRouter(bot *Bot, tradeC chan<- *response.MyTrade) (*Router, error) {
	return &Router{
		bot:    bot,
		tradeC: tradeC,
	}, nil
}

func (r *Router) Do(raw []byte) {
	baseMsg := response.BaseResponse{}
	err := json.Unmarshal(raw, &baseMsg)
	if err != nil {
		log.Println("unmarshal:", err)
		return
	}

	switch baseMsg.Event {
	case "bts:heartbeat":
		err = heartbeatHandler(raw)
	case "order_created", "order_changed", "order_deleted":
		err = myOrderHandler(raw, r.bot.account, baseMsg.Event)
	case "trade":
		err = myTradeHandler(raw, r.bot.account, r.bot.tradeC)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func heartbeatHandler(raw []byte) error {
	hb := response.Heartbeat{}

	err := json.Unmarshal(raw, &hb)
	if err != nil {
		return err
	}

	if hb.Data.Status != "success" {
		log.Print("pong not received")
	}

	return nil
}

func myTradeHandler(raw []byte, account string, tradeC chan<- *response.MyTrade) error {
	myTrade := &response.MyTrade{}

	err := json.Unmarshal(raw, myTrade)
	if err != nil {
		return err
	}

	tradeAcc, err := helper.GetAccountFrom(myTrade.Data.ClientOrderID)
	if err != nil {
		return nil
	}

	if tradeAcc != account {
		return nil
	}

	log.Printf("Trade-> %v: %v @ %v", myTrade.Data.Side, myTrade.Data.Amount, myTrade.Data.Price)

	tradeC <- myTrade

	return nil
}

func myOrderHandler(raw []byte, account, event string) error {
	myOrder := response.MyOrder{}

	err := json.Unmarshal(raw, &myOrder)
	if err != nil {
		return err
	}

	orderAcc, err := helper.GetAccountFrom(myOrder.Data.ClientOrderID)
	if err != nil {
		return nil
	}

	if orderAcc != account {
		return nil
	}

	log.Printf("Order->%s-> %s: %s @ %s", event, myOrder.GetOrderType(), myOrder.Data.AmountStr, myOrder.Data.PriceStr)

	return nil
}
