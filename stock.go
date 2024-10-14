package util

import (
	"math"
	"math/rand"
	"time"
)

func RoundStock(stock int) int {
	mmin := 20
	mmax := stock

	rand.Seed(time.Now().UnixNano())

	if stock == 0 {
		return stock
	}

	if stock > 5000 {
		mmin = 1450
		mmax = 2000
	} else if stock > 1000 {
		mmin = 500
		mmax = 1001
	} else if stock > 100 {
		mmin = 90
	} else if stock > 50 {
		mmin = 50
	} else if stock > 20 {
		mmin = 21
	} else if stock > 10 {
		mmin = 9
	} else if stock > 5 {
		mmin = stock
		mmax = stock + rand.Intn(3)
	} else {
		mmin = stock
		mmax = stock + 1
	}

	return rand.Intn(mmax-mmin+1) + mmin
}

func ResellerRoundStock(stock int, stock_if_zero int) int {
	var mmin int
	mmin = 20
	mmax := stock

	rand.Seed(time.Now().UnixNano())

	if stock == 0 {
		return stock_if_zero
	}

	if stock > 5000 {
		mmin = 300
		mmax = 500
	} else if stock > 1000 {
		mmin = 150
		mmax = 300
	} else if stock > 100 {
		mmin = 90
	} else if stock > 50 {
		mmin = 50
	} else if stock > 20 {
		mmin = 21
	} else if stock > 10 {
		mmin = 9
	} else if stock > 5 {
		mmin = stock
		mmax = stock + 3
	} else {
		mmin = stock
		mmax = stock + 1
	}

	return rand.Intn(mmax-mmin+1) + mmin
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
