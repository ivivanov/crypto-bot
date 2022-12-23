package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	"github.com/ivivanov/crypto-bot/bitstamp"
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

type Trader interface {
	Start(tradeC <-chan *response.MyTrade)
}

type BotCtx struct {
	ApiConn    *bitstamp.Conn
	Debug      bool
	Account    string
	Pair       string
	WSScheme   string
	WSAddr     string
	APIKey     string
	APISecret  string
	CustomerID string
}

type Bot struct {
	ctx    *BotCtx
	wsConn *websocket.Conn

	// channels
	interruptC chan os.Signal
	messageC   chan []byte
	doneC      chan struct{}
	tradeC     chan *response.MyTrade

	// bitstamp
	publicGetter PublicGetter
	wsUrl        string

	// messages
	heartbeatMessage []byte
	myOrdersMessage  []byte
	myTradesMessage  []byte

	router *Router
	trader Trader
}

func NewBot(ctx *BotCtx, trader Trader) (*Bot, error) {
	if ctx == nil || trader == nil {
		return nil, errors.New("all args are required")
	}
	wsUrl := url.URL{Scheme: ctx.WSScheme, Host: ctx.WSAddr}

	wsToken, err := ctx.ApiConn.WebsocketToken()
	if err != nil {
		return nil, err
	}

	heartbeatMessage, err := message.NewHeartbeat()
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		ctx:        ctx,
		interruptC: make(chan os.Signal, 1),
		messageC:   make(chan []byte),
		doneC:      make(chan struct{}),
		tradeC:     make(chan *response.MyTrade),

		publicGetter: ctx.ApiConn,
		wsUrl:        wsUrl.String(),

		heartbeatMessage: heartbeatMessage,

		trader: trader,
	}

	bot.router, err = NewRouter(bot, bot.tradeC)
	if err != nil {
		return nil, err
	}

	myOrders := message.MyOrdersMessage(ctx.Pair, wsToken.Token, wsToken.UserID)
	myOrdersMessage, err := json.Marshal(myOrders)
	if err != nil {
		return nil, err
	}

	bot.myOrdersMessage = myOrdersMessage

	myTrades := message.MyTradesMessage(ctx.Pair, wsToken.Token, wsToken.UserID)
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

	// routine: subscribe
	go func() {
		b.messageC <- b.myOrdersMessage
		b.messageC <- b.myTradesMessage
	}()

	go b.trader.Start(b.tradeC)

	// hold the main routine
	b.writePump()

	return nil
}

func (b *Bot) readPump() {
	defer b.wsConn.Close()

	b.wsConn.SetReadDeadline(time.Now().Add(pongWait))
	b.wsConn.SetPongHandler(func(string) error {
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
