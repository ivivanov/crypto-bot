package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	customerIDFlag *string
	apiKeyFlag     *string
	apiSecretFlag  *string
	wsAddrFlag     *string
	wsSchemeFlag   *string
	pairFlag       *string
	feeFlag        *float64
)

func ParseFlags() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env missing")
	}

	customerIDFlag = flag.String("id", os.Getenv("BS_CUST_ID"), "customer id")
	apiKeyFlag = flag.String("key", os.Getenv("BS_API_KEY"), "api key")
	apiSecretFlag = flag.String("secret", os.Getenv("BS_API_SECRET"), "api secret")
	wsAddrFlag = flag.String("addr", "ws.bitstamp.net", "websocket address")
	wsSchemeFlag = flag.String("scheme", "wss", "websocket scheme")
	pairFlag = flag.String("pair", "usdtusd", "trading pair")
	feeFlag = flag.Float64("fee", 0.02, "maker fee")

	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			log.Printf("Flag: %v is required", f.Name)
			flag.PrintDefaults()
			os.Exit(1)
		}
	})
}
