package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	bs "github.com/ivivanov/crypto-socks/bitstamp"
	bsresponse "github.com/ivivanov/crypto-socks/bitstamp/response"

	"github.com/ivivanov/crypto-socks/message"
	"github.com/ivivanov/crypto-socks/response"
)

type OrdersCreator interface {
	PostSellLimitOrder(currencyPair string, amount float64, price float64) (*bsresponse.SellLimitOrder, error)
	PostBuyLimitOrder(currencyPair string, amount float64, price float64) (*bsresponse.BuyLimitOrder, error)
}

type App struct {
	addr   string
	scheme string
	pair   string
	fee    float64

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

func NewApp() (*App, error) {
	addr := "ws.bitstamp.net"
	scheme := "wss"
	pair := "usdtusd"
	fee := 0.02

	secret, err := bs.GetSecret()
	if err != nil {
		return nil, err
	}

	wsUrl := url.URL{Scheme: scheme, Host: addr}

	apiConn, err := bs.NewAuthConn(secret.Key, secret.Secret, secret.CustomerID)
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

	app := &App{
		addr:   addr,
		scheme: scheme,
		pair:   pair,
		fee:    fee,

		interruptC: make(chan os.Signal, 1),
		messageC:   make(chan []byte),
		doneC:      make(chan struct{}),
		tradeC:     make(chan *response.MyTrade),

		wsConn:        wsConn,
		ordersCreator: apiConn,

		heartbeatMessage: heartbeatMessage,
	}

	app.trader = NewTrader(app, app.tradeC)

	myOrders := message.MyOrdersMessage(pair, wsToken.Token, wsToken.UserID)
	myOrdersMessage, err := json.Marshal(myOrders)
	if err != nil {
		return nil, err
	}

	app.myOrdersMessage = myOrdersMessage

	myTrades := message.MyTradesMessage(pair, wsToken.Token, wsToken.UserID)
	myTradesMessage, err := json.Marshal(myTrades)
	if err != nil {
		return nil, err
	}

	app.myTradesMessage = myTradesMessage

	return app, nil
}

func (app *App) CloseAll() {
	app.wsConn.Close()
	close(app.interruptC)
	close(app.messageC)
	close(app.doneC)
	close(app.tradeC)
}

func (app *App) Run() error {
	defer app.CloseAll()
	signal.Notify(app.interruptC, os.Interrupt)

	// routine: read
	go func() {
		for {
			_, msg, err := app.wsConn.ReadMessage()
			if err != nil {
				log.Fatal("read:", err)
				return
			}

			ResponseHandler(msg, app.tradeC)
		}
	}()

	// routine: ping/pong
	go func() {
		for {
			app.messageC <- app.heartbeatMessage
			time.Sleep(15 * time.Second)
		}
	}()

	// routine: subscribe
	go func() {
		app.messageC <- app.myOrdersMessage
		app.messageC <- app.myTradesMessage
	}()

	// routine: trader
	go func() {
		app.trader.Start()
	}()

	// routine: main
	// write, interupt, done
	for {
		select {
		case msg := <-app.messageC:
			err := app.wsConn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return fmt.Errorf("write: %w", err)
			}
		case <-app.interruptC:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := app.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return fmt.Errorf("write close: %w", err)
			}

			select {
			case <-app.doneC:
			case <-time.After(time.Second):
			}

			return nil
		case <-app.doneC:
			return nil
		}
	}
}
