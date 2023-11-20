package util

import (
	"math"
	"math/rand"
	"time"
)

func RoundStock(stock int) int {
	min := 20
	max := stock

	rand.Seed(time.Now().UnixNano())

	if stock == 0 {
		return stock
	}

	if stock > 5000 {
		min = 1450
		max = 2000
	} else if stock > 1000 {
		min = 500
		max = 1001
	} else if stock > 100 {
		min = 90
	} else if stock > 50 {
		min = 50
	} else if stock > 20 {
		min = 21
	} else if stock > 10 {
		min = 9
	} else if stock > 5 {
		min = stock
		max = stock + rand.Intn(3)
	} else {
		min = stock
		max = stock + 1
	}

	return rand.Intn(max - min + 1) + min
}

func ResellerRoundStock(stock int, stock_if_zero int) int {
	var min int
	min = 20
	max := stock

	rand.Seed(time.Now().UnixNano())

	if stock == 0 {
		return stock_if_zero
	}

	if stock > 5000 {
		min = 300
		max = 500
	} else if stock > 1000 {
		min = 150
		max = 300
	} else if stock > 100 {
		min = 90
	} else if stock > 50 {
		min = 50
	} else if stock > 20 {
		min = 21
	} else if stock > 10 {
		min = 9
	} else if stock > 5 {
		min = stock
		max = stock + 3
	} else {
		min = stock
		max = stock + 1
	}

	return rand.Intn(max - min + 1) + min
}

func ResellerRoundStock64(stock float64) float64 {
	rounded := ResellerRoundStock(int(stock), 0)
	return float64(rounded)
}

func RoundPrice(sellPrice float64) float64 {
	if sellPrice < 1000000 {
		return math.Round(float64(sellPrice)/1000) * 1000
	} else if sellPrice < 5000000 {
		return math.Round(float64(sellPrice)/10000) * 10000
	}

	return math.Round(sellPrice/100000) * 100000
}


