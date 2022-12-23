/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// balanceCmd represents the balance command
var queryCmd = &cobra.Command{
	Use:   "q",
	Short: "Queries data: balance, OHLC, SMA",
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
