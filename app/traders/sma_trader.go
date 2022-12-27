package traders

import (
	"time"

	bsre "github.com/ivivanov/crypto-bot/bitstamp/response"

	"github.com/ivivanov/crypto-bot/app"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/ivivanov/crypto-bot/response"
)

const (
	Buy  int = 0
	Sell int = 1
)

type ApiSMA interface {
	PostSellLimitOrder(currencyPair, clientOrderID string, amount, price float64) (*bsre.LimitOrder, error)
	PostBuyLimitOrder(currencyPair, clientOrderID string, amount, price float64) (*bsre.LimitOrder, error)
	PostCancelOrder(id int64) (*bsre.CancelOrder, error)
	PostOpenOrders(currencyPair string) (*[]bsre.OpenOrder, error)
	GetOHLC(currencyPair string, step, limit int) (*[]bsre.OHLC, error)
}

type SMATrader struct {
	Ctx       *app.BotCtx
	ApiSMA    ApiSMA
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

	go st.SMAPolling(smaC)
	sma := 0.0

	for {
		select {
		case trade := <-tradeC:
			//	buy trade
			//		create sell apply fixed profit
			//	sell trade (1 cycle done)
			//		create buy order (with current SMA)

			switch trade.Data.Side {
			case "buy":
				sellPrice := helper.Round5dec(trade.Price() + st.Profit)
				_, err := st.ApiSMA.PostSellLimitOrder(st.Ctx.Pair, helper.GetClientOrderID(st.Ctx.Account, sellPrice), trade.Amount(), sellPrice)
				helper.HandleFatalError(err)
			case "sell":
				buyPrice := helper.Round5dec(sma - st.Offset)
				_, err := st.ApiSMA.PostBuyLimitOrder(st.Ctx.Pair, helper.GetClientOrderID(st.Ctx.Account, buyPrice), trade.Amount(), buyPrice)
				helper.HandleFatalError(err)
			}
		case newSMA := <-smaC:
			//	get open orders
			//	no open orders
			//		create buy order (with new SMA)
			// 	loop orders and check if limit price != new SMA
			//		re-create buy order (with new SMA)

			sma = newSMA

			openOrders, err := st.ApiSMA.PostOpenOrders(st.Ctx.Pair)
			helper.HandleFatalError(err)

			buyPrice := helper.Round5dec(sma - st.Offset)

			if len(*openOrders) == 0 { // Initial order
				_, err := st.ApiSMA.PostBuyLimitOrder(st.Ctx.Pair, helper.GetClientOrderID(st.Ctx.Account, buyPrice), st.Amount, buyPrice)
				helper.HandleFatalError(err)
			} else {
				// TODO: if there are multiple orders on the same price aggregate them into 1 order
				for _, order := range *openOrders {
					if order.Type == Buy && order.Price != buyPrice {
						_, err := st.ApiSMA.PostCancelOrder(order.ID)
						helper.HandleFatalError(err)
						_, err = st.ApiSMA.PostBuyLimitOrder(st.Ctx.Pair, helper.GetClientOrderID(st.Ctx.Account, buyPrice), st.Amount, buyPrice)
						helper.HandleFatalError(err)
					}
				}
			}
		}
	}
}

func (st *SMATrader) SMAPolling(smaC chan<- float64) {
	updateInterval := time.Duration(st.Timeframe) * time.Second
	ticker := time.NewTicker(updateInterval)

	smaGetter := func() {
		ohlc, err := st.ApiSMA.GetOHLC(st.Ctx.Pair, st.Timeframe, st.OHLCLimit)
		helper.HandleFatalError(err)
		sma := helper.SMAFrom(st.Length, *ohlc, func(v bsre.OHLC) float64 { return v.Close })

		smaC <- sma[len(sma)-1]
	}

	smaGetter()

	for range ticker.C {
		smaGetter()
	}
}
