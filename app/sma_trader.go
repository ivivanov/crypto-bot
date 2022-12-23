package app

import (
	"log"
	"time"

	bsre "github.com/ivivanov/crypto-bot/bitstamp/response"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/ivivanov/crypto-bot/response"
)

type SMATrader struct {
	Ctx            *BotCtx
	OrderCreator   OrderCreator
	OrderCanceller OrderCanceller
	PrivateGetter  PrivateGetter
	PublicGetter   PublicGetter

	Timeframe int
	OHLCLimit int
	Length    int
	Offset    float64
	Profit    float64
	Amount    float64
}

func (st *SMATrader) Start(tradeC <-chan *response.MyTrade) {
	smaC := make(chan float64)
	defer close(smaC)

	go st.smaPolling(smaC)
	sma := 0.0

	for {
		select {
		case _ = <-tradeC:
			//	buy trade
			//		create sell apply fixed profit
			//	sell trade (1 cycle done)
			//		create buy order (with current SMA)

			// switch trade.Data.Side {
			// case "buy":
			// 	_, err := gt.PostSellCounterTrade(trade)
			// 	helper.HandleFatalError(err)
			// case "sell":
			// 	_, err := gt.PostBuyCounterTrade(trade)
			// 	helper.HandleFatalError(err)
			// }
		case newSMA := <-smaC:
			//	get open orders
			//	no open orders
			//		create buy order (with new SMA)
			// 	loop orders and check if limit price != new SMA
			//		re-create buy order (with new SMA)
			
			log.Print(newSMA)
			sma = newSMA

			openOrders, err := st.PrivateGetter.PostOpenOrders(st.Ctx.Pair)
			helper.HandleFatalError(err)
			log.Print(openOrders)

			buyPrice := helper.Round5dec(sma - st.Offset)

			if len(*openOrders) == 0 {
				log.Print("Initial order")
				// r, err := st.OrderCreator.PostBuyLimitOrder(st.Ctx.Pair, helper.GetClientOrderID(st.Ctx.Account, buyPrice), st.Amount, buyPrice)
				// helper.HandleFatalError(err)
				// helper.PrintIdent(r)
			} else if newSMA != sma {
				// TODO: if there are multiple orders on the same price aggregate them into 1 order
				for _, order := range *openOrders {
					if order.Price != newSMA {
						cancel, err := st.OrderCanceller.PostCancelOrder(order.ID)
						helper.HandleFatalError(err)
						helper.PrintIdent(cancel)
						order, err := st.OrderCreator.PostBuyLimitOrder(st.Ctx.Pair, helper.GetClientOrderID(st.Ctx.Account, buyPrice), st.Amount, buyPrice)
						helper.HandleFatalError(err)
						helper.PrintIdent(order)
					}
				}
			}
		}
	}
}

func (st *SMATrader) smaPolling(smaC chan<- float64) {
	updateInterval := time.Duration(st.Timeframe) * time.Second
	ticker := time.NewTicker(updateInterval)

	smaGetter := func() {
		ohlc, err := st.PublicGetter.GetOHLC(st.Ctx.Pair, st.Timeframe, st.OHLCLimit)
		helper.HandleFatalError(err)
		sma := smaFrom(st.Length, *ohlc, func(v bsre.OHLC) float64 { return v.Close })

		smaC <- sma[len(sma)-1]
	}

	smaGetter()

	for range ticker.C {
		smaGetter()
	}
}
