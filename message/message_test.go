package message

import (
	"encoding/json"
	"testing"
)

func TestMyOrdersMessage(t *testing.T) {
	myOrdersMessage := MyOrdersMessage("usdtusd", "wstoken", 1234)
	b, err := json.Marshal(myOrdersMessage)
	if err != nil {
		t.Error(err)
	}

	exp := `{"event":"bts:subscribe","data":{"channel":"private-my_orders_usdtusd-1234","auth":"wstoken"}}`
	if string(b) != exp {
		t.Fatal("Unexpected my orders message")
	}
}
