package cmd

import (
	"os"

	"github.com/ivivanov/crypto-bot/helper"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	account    string
	customerID string
	apiKey     string
	apiSecret  string
	verbose    bool
	timeframe  int
	ohlcLimit  int
	smaLenght  int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crypto-bot",
	Short: "Automated grid trading strategy",
	Long: `Run prepare first. This will open tagged buy/sell orders into the sub-account.
Tags are used to route the orders to their specific sub-account.
Manually created orders will not be processed.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	helper.HandleFatalError(rootCmd.Execute())
}

func init() {
	helper.HandleFatalError(godotenv.Load())

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&account, "acc", os.Getenv("BS_ACCOUNT"), "Bitstamp account/sub-account")
	rootCmd.PersistentFlags().StringVar(&customerID, "id", os.Getenv("BS_CUSTOMER_ID"), "Bitstamp customer ID")
	rootCmd.PersistentFlags().StringVar(&apiKey, "key", os.Getenv("BS_API_KEY"), "Bitstamp API key")
	rootCmd.PersistentFlags().StringVar(&apiSecret, "secret", os.Getenv("BS_API_SECRET"), "Bitstamp API secret")
	rootCmd.PersistentFlags().StringVar(&pair, "pair", os.Getenv("BS_PAIR"), "Trading pair")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Dumps requests and other debug info")
}
