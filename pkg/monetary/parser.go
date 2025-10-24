package monetary

import "math"

func ConvertToMinorUnitInt64(amount float64, decimals uint) int64 {
	scaleFactor := math.Pow(10, float64(decimals))

	// Scale and round
	scaled := amount * scaleFactor
	rounded := math.Round(scaled)

	// Convert the rounded float64 to int64.
	return int64(rounded)
}
