package util

import (
	"github.com/thoas/go-funk"
	"net/url"
	"regexp"
	"strings"
)

var DotIdExtensions = []string{
	".id", ".co.id", ".biz.id", ".or.id", ".web.id",
	".my.id", ".ac.id", ".sch.id", ".desa.id", ".ponpes.id",
}

var (
	domainRegExp = regexp.MustCompile(`[^a-zA-Z0-9-.]`) // Domain accepted characters
	wwwRegExp    = regexp.MustCompile(`(?i)www.`)       // For removing www
)

// emptySpace is an empty space for replacing
var emptySpace = []byte("")

// Domain returns a proper hostname / domain name. Preserve case is to flag keeping the case
// versus forcing to lowercase. Use the removeWww flag to strip the www sub-domain.
// This method returns an error if parse critically fails.
//
//	View examples: sanitize_test.go
func SanitizeDomain(original string, preserveCase bool, removeWww bool) (string, error) {

	// Try to see if we have a host
	if len(original) == 0 {
		return original, nil
	}

	// Missing http?
	if !strings.Contains(original, "http") {
		original = "http://" + strings.TrimSpace(original)
	}

	// Try to parse the url
	u, err := url.Parse(original)
	if err != nil {
		return original, err
	}

	// Remove leading www.
	if removeWww {
		u.Host = wwwRegExp.ReplaceAllString(u.Host, "")
	}

	// Keeps the exact case of the original input string
	if preserveCase {
		return string(domainRegExp.ReplaceAll([]byte(u.Host), emptySpace)), nil
	}

	// Generally all domains should be uniform and lowercase
	return string(domainRegExp.ReplaceAll([]byte(strings.ToLower(u.Host)), emptySpace)), nil
}

// Make sure to SanitizeDomain before call this
func ExtractTLD(name string) string {
	if strings.Contains(name, ".") {
		splitted := strings.Split(name, ".")
		tld1 := "." + splitted[len(splitted)-1]
		if len(splitted) == 2 {
			return tld1
		}

		tld2 := "." + splitted[len(splitted)-2] + tld1
		tlds := []string{
			".ac.id", ".ae.org", ".biz.id", ".br.com", ".cn.com", ".co.com", ".co.gg", ".co.id", ".co.in", ".co.je", ".co.uk", ".com.ai", ".com.de", ".com.se",
			".de.com", ".desa.id", ".eu.com", ".firm.in", ".gb.net", ".gen.in", ".gr.com", ".hu.net", ".in.net", ".ind.in", ".jp.net", ".jpn.com", ".me.uk",
			".mex.com", ".my.id", ".net.ai", ".net.gg", ".net.in", ".net.je", ".off.ai", ".or.id", ".org.ai", ".org.gg", ".org.in", ".org.je", ".org.uk",
			".ponpes.id", ".radio.am", ".radio.fm", ".ru.com", ".sa.com", ".sch.id", ".se.net", ".uk.com", ".uk.net", ".us.com", ".us.org", ".web.id", ".za.com",
		}
		if funk.ContainsString(tlds, tld2) {
			return tld2
		}
		return tld1
	}

	return ""
}

func ExtractDomainNameWithoutExt(domainName string) string {
	san := SanitizeDotExtensions(domainName)
	if strings.Contains(san, ".") {
		idx := strings.Index(san, ".")
		san = san[:idx]
	}

	//Remove TLD from sanitization
	key := KebabCase(san)
	key = strings.ToLower(key)
	return key
}

func IsDomainIsDotIdExtension(domainName string) bool {
	ext := ExtractTLD(domainName)
	return funk.ContainsString(DotIdExtensions, ext)
}

func IsDomainSubdomain(name string) (bool, error) {
	var err error
	name, err = SanitizeDomain(name, false, true)
	if err != nil {
		return false, err
	}

	tld := ExtractTLD(name)
	cleanWithoutTld := TrimSuffix(name, tld)
	return strings.Contains(cleanWithoutTld, "."), nil
}

func SanitizeDotExtensions(text string) string {
	text = strings.ToLower(text)
	regex := regexp.MustCompile(`[^a-z0-9.-]+`)
	return regex.ReplaceAllString(text, "")
}
