package helper

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
)

func TestGetClientOrderID(t *testing.T) {
	expMaxLen := 36                                                                                                             // len: 36
	price := "0.99950"                                                                                                          // len: 7 (we round to the 5th dec digit)
	maxAccLen := expMaxLen - len(fmt.Sprintf("%s%v%s%v", ClientOrderIdDelimiter, price, ClientOrderIdDelimiter, math.MaxInt32)) // len: 17

	for _, tc := range []struct {
		account  string
		price    string
		rnSuffix int32
	}{
		{
			account:  strings.Repeat("c", maxAccLen),
			price:    "0.12345",
			rnSuffix: rand.Int31n(math.MaxInt32),
		},
		{
			account:  "test123",
			price:    "1.00",
			rnSuffix: rand.Int31n(math.MaxInt32),
		},
		{
			account:  "t",
			price:    "1.00",
			rnSuffix: rand.Int31n(math.MaxInt32),
		},
		{
			account:  "testsafsaf",
			price:    "0.00000",
			rnSuffix: rand.Int31n(math.MaxInt32),
		},
		{
			account:  "t",
			price:    "0.99999",
			rnSuffix: rand.Int31n(math.MaxInt32),
		},
	} {
		if len(tc.account) > maxAccLen {
			t.Errorf("exp: %v, act: %v", len(tc.account), maxAccLen)
		}

		clientOrderID := strings.ToLower(fmt.Sprintf("%v%s%v%s%v", tc.account, ClientOrderIdDelimiter, tc.price, ClientOrderIdDelimiter, tc.rnSuffix))
		if len(clientOrderID) > expMaxLen {
			t.Errorf("exp: %v, act: %v, %v", expMaxLen, len(clientOrderID), clientOrderID)
		}
	}
}

func TestGetPriceFrom(t *testing.T) {
	expPrice := 0.99981
	clientOrderID := GetClientOrderID("account", expPrice)
	actPrice, err := GetPriceFrom(clientOrderID)
	if err != nil {
		t.Fatalf("error parsing price: %v", err)
	}

	if expPrice != actPrice {
		t.Errorf("ext: %v, act: %v", expPrice, actPrice)
	}
}
