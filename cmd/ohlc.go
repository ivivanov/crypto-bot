/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ivivanov/crypto-bot/app"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/spf13/cobra"
)

var (
	step  int
	limit int
)

// ohlcCmd represents the ohlc command
var ohlcCmd = &cobra.Command{
	Use:   "ohlc",
	Short: "OHLC data for a timeframe",
	Run: func(cmd *cobra.Command, args []string) {
		querier, err := app.NewQuerier(apiKey, apiSecret, customerID, verbose)
		helper.HandleFatalError(err)
		helper.HandleFatalError(querier.OHLC(pair, step, limit))
	},
}

func init() {
	queryCmd.AddCommand(ohlcCmd)

	ohlcCmd.Flags().IntVar(&step, "step", 3600, "Timeframe in seconds. Possible options are 60, 180, 300, 900, 1800, 3600, 7200, 14400, 21600, 43200, 86400, 259200")
	ohlcCmd.Flags().IntVar(&limit, "limit", 1000, "Limit OHLC results (minimum: 1; maximum: 1000)")
}
