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

// Check if phone begin with zero, then add country code
func PhonePrefix(phone string) string {
	phone = SanitizeStringNumber(phone)
	if len(phone) >= 1 {
		if phone[:1] == "0" {
			phone = "62" + phone[1:]
		}
	}
	return phone
}

func ValidatePhoneNo(n interface{}, allowedPrefix []string) error {
	_, err := ValidatePhoneNumber(n, allowedPrefix)
	return err
}

func ValidatePhoneNumber(n interface{}, allowedPrefix []string) (string, error) {
	var ps string
	if w, ok := n.(string); ok {
		ps = w
	}
	if w, ok := n.(*string); ok {
		ps = *w
	}

	ps = SanitizeStringNumber(ps)
	//Minimum Phone String is 5 digit , for PhonePrefix Prerequisities
	if len(ps) <= 5 {
		return ps, fmt.Errorf("Nomor handphone minimal adalah 5 digit")
	}

	ps = PhonePrefix(ps)
	for _, pf := range allowedPrefix {
		pfl := len(pf)
		if len(ps) < pfl {
			return ps, fmt.Errorf(fmt.Sprintf("Nomor handphone minimal adalah %d digit", pfl))
		}
		phonePrefix := ps[:pfl]
		if phonePrefix == pf {
			if pf == "62" {
				if len(ps) < 10 {
					return ps, fmt.Errorf("Nomor handphone minimal adalah 10 digit")
				}
				if funk.ContainsString(INDONESIA_PROVIDERS, ps[:5]) {
					return ps, nil
				}
				return ps, fmt.Errorf("Provider Nomor handphone tidak dikenali (Indonesia)")
			}

			//Phone prefix is passed or whitelisted in allowedPrefix
			return ps, nil
		}

	}

	if len(allowedPrefix) > 0 {
		return ps, fmt.Errorf("Nomor handphone `%s` tidak dapat digunakan.", ps)
	}

	return ps, nil
}
