package message

import (
	"encoding/json"
	"fmt"
)

type Heartbeat struct {
	Event string `json:"event"`
}

//	{
//	    "event": "bts:heartbeat"
//	}
func NewHeartbeat() ([]byte, error) {
	b, err := json.Marshal(Heartbeat{"bts:heartbeat"})
	if err != nil {
		return nil, err
	}

	return b, nil
}

type Subscribe struct {
	Event string        `json:"event"`
	Data  SubscribeData `json:"data"`
}

//	{
//	    "event": "bts:subscribe",
//	    "data": {
//	        "channel": "[channel_name]-[user-id]",
//			"auth": "[token]",
//	    }
//	}
func NewSubscribe(data SubscribeData) Subscribe {
	return Subscribe{
		Event: "bts:subscribe",
		Data:  data,
	}
}

type SubscribeData struct {
	Channel string `json:"channel"`
	Auth    string `json:"auth"`
}

// "channel": "[channel_name]-[user-id]"
func NewSubscribeData(channelName, token string, userID int) SubscribeData {
	return SubscribeData{
		Channel: fmt.Sprintf("%s-%v", channelName, userID),
		Auth:    token,
	}
}

// Private My Orders
// channel: private-my_orders_[currency_pair]
func MyOrdersMessage(pair, token string, userID int) Subscribe {
	ch := fmt.Sprintf("private-my_orders_%s", pair)
	return NewSubscribe(NewSubscribeData(ch, token, userID))
}

// Private My Trades
// channel: private-my_trades_[currency_pair]
func MyTradesMessage(pair, token string, userID int) Subscribe {
	ch := fmt.Sprintf("private-my_trades_%s", pair)
	return NewSubscribe(NewSubscribeData(ch, token, userID))
}
