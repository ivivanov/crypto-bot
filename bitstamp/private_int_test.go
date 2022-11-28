package bitstamp

import (
	_ "embed"

	"testing"
)

func TestWsToken(t *testing.T) {
	secret, err := GetSecret()
	if err != nil {
		t.FailNow()
	}

	api, err := NewAuthConn(secret.Key, secret.Secret, secret.CustomerID)
	if err != nil {
		t.FailNow()
	}
	defer api.Close()

	token, err := api.WebsocketToken()
	if err != nil {
		t.FailNow()
	}

	if token.Token == "" || token.UserID == 0 || token.ValidSec == 0 {
		t.Errorf("Invalid token response")
	}
}

// func TestPostSellLimitOrder(t *testing.T) {
// 	secret, err := GetSecret()
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	api, err := NewAuthConn(secret.Key, secret.Secret, secret.CustomerID)
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	defer api.Close()

// 	// &response.SellLimitOrder{Price:1.00002, Amount:15.003, Type:1, ID:0, DateTime:response.DateTime{wall:0xc0d8f9770996fcb7, ext:6919593, loc:(*time.Location)(0x7bfd40)}}
// 	resp, err := api.PostSellLimitOrder("usdtusd", 12.00000, 1.00250)
// 	t.Log(resp)
// 	t.Log(err)

// 	t.Fail()
// }
