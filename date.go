package util

import (
	"time"
)

func CalculateDays(d1 time.Time, d2 time.Time) int {
	diff := d2.Sub(d1)

	return int(diff.Hours() / 24)
}

func CalculateMonths(d1, d2 time.Time) int {
	diff := d2.Sub(d1)

	return int(diff.Hours() / 24 / 31)
}

func CalculateAge(birthdate time.Time) int {
	today := time.Now().In(birthdate.Location())
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	by, bm, bd := birthdate.Date()
	birthdate = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	if today.Before(birthdate) {
		return 0
	}
	age := ty - by
	anniversary := birthdate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}
	return age
}

func GetMonthTxt(idx int) string {
	months := []string{
		"Jan",
		"Feb",
		"Mar",
		"Apr",
		"May",
		"Jun",
		"Jul",
		"Aug",
		"Sep",
		"Oct",
		"Nov",
		"Des",
	}

	if idx < 0 && idx > 11 {
		return ""
	}

	return months[idx]
}
