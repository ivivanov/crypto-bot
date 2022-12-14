/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ivivanov/crypto-bot/app"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/spf13/cobra"
)

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Get account balance",
	Run: func(cmd *cobra.Command, args []string) {
		balancer, err := app.NewBalancer(apiKey, apiSecret, customerID, verbose)
		helper.HandleFatalError(err)
		helper.HandleFatalError(balancer.BalanceAll(""))
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
