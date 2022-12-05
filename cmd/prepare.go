/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ivivanov/crypto-bot/app"
	"github.com/ivivanov/crypto-bot/helper"
	"github.com/spf13/cobra"
)

type OrderType int8

const (
	Buy  uint8 = 0
	Sell uint8 = 1
	Both uint8 = 2
)

var (
	price      float64
	bank       float64
	grid       float64
	orders     float64
	ordersType uint8
)

var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Opens the initial buy orders",
	Long: `Grid bot automatically opens new order when the order is filled.
For this reason we need to create the initial orders based on few configurations.`,
	Run: func(cmd *cobra.Command, args []string) {
		preparer, err := app.NewPreparer(account, apiKey, apiSecret, customerID, pair, verbose)
		helper.HandleFatalError(err)

		switch ordersType {
		case Buy:
			helper.HandleFatalError(preparer.OpenBuyOrders(bank, price, grid, orders))
		case Sell:
			helper.HandleFatalError(preparer.OpenSellOrders(bank, price, grid, orders))
		case Both:
			helper.HandleFatalError(preparer.OpenBuySellOrders(bank, price, grid, orders))
		}
	},
}

func init() {
	rootCmd.AddCommand(prepareCmd)

	prepareCmd.Flags().Float64Var(&price, "price", 0, "Entry price")
	prepareCmd.Flags().Float64Var(&bank, "bank", 0, "Start capital")
	prepareCmd.Flags().Float64Var(&grid, "grid", 0.04, "Grid step %")
	prepareCmd.Flags().Float64Var(&orders, "orders", 6, "Initial orders count")
	prepareCmd.Flags().Uint8Var(&ordersType, "type", 0, "0 - only buy; 1 - only sell; 2 - both")

	prepareCmd.MarkFlagRequired("price")
	prepareCmd.MarkFlagRequired("bank")
}
