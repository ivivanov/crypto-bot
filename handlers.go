package main

import (
	"encoding/json"
	"log"

	"github.com/ivivanov/crypto-bot/response"
)

func ResponseHandler(raw []byte, tradeC chan<- *response.MyTrade) {
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
		err = myOrderHandler(raw, baseMsg.Event)
	case "trade":
		err = myTradeHandler(raw, tradeC)
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

func myTradeHandler(raw []byte, tradeC chan<- *response.MyTrade) error {
	myTrade := &response.MyTrade{}

	err := json.Unmarshal(raw, myTrade)
	if err != nil {
		return err
	}

	log.Printf("Trade-> %v: %v @ %v", myTrade.Data.Side, myTrade.Data.Amount, myTrade.Data.Price)

	tradeC <- myTrade

	return nil
}

func myOrderHandler(raw []byte, event string) error {
	myOrder := response.MyOrder{}

	err := json.Unmarshal(raw, &myOrder)
	if err != nil {
		return err
	}

	log.Printf("Order->%s-> %s: %s @ %s", event, myOrder.GetOrderType(), myOrder.Data.AmountStr, myOrder.Data.PriceStr)

	return nil
}
