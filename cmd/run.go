/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/ivivanov/crypto-bot/app"
	"github.com/ivivanov/crypto-bot/bitstamp"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/spf13/cobra"
)

const (
	SMA  uint8 = 1
	GRID uint8 = 2
)

var (
	wsAddr   string
	wsScheme string
	pair     string
	profit   float64
	maker    float64
	taker    float64
	strategy uint8
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the bot",
	Long:  "You have to choose which trading strategy to use",
	Run: func(cmd *cobra.Command, args []string) {
		apiConn, err := bitstamp.NewAuthConn(apiKey, apiSecret, customerID, verbose)
		helper.HandleFatalError(err)

		config := &app.BotCtx{
			ApiConn:    apiConn,
			Debug:      verbose,
			Account:    account,
			Pair:       pair,
			Profit:     profit,
			MakerFee:   maker,
			TakerFee:   taker,
			WSScheme:   wsScheme,
			WSAddr:     wsAddr,
			APIKey:     apiKey,
			APISecret:  apiSecret,
			CustomerID: customerID,
			Timeframe:  timeframe,
			OHLCLimit:  ohlcLimit,
			SMALength:  smaLenght,
		}

		// if sma or grid {

		// }

		var trader app.Trader
		switch strategy {
		case SMA:
			trader = &app.SMATrader{
				Config: config,
				Order:  apiConn,
			}
		case GRID:
			trader = &app.GridTrader{
				Config:       config,
				OrderCreator: apiConn,
			}
		}

		bot, err := app.NewBot(config, trader)
		helper.HandleFatalError(err)

		log.Printf("pair: %v, profit: %v", pair, profit)

		helper.HandleFatalError(bot.Run())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&wsAddr, "ws-addr", "ws.bitstamp.net", "Bitstamp websocket address")
	runCmd.Flags().StringVar(&wsScheme, "ws-scheme", "wss", "Bitstamp websocket scheme")
	runCmd.Flags().Float64Var(&profit, "profit", 0.01, "Profit applied on each trade")
	runCmd.Flags().Float64Var(&maker, "maker", 0.02, "Maker fee %")
	runCmd.Flags().Float64Var(&taker, "taker", 0.03, "Taker fee %")
	runCmd.Flags().IntVar(&timeframe, "step", 3600, "Source Timeframe in seconds. Possible options are 60, 180, 300, 900, 1800, 3600, 7200, 14400, 21600, 43200, 86400, 259200")
	runCmd.Flags().IntVar(&ohlcLimit, "limit", 1000, "Source Limit OHLC results (minimum: 1; maximum: 1000)")
	runCmd.Flags().IntVar(&smaLenght, "length", 20, "SMA length")
	runCmd.Flags().Uint8Var(&strategy, "type", 1, "1 - SMA; 2 - GRID;")
}
