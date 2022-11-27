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
	"github.com/ivivanov/crypto-socks/bitstamp"
	bitstampResponse "github.com/ivivanov/crypto-socks/bitstamp/response"
	"github.com/ivivanov/crypto-socks/message"
	"github.com/ivivanov/crypto-socks/response"
)

type App struct {
	addr   string
	scheme string
	pair   string
	fee    float64

	// channels
	interrupt chan os.Signal
	messageC  chan []byte
	doneC     chan struct{}
	tradeC    chan response.MyTrade

	// bitstamp
	apiConn *bitstamp.Conn
	secret  *bitstamp.Secret
	wsToken *bitstampResponse.WebsocketToken
	wsConn  *websocket.Conn

	// messages
	heartbeatMessage  []byte
	myOrdersSubscribe []byte
	myTradesSubscribe []byte

	trader *Trader
}

func NewApp() (*App, error) {
	addr := "ws.bitstamp.net"
	scheme := "wss"
	pair := "usdtusd"
	fee := 0.02

	secret, err := bitstamp.GetSecret()
	if err != nil {
		return nil, err
	}

	wsUrl := url.URL{Scheme: scheme, Host: addr}

	apiConn, err := bitstamp.NewAuthConn(secret.Key, secret.Secret, secret.CustomerID)
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

	myOrders := message.MyOrdersMessage(pair, wsToken.Token, wsToken.UserID)
	myOrdersMessage, err := json.Marshal(myOrders)
	if err != nil {
		return nil, err
	}

	myTrades := message.MyTradesMessage(pair, wsToken.Token, wsToken.UserID)
	myTradesMessage, err := json.Marshal(myTrades)
	if err != nil {
		return nil, err
	}

	app := &App{
		addr:   addr,
		scheme: scheme,
		pair:   pair,
		fee:    fee,

		interrupt: make(chan os.Signal, 1),
		messageC:  make(chan []byte),
		doneC:     make(chan struct{}),
		tradeC:    make(chan response.MyTrade),

		apiConn: apiConn,
		secret:  secret,
		wsToken: wsToken,
		wsConn:  wsConn,

		heartbeatMessage:  heartbeatMessage,
		myOrdersSubscribe: myOrdersMessage,
		myTradesSubscribe: myTradesMessage,
	}

	app.trader = NewTrader(app, app.tradeC)

	return app, nil
}

func (app *App) CloseAll() {
	app.apiConn.Close()
	app.wsConn.Close()
	close(app.interrupt)
	close(app.messageC)
	close(app.doneC)
	close(app.tradeC)
}

func (app *App) Run() error {
	defer app.CloseAll()
	signal.Notify(app.interrupt, os.Interrupt)

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
		app.messageC <- app.myOrdersSubscribe
		app.messageC <- app.myTradesSubscribe
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
		case <-app.interrupt:
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
