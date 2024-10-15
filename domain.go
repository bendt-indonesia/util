package util

import (
	"github.com/thoas/go-funk"
	"regexp"
	"strings"
)

var DotIdExtensions = []string{
	".id", ".co.id", ".biz.id", ".or.id", ".web.id",
	".my.id", ".ac.id", ".sch.id", ".desa.id", ".ponpes.id",
}

func ExtractTLD(domainName string) string {
	san := SanitizeDotExtensions(domainName)
	if strings.Contains(san, ".") {
		idx := strings.Index(san, ".")
		return san[idx:]
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

func SanitizeDotExtensions(text string) string {
	text = strings.ToLower(text)
	regex := regexp.MustCompile(`[^a-z0-9.-]+`)
	return regex.ReplaceAllString(text, "")
}
