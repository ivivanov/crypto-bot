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

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type LimitOrdersCreator interface {
	PostSellLimitOrder(currencyPair, clientOrderID string, amount float64, price float64) (*bsre.SellLimitOrder, error)
	PostBuyLimitOrder(currencyPair, clientOrderID string, amount float64, price float64) (*bsre.BuyLimitOrder, error)
}

type Bot struct {
	account string
	pair    string
	profit  float64
	maker   float64
	taker   float64

	wsConn *websocket.Conn

	// channels
	interruptC chan os.Signal
	messageC   chan []byte
	doneC      chan struct{}
	tradeC     chan *response.MyTrade

	// bitstamp
	limitOrdersCreator LimitOrdersCreator
	wsUrl              string

	// messages
	heartbeatMessage []byte
	myOrdersMessage  []byte
	myTradesMessage  []byte

	router *Router
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
	maker float64,
	taker float64,
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

	heartbeatMessage, err := message.NewHeartbeat()
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		account: account,
		pair:    pair,
		profit:  profit,
		maker:   maker,
		taker:   taker,

		interruptC: make(chan os.Signal, 1),
		messageC:   make(chan []byte),
		doneC:      make(chan struct{}),
		tradeC:     make(chan *response.MyTrade),

		limitOrdersCreator: apiConn,
		wsUrl:              wsUrl.String(),

		heartbeatMessage: heartbeatMessage,
	}
	bot.router, err = NewRouter(bot, bot.tradeC)
	if err != nil {
		return nil, err
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
	close(b.interruptC)
	close(b.messageC)
	close(b.doneC)
	close(b.tradeC)
}

func (b *Bot) Run() error {
	defer b.CloseAll()
	signal.Notify(b.interruptC, os.Interrupt)

	log.Printf("connecting to %s", b.wsUrl)
	wsConn, httpResp, err := websocket.DefaultDialer.Dial(b.wsUrl, nil)
	if err != nil {
		return fmt.Errorf("dial ws: %w", err)
	}

	b.wsConn = wsConn

	log.Printf("dial status: %v", httpResp.StatusCode)

	go b.readPump()

	// routine: ping/pong
	// go func() {
	// 	for {
	// 		b.messageC <- b.heartbeatMessage
	// 		time.Sleep(30 * time.Second)
	// 	}
	// }()

	// routine: subscribe
	go func() {
		b.messageC <- b.myOrdersMessage
		b.messageC <- b.myTradesMessage
	}()

	go b.trader.Start()

	// hold the main routine
	b.writePump()

	return nil
}

func (b *Bot) readPump() {
	defer b.wsConn.Close()

	b.wsConn.SetReadDeadline(time.Now().Add(pongWait))
	b.wsConn.SetPongHandler(func(string) error {
		log.Print("pong")
		b.wsConn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := b.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			log.Printf("read error: %v", err)
			break
		}

		b.router.Do(msg)
	}
}

func (b *Bot) writePump() {
	defer b.wsConn.Close()
	ticker := time.NewTicker(pingPeriod)

	for {
		select {
		case msg := <-b.messageC:
			b.wsConn.SetWriteDeadline(time.Now().Add(writeWait))

			err := b.wsConn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		case <-b.interruptC:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := b.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return
			}

			select {
			case <-b.doneC:
			case <-time.After(time.Second):
			}

			return
		case <-ticker.C:
			b.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			err := b.wsConn.WriteMessage(websocket.TextMessage, b.heartbeatMessage)
			if err != nil {
				log.Printf("write error: %v", err)
				return
			}
		case <-b.doneC:
			return
		}
	}
}
