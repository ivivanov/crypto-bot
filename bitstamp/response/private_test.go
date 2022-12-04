package response

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalCancelAllOrders(t *testing.T) {
	resJson := `
	{
		"success": true,
		"canceled": [
			{
				"currency_pair": "USDT/USD",
				"price": 1.0001,
				"amount": 33.3,
				"type": 1,
				"id": 1561930221211649
			}
		]
	}
	`

	res := &CancelAllOrders{}
	err := json.Unmarshal([]byte(resJson), res)
	if err != nil {
		t.Fatal(err)
	}

	if !res.Success {
		t.Error("exp success to be true")
	}

	if len(res.Canceled) != 1 {
		t.Fatalf("exp: %v, act: %v", 1, len(res.Canceled))
	}

	if res.Canceled[0].Amount != 33.3 {
		t.Errorf("exp: %v, act: %v", 33.3, res.Canceled[0].Amount)
	}
}
