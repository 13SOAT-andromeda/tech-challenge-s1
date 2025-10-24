package monetary

import "testing"

func TestConvertToMinorUnitInt64_BrazilianExamples(t *testing.T) {
	cases := []struct {
		name     string
		amount   float64
		decimals uint
		want     int64
	}{
		{"R$ 1.234,56", 1234.56, 2, 123456},
		{"R$ 0,99", 0.99, 2, 99},
		{"R$ 0,005 (3 decimals)", 0.005, 3, 5},
		// @TODO: after add field to uint64 {"-R$ 12,34 (negative)", -12.34, 2, -1234},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ConvertToMinorUnitInt64(tc.amount, tc.decimals)
			if got != tc.want {
				t.Fatalf("%s: ConvertToMinorUnitInt64(%v, %d) = %d; want %d", tc.name, tc.amount, tc.decimals, got, tc.want)
			}
		})
	}
}
