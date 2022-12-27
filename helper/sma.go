package helper

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

func SMAFrom[T any](length int, history []T, getVal func(v T) float64) []float64 {
	result := make([]float64, len(history))
	sum := float64(0)

	for i, ohlc := range history {
		count := i + 1
		sum += getVal(ohlc)

		if i >= length {
			sum -= getVal(history[i-length])
			count = length
		}

		result[i] = sum / float64(count)
	}

	return result
}
