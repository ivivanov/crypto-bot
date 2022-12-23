/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ivivanov/crypto-bot/app"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/spf13/cobra"
)

// ohlcCmd represents the ohlc command
var ohlcCmd = &cobra.Command{
	Use:   "ohlc",
	Short: "OHLC data for a timeframe",
	Run: func(cmd *cobra.Command, args []string) {
		querier, err := app.NewQuerier(apiKey, apiSecret, customerID, verbose)
		helper.HandleFatalError(err)
		res, err := querier.OHLC(pair, timeframe, ohlcLimit)
		helper.HandleFatalError(err)
		helper.PrintIdent(res)
	},
}

func init() {
	queryCmd.AddCommand(ohlcCmd)

	ohlcCmd.Flags().IntVar(&timeframe, "timeframe", 3600, "Timeframe in seconds. Possible options are 60, 180, 300, 900, 1800, 3600, 7200, 14400, 21600, 43200, 86400, 259200")
	ohlcCmd.Flags().IntVar(&ohlcLimit, "limit", 1000, "Limit OHLC results (minimum: 1; maximum: 1000)")
}
