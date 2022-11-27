package bitstamp

import "testing"

func TestReadSecret(t *testing.T) {
	secret, err := GetSecret()
	if err != nil {
		t.Error(err)
	}

	if secret.CustomerID == "" || secret.Key == "" || secret.Secret == "" {
		t.Fatal("Missing secret params.")
	}
}
