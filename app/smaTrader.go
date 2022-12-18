package app

import (
	"github.com/ivivanov/crypto-bot/response"
)

type SMATrader struct {
	Config *BotCtx
	Order  OrderCreator

	smaC chan<- float64
}

func (st *SMATrader) Start(tradeC <-chan *response.MyTrade) {
	// go st.smaPolling()

	for {
		_ = <-tradeC

		// switch trade.Data.Side {
		// case "buy":
		// 	_, err := gt.PostSellCounterTrade(trade)
		// 	helper.HandleFatalError(err)
		// case "sell":
		// 	_, err := gt.PostBuyCounterTrade(trade)
		// 	helper.HandleFatalError(err)
		// }
	}
}

// func (t *Trader) PostBuySmaTrade(trade *response.MyTrade) (*bsresponse.BuyLimitOrder, error) {

// }

// func (st *SMATrader) smaPolling() {
// 	ticker := time.NewTicker(time.Duration(b.ctx.Timeframe))

// 	smaGetter := func() {
// 		ohlc, err := b.publicGetter.GetOHLC(b.ctx.Pair, b.ctx.Timeframe, b.ctx.OHLCLimit)
// 		helper.HandleFatalError(err)
// 		sma := smaFrom(b.ctx.SMALength, *ohlc, func(v bsre.OHLC) float64 { return v.Close })

// 		st.smaC <- sma[0]
// 	}

// 	smaGetter()

// 	for {
// 		<-ticker.C
// 		smaGetter()
// 	}
// }

// smaC       chan float64

// smaC:       make(chan float64),

// close(b.smaC)

// func ReorderOnSMA()  {

// }
