package monetary

func SumPrices(prices []float64) float64 {
	var total float64

	for _, price := range prices {
		total += price
	}

	return total
}
