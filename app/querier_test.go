package app

import (
	"encoding/json"
	"testing"

	_ "embed"

	bsresponse "github.com/ivivanov/crypto-bot/bitstamp/response"
)

//go:embed ohlc_test.json
var ohlcData string

func TestSmaFrom(t *testing.T) {
	ohlc := []bsresponse.OHLC{}
	err := json.Unmarshal([]byte(ohlcData), &ohlc)
	if err != nil {
		t.Fatal(err)
	}

	sma := smaFrom(20, ohlc, func(v bsresponse.OHLC) float64 { return v.Close })

	if len(sma) != 1000 {
		t.Error("unexpected value")
	}
}
