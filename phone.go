package util

import (
	"fmt"
	"github.com/thoas/go-funk"
)

var INDONESIA_PROVIDERS = []string{
	//Operator XL Axiata
	"62859",
	"62877",
	"62878",
	"62817",
	"62818",
	"62819",

	//Operator Telkomsel
	"62811",
	"62812",
	"62813",
	"62821",
	"62822",
	"62823",
	"62852",
	"62853",
	"62851",

	//Operator 3
	"62898",
	"62899",
	"62895",
	"62896",
	"62897",

	//Oreedo
	"62814",
	"62815",
	"62816",
	"62855",
	"62856",
	"62857",
	"62858",

	//SmartFren
	"62889",
	"62881",
	"62882",
	"62883",
	"62886",
	"62887",
	"62888",
	"62884",
	"62885",

	//Providers
	"62832",
	"62833",
	"62838",
	"62831",
}

// allowedPrefix = country code ( ex 62 )
func ValidatePhoneNo(n string, allowedPrefix []string) error {
	ph := SanitizeStringNumber(n)
	ph = PhonePrefix(ph)

	if len(ph) < 5 {
		return fmt.Errorf("Nomor handphone minimal adalah 5 digit")
	}

	for _, pf := range allowedPrefix {
		pfl := len(pf)
		phonePrefix := ph[:pfl]
		if phonePrefix == pf {
			if pf == "62" {
				if len(ph) < 10 {
					return fmt.Errorf("Nomor handphone minimal adalah 10 digit")
				}
				if funk.ContainsString(INDONESIA_PROVIDERS, ph[:5]) {
					return nil
				}
			}

			//Phone prefix is passed or whitelisted in allowedPrefix
			return nil
		}

	}

	if len(allowedPrefix) > 0 {
		return fmt.Errorf("Nomor handphone `%s` tidak dapat digunakan.")
	}

	return nil
}
