package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	bs "github.com/ivivanov/crypto-bot/bitstamp"
	bsre "github.com/ivivanov/crypto-bot/bitstamp/response"

	"github.com/ivivanov/crypto-bot/message"
	"github.com/ivivanov/crypto-bot/response"
)

type OrdersCreator interface {
	PostSellLimitOrder(currencyPair, clientOrderID string, amount float64, price float64) (*bsre.SellLimitOrder, error)
	PostBuyLimitOrder(currencyPair, clientOrderID string, amount float64, price float64) (*bsre.BuyLimitOrder, error)
}

type Bot struct {
	account string
	pair    string
	profit  float64

	// channels
	interruptC chan os.Signal
	messageC   chan []byte
	doneC      chan struct{}
	tradeC     chan *response.MyTrade

	// bitstamp
	ordersCreator OrdersCreator
	wsConn        *websocket.Conn

	// messages
	heartbeatMessage []byte
	myOrdersMessage  []byte
	myTradesMessage  []byte

	trader *Trader
}

func NewBot(
	account string,
	wsScheme string,
	wsAddr string,
	apiKey string,
	apiSecret string,
	customerID string,
	pair string,
	profit float64,
) (*Bot, error) {

	wsUrl := url.URL{Scheme: wsScheme, Host: wsAddr}

	apiConn, err := bs.NewAuthConn(apiKey, apiSecret, customerID)
	if err != nil {
		return nil, err
	}

	wsToken, err := apiConn.WebsocketToken()
	if err != nil {
		return nil, err
	}

	log.Printf("connecting to %s", wsUrl.String())
	wsConn, httpResp, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("dial ws: %w", err)
	}

	log.Printf("dial status: %v", httpResp.StatusCode)

	heartbeatMessage, err := message.NewHeartbeat()
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		account: account,
		pair:    pair,
		profit:  profit,

		interruptC: make(chan os.Signal, 1),
		messageC:   make(chan []byte),
		doneC:      make(chan struct{}),
		tradeC:     make(chan *response.MyTrade),

		wsConn:        wsConn,
		ordersCreator: apiConn,

		heartbeatMessage: heartbeatMessage,
	}

	bot.trader, err = NewTrader(bot, bot.tradeC)
	if err != nil {
		return nil, err
	}

	myOrders := message.MyOrdersMessage(bot.pair, wsToken.Token, wsToken.UserID)
	myOrdersMessage, err := json.Marshal(myOrders)
	if err != nil {
		return nil, err
	}

	bot.myOrdersMessage = myOrdersMessage

	myTrades := message.MyTradesMessage(bot.pair, wsToken.Token, wsToken.UserID)
	myTradesMessage, err := json.Marshal(myTrades)
	if err != nil {
		return nil, err
	}

	bot.myTradesMessage = myTradesMessage

	return bot, nil
}

func (b *Bot) CloseAll() {
	b.wsConn.Close()
	close(b.interruptC)
	close(b.messageC)
	close(b.doneC)
	close(b.tradeC)
}

func (b *Bot) Run() error {
	defer b.CloseAll()
	signal.Notify(b.interruptC, os.Interrupt)

	// routine: read
	go func() {
		for {
			_, msg, err := b.wsConn.ReadMessage()
			if err != nil {
				log.Fatal("read:", err)
				return
			}

			ResponseHandler(msg, b.tradeC)
		}
	}()

	// routine: ping/pong
	go func() {
		for {
			b.messageC <- b.heartbeatMessage
			time.Sleep(15 * time.Second)
		}
	}()

	// routine: subscribe
	go func() {
		b.messageC <- b.myOrdersMessage
		b.messageC <- b.myTradesMessage
	}()

	// routine: trader
	go func() {
		b.trader.Start()
	}()

	// routine: main
	// write, interupt, done
	for {
		select {
		case msg := <-b.messageC:
			err := b.wsConn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return fmt.Errorf("write: %w", err)
			}
		case <-b.interruptC:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := b.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return fmt.Errorf("write close: %w", err)
			}

			select {
			case <-b.doneC:
			case <-time.After(time.Second):
			}

			return nil
		case <-b.doneC:
			return nil
		}
	}
}
