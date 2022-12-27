/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/ivivanov/crypto-bot/app"
	"github.com/ivivanov/crypto-bot/app/traders"
	"github.com/ivivanov/crypto-bot/bitstamp"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/spf13/cobra"
)

const (
	SMA  uint8 = 1
	GRID uint8 = 2
)

var (
	wsAddr         string
	wsScheme       string
	pair           string
	gridProfit     float64
	maker          float64
	taker          float64
	strategy       uint8
	smaOffset      float64
	smaProfit      float64
	smaOrderAmount float64
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
			WSScheme:   wsScheme,
			WSAddr:     wsAddr,
			APIKey:     apiKey,
			APISecret:  apiSecret,
			CustomerID: customerID,
		}

		var trader app.Trader
		switch strategy {
		case SMA:
			trader = &traders.SMATrader{
				Ctx:       config,
				ApiSMA:    apiConn,
				Timeframe: timeframe,
				OHLCLimit: ohlcLimit,
				Length:    smaLength,
				Offset:    smaOffset,
				Profit:    smaProfit,
				Amount:    smaOrderAmount,
			}

			log.Printf("Starting SMA strategy, pair: %v, offset: %v, profit: %v, timeframe: %v", pair, smaOffset, smaProfit, smaLength)
		case GRID:
			trader = &traders.GridTrader{
				Ctx:          config,
				OrderCreator: apiConn,
				Profit:       gridProfit,
				MakerFee:     maker,
				TakerFee:     taker,
			}

			log.Printf("Starting grid strategy, pair: %v, profit: %v", pair, gridProfit)
		}

		bot, err := app.NewBot(config, trader)
		helper.HandleFatalError(err)
		helper.HandleFatalError(bot.Run())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&wsAddr, "ws-addr", "ws.bitstamp.net", "Bitstamp websocket address")
	runCmd.Flags().StringVar(&wsScheme, "ws-scheme", "wss", "Bitstamp websocket scheme")
	runCmd.Flags().Float64Var(&gridProfit, "grid-profi", 0.01, "Grid trader profit in %")
	runCmd.Flags().Float64Var(&maker, "maker", 0.04, "Maker fee %")
	runCmd.Flags().Float64Var(&taker, "taker", 0.06, "Taker fee %")
	runCmd.Flags().IntVar(&timeframe, "timeframe", 60, "Source Timeframe in seconds. Possible options are 60, 180, 300, 900, 1800, 3600, 7200, 14400, 21600, 43200, 86400, 259200")
	runCmd.Flags().IntVar(&ohlcLimit, "limit", 10, "Source Limit OHLC results (minimum: 1; maximum: 1000)")
	runCmd.Flags().Uint8Var(&strategy, "type", 1, "1 - SMA; 2 - GRID;")
	runCmd.Flags().IntVar(&smaLength, "sma-length", 8, "SMA length")
	runCmd.Flags().Float64Var(&smaOffset, "sma-offset", 0.0003, "SMA trader offset")
	runCmd.Flags().Float64Var(&smaProfit, "sma-profit", 0.0009, "SMA trader profit in cents")
	runCmd.Flags().Float64Var(&smaOrderAmount, "amount", 0, "SMA trade order amount")
}
