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
	Short: "Queries data which does not require authentication",
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
