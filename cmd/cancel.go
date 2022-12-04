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
	cancelAll bool
)

// cancelCmd represents the cancel command
var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Close order/s",
	Run: func(cmd *cobra.Command, args []string) {
		canceller, err := app.NewCanceler(apiKey, apiSecret, customerID, pair, verbose)
		helper.HandleFatalError(err)
		if cancelAll {
			helper.HandleFatalError(canceller.CancelAll())
		}
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)

	cancelCmd.Flags().BoolVar(&cancelAll, "all", true, "Cancel all open orders")
}
