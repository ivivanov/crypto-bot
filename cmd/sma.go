/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// smaCmd represents the sma command
var smaCmd = &cobra.Command{
	Use:   "sma",
	Short: "Calculates SMA",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	queryCmd.AddCommand(smaCmd)

	smaCmd.Flags().IntVar(&step, "step", 3600, "Timeframe in seconds. Possible options are 60, 180, 300, 900, 1800, 3600, 7200, 14400, 21600, 43200, 86400, 259200")
	smaCmd.Flags().IntVar(&limit, "limit", 1000, "Limit OHLC results (minimum: 1; maximum: 1000)")
}
