package helper

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const ClientOrderIdDelimiter = "_"

func HandleFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Round5dec(num float64) float64 {
	precision := 100_000
	return math.Round(num*float64(precision)) / float64(precision)
}

func Round2dec(num float64) float64 {
	precision := 100
	return math.Round(num*float64(precision)) / float64(precision)
}

// [account]_[target-price]_[random] - ClientOrderID format
func GetClientOrderID(account string, price float64) string {
	seed := rand.NewSource(time.Now().UnixNano())
	rn := rand.New(seed)
	rnSuffix := rn.Int31n(math.MaxInt32)

	return strings.ToLower(fmt.Sprintf("%v%s%v%s%v", account, ClientOrderIdDelimiter, price, ClientOrderIdDelimiter, rnSuffix))
}

// [account]_[target-price]_[random] - ClientOrderID format
func GetPriceFrom(clientOrderID string) (float64, error) {
	split := strings.Split(clientOrderID, ClientOrderIdDelimiter)
	if len(split) != 3 {
		return 0, fmt.Errorf("invalid client order id: %v", clientOrderID)
	}

	if split[0] == "" {
		return 0, fmt.Errorf("invalid account")
	}

	price, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price: %v", split[1])
	}

	return price, nil
}

func GetAccountFrom(clientOrderID string) (string, error) {
	split := strings.Split(clientOrderID, ClientOrderIdDelimiter)
	if len(split) != 3 {
		return "", fmt.Errorf("invalid client order id: %v", clientOrderID)
	}

	return split[0], nil
}

// Simple Moving Average (SMA).
func Sma(period int, prices []float64) []float64 {
	result := make([]float64, len(prices))
	sum := float64(0)

	for i, p := range prices {
		count := i + 1
		sum += p

		if i >= period {
			sum -= prices[i-period]
			count = period
		}

		result[i] = sum / float64(count)
	}

	return result
}
