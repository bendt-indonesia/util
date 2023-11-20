package util

import (
	"fmt"
	"strconv"
	"time"
)

func Float64ToInt64(f float64) int64 {
	return int64(f)
}

func Float64ToString(s float64, p *int) string {
	f := "%"
	if p != nil {
		if *p == 0 {
			f += "d"
		} else {
			f += "."+IntToStr(*p)
		}
	} else {
		f += "f"
	}

	str := fmt.Sprintf(f, s)
	return str
}

func IntToLetters(number int32) (letters string){
	number--
	if firstLetter := number/26; firstLetter >0{
		letters += IntToLetters(firstLetter)
		letters += string('A' + number%26)
	} else {
		letters += string('A' + number)
	}

	return
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func IntToUint(n int) uint {
	return uint(n)
}

func IntToUint64(n int) uint64 {
	return uint64(n)
}

//Boolean func returns boolean value of string value like on, off, 0, 1, yes, no
func StringToInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)

	return n
}

//Boolean func returns boolean value of string value like on, off, 0, 1, yes, no
func StringToUint(s string) uint {
	var branchUint uint
	branchU64, _ := strconv.ParseUint(s, 10, 32)
	branchUint = uint(branchU64)

	return branchUint
}

func StringToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s,64)
	return f
}

func UintToString(s uint) string {
	return fmt.Sprintf("%d", s)
}

func StringToTime(timeStr string, adjustTimeZone int) *time.Time {

	t2, err := time.Parse("2006-01-02", timeStr)
	if err != nil {
		//panic(err)
		return nil
	}

	nt := t2.Add(time.Hour * time.Duration(adjustTimeZone))

	return &nt
}

func WeightTypeConvert(wt string, before float64) float64 {
	if wt == "G" {
		return before
	}
	if wt == "KG" {
		return before * 1000
	}
	if wt == "OZ" {
		return before * 28.3495
	}
	if wt == "LB" {
		return before * 453.592
	}

	return before
}
