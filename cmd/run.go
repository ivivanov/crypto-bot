/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	wsAddr   string
	wsScheme string
	pair     string
	profit   float64
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the bot",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// app, err := NewApp()
		// helper.HandleFatalError(err)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	rootCmd.Flags().StringVar(&wsAddr, "ws-addr", "ws.bitstamp.net", "Bitstamp websocket address")
	rootCmd.Flags().StringVar(&wsScheme, "ws-scheme", "wss", "Bitstamp websocket scheme")
	rootCmd.Flags().Float64Var(&profit, "profit", 0.0, "Profit applied on each trade")
}
