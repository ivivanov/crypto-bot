package app

import (
	"encoding/json"
	"strconv"
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

	for i := 0; i < len(ohlc)-1; i++ {
		curr, _ := strconv.ParseInt(ohlc[i].Timestamp, 10, 64)
		next, _ := strconv.ParseInt(ohlc[i+1].Timestamp, 10, 64)

		if next-curr < 0 {
			t.Error("timestamp should be in chronological order")
		}

		if next-curr != 3600 {
			t.Errorf("exp: %v, act: %v", 3600, next-curr)
		}
	}

	sma := smaFrom(20, ohlc, func(v bsresponse.OHLC) float64 { return v.Close })

	if len(sma) != 1000 {
		t.Error("unexpected value")
	}
}
