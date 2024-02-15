package util

import (
	"math"
	"strings"
)

// Round
func RoundFloat64(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func FloorFloat64(val float64, precision uint) float64 {
	rounded := RoundFloat64(val, precision)
	check := val - rounded
	//Negatif Value, means rounded to up
	if check < 0 {
		return rounded - (1 / math.Pow(10, float64(precision)))
	}

	return rounded
}

func PrecisionStr(s string, precision int) string {
	findIdx := strings.Index(s, ".")
	if findIdx == -1 {
		return s
	}
	lenLat := len(s)
	decimals := lenLat - 1 - findIdx
	if decimals > precision {
		until := lenLat - (decimals - precision)
		s = s[:until]
	}
	return s
}
