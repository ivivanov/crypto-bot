package response

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalOHLC(t *testing.T) {
	resJson := `
	{
		"data": {
		  "ohlc": [
			{
			  "close": "1.00000",
			  "high": "1.00001",
			  "low": "1.00000",
			  "open": "1.00001",
			  "timestamp": "1667678400",
			  "volume": "5550.24470"
			},
			{
			  "close": "1.00000",
			  "high": "1.00002",
			  "low": "1.00000",
			  "open": "1.00001",
			  "timestamp": "1667682000",
			  "volume": "25753.43649"
			},
			{
			  "close": "1.00000",
			  "high": "1.00001",
			  "low": "1.00000",
			  "open": "1.00000",
			  "timestamp": "1667685600",
			  "volume": "4275.08018"
			}
		  ],
		  "pair": "USDT/USD"
		}
	  }
	`

	res := map[string]OHLCData{}
	err := json.Unmarshal([]byte(resJson), &res)
	if err != nil {
		t.Fatal(err)
	}

	if res["data"].OHCL[0].Close != 1.00000 {
		t.Error("unexpected value")
	}

	if res["data"].OHCL[0].Volume != 5550.24470 {
		t.Error("unexpected value")
	}
}
