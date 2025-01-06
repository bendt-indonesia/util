package util

import (
	"math"
	"strings"
)

// Ceil
func CeilFloat64(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Ceil(val*ratio) / ratio
}

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

func Markup(markupType string, markupBy, baseNumber float64) float64 {
	if markupType == "FLAT_PRICE" {
		return baseNumber
	} else if markupType == "FIXED" {
		baseNumber += markupBy
	} else if markupType == "PERCENTAGE" {
		baseNumber = math.Round((100 + markupBy) / 100 * baseNumber)
	} else if markupType == "MARGIN_PERCENTAGE" {
		baseNumber = baseNumber / ((100 - markupBy) / 100)
	}

	return baseNumber
}

func CalculateMarkup(markupType string, markupBy, baseNumber float64) float64 {
	baseNumber = Markup(markupType, markupBy, baseNumber)
	baseNumber = RoundPrice(baseNumber)
	return baseNumber
}
