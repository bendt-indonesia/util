package util

import (
	"strings"
)

func UnhexUuidToBinary(f string, prefix string) string {
	key := strings.ReplaceAll(f, "-","")
	q := "UNHEX(REPLACE("
	q += prefix + key
	q += ", '-',''))"

	return q
}
