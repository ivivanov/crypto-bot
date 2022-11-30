package helper

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

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

func GetRnClientOrderID(account string) string {
	return fmt.Sprintf("%v-%v", account, time.Now().UnixMicro())
}

func GetAccountFrom(clientOrderID string) (string, error) {
	split := strings.Split(clientOrderID, "-")
	if len(split) != 2 {
		return "", fmt.Errorf("invalid client order id: %v", clientOrderID)
	}

	return split[0], nil
}
