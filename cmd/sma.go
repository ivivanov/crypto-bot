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

// smaCmd represents the sma command
var smaCmd = &cobra.Command{
	Use:   "sma",
	Short: "Calculates SMA",
	Long:  "Simple moving average of `source` for `period` bars back.",
	Run: func(cmd *cobra.Command, args []string) {
		querier, err := app.NewQuerier(apiKey, apiSecret, customerID, verbose)
		helper.HandleFatalError(err)
		result, err := querier.SMA(pair, timeframe, ohlcLimit, smaLength)
		helper.HandleFatalError(err)

		count := len(result)
		log.Printf("now-5: %v, now-4: %v, now-3: %v, now-2: %v, now-1: %v, now: %v",
			helper.Round5dec(result[count-6]),
			helper.Round5dec(result[count-5]),
			helper.Round5dec(result[count-4]),
			helper.Round5dec(result[count-3]),
			helper.Round5dec(result[count-2]),
			helper.Round5dec(result[count-1]),
		)
	},
}

func init() {
	queryCmd.AddCommand(smaCmd)

	smaCmd.Flags().IntVar(&timeframe, "timeframe", 3600, "Source Timeframe in seconds. Possible options are 60, 180, 300, 900, 1800, 3600, 7200, 14400, 21600, 43200, 86400, 259200")
	smaCmd.Flags().IntVar(&ohlcLimit, "limit", 1000, "Source Limit OHLC results (minimum: 1; maximum: 1000)")
	smaCmd.Flags().IntVar(&smaLength, "length", 20, "SMA length")
}
