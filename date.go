package util

import (
	"time"
)

func CalculateDays(d1 time.Time, d2 time.Time) int {
	diff := d2.Sub(d1)

	return int(diff.Hours()/24)
}
