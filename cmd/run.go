/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/ivivanov/crypto-bot/app"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/spf13/cobra"
)

var (
	wsAddr   string
	wsScheme string
	pair     string
	profit   float64
	maker    float64
	taker    float64
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the bot",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		bot, err := app.NewBot(account, wsScheme, wsAddr, apiKey, apiSecret, customerID, pair, profit, maker, taker, verbose)
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
}
