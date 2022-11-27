package bitstamp

import (
	_ "embed"
	"encoding/json"
)

//go:embed config/secret.json
var secretF []byte

type Secret struct {
	CustomerID string `json:"customerId"`
	Key        string `json:"key"`
	Secret     string `json:"secret"`
}

func GetSecret() (*Secret, error) {
	secret := &Secret{}

	err := json.Unmarshal(secretF, secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
