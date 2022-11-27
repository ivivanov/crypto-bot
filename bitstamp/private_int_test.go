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
