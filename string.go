package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gertd/go-pluralize"
	"github.com/gobeam/stringy"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Boolean func returns boolean value of string value like on, off, 0, 1, yes, no
func Boolean(str string) bool {
	b := stringy.New(str)
	return b.Boolean()
}

// HelloWorld
func CamelCase(str string) string {
	camelCase := stringy.New(str)
	return camelCase.CamelCase()
}

// HelloWorlds
func CamelCasePlural(str string) string {
	camelCase := stringy.New(Plural(str))
	return camelCase.CamelCase()
}

func CleanStringLowered(str string) string {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	return str
}
func ContainsAll(str string, in ...string) bool {
	contains := stringy.New(str)
	return contains.ContainsAll(in...)
}

// This paragraph is too long...
func Excerpt(paragraph string, maxLength int) string {
	s := stringy.New(paragraph)
	return s.Tease(maxLength, "...")
}

// ExtractPrefix("Ben Walandow",3) => will return "Ben"
// ExtractPrefix("Ben",5) => will return "Ben"
func ExtractPrefix(s string, chars int) string {
	if chars > len(s) {
		return s
	}

	return s[:chars]
}

// ExtractSuffix("Ben Walandow",3) => will return "dow"
// ExtractSuffix("Ben",5) => will return "Ben"
func ExtractSuffix(s string, chars int) string {
	l := len(s)
	if chars > l {
		return s
	}
	s = s[l-chars : l]
	return s
}

// hello-world-morning
func KebabCase(str string) string {
	c := stringy.New(str)
	return c.KebabCase("?", "").ToLower()
}

// HELLO-WORLD-MORNING
func KebabCaseUpper(str string) string {
	c := stringy.New(str)
	return c.KebabCase("?", "").ToUpper()
}

// helloWorld
func LowerCamelCase(str string) string {
	cc := CamelCase(str)
	return stringy.New(cc).LcFirst()
}

// helloWorlds
func LowerCamelCasePlural(str string) string {
	cc := CamelCase(Plural(str))
	return stringy.New(cc).LcFirst()
}

// HelloWorlds
func Plural(str string) string {
	return pluralize.NewClient().Plural(str)
}

// Prefix makes sure string has been prefixed with a given string and avoids adding it again if it has
func Prefix(str string, prefix string) string {
	pre := stringy.New(str)
	return pre.Prefix(prefix)
}

func RemoveSpecialCharacter(str string) string {
	s := stringy.New(str)
	return s.RemoveSpecialCharacter()
}

func ReplaceFirst(str string, search string, replace string) string {
	s := stringy.New(str)
	return s.ReplaceFirst(search, replace)
}

func Singular(str string) string {
	return pluralize.NewClient().Singular(str)
}

// hello_world_morning
func SnakeCase(str string) string {
	snakeCase := stringy.New(str)
	return snakeCase.SnakeCase("?", "").ToLower()
}

// HELLO_WORLD_MORNING
func SnakeCaseUpper(str string) string {
	snakeCase := stringy.New(str)
	return snakeCase.SnakeCase("?", "").ToUpper()
}

func StrPad(str string, len int, ltr string, padStyle string) string {
	s := stringy.New(str)

	//It can be right / center
	//right: 	hello0000
	//both: 	00hello00
	if padStyle == "both" || padStyle == "right" {
		return s.Pad(0, ltr, padStyle)
	}

	//return 0000hello (default)
	return s.Pad(len, ltr, "left")
}

// 0000Hello
func StrPadLeft(str string, len int, ltr string) string {
	s := stringy.New(str)
	return s.Pad(len, ltr, "left")
}

// Hello0000
func StrPadRight(str string, len int, ltr string) string {
	s := stringy.New(str)
	return s.Pad(len, ltr, "right")
}

// Suffix makes sure string has been suffixed with a given string and avoids adding it again if it has.
func Suffix(str string, postfix string) string {
	suf := stringy.New(str)
	return suf.Suffix(postfix)
}

// ben + "'" => returned 'ben'
func StrQuoted(str string, by string) string {
	s := stringy.New(str)
	return s.Surround(by)
}

// Remove the last string based on suffix provided
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func TrimPrefix(s, prefix string) string {
	if len(prefix) > len(s) {
		return s
	}

	sPrefix := s[:len(prefix)]
	if sPrefix == prefix {
		return s[len(prefix):]
	}

	return s
}

func TrimAllWhiteSpaces(s string) string {
	s = strings.TrimSpace(s)
	pattern := regexp.MustCompile(`\s+`)
	return pattern.ReplaceAllString(s, " ")
}

// First letter capitalize
func UcFirst(str string) string {
	c := stringy.New(str)
	return c.UcFirst()
}

// All letter capitalize
func StrToUpper(str string) string {
	return strings.ToUpper(str)
}

// eb0ae4d9a87f4fd6b2e17d7cbab71853
// output: "eb0ae4d9-a87f-4fd6-b2e1-7d7cbab71853"
func HexToUUID(hex string) string {
	uuid := hex[0:8] + "-" + hex[8:12] + "-" + hex[12:16] + "-" + hex[16:20] + "-" + hex[20:]
	return strings.ToLower(uuid)
}

func ReplaceAll(text string, oldStr string, newStr string) string {
	return strings.ReplaceAll(text, oldStr, newStr)
}

// Random string with length input
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// Remove All WhiteSpace
func RemoveWhitespace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func NumberFormat(num int) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", num)
}

func NumberFormat64(num int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", num)
}

func FloatNumberFormat(num float64, precisions int) string {
	ps := ""
	if precisions >= 0 {
		ps = fmt.Sprintf(".%d", precisions)
	}
	p := message.NewPrinter(language.English)
	return p.Sprintf("%"+ps+"f", num)
}

func RemoveString(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func TrimArrayString(arr []string) []string {
	var newArr []string
	for _, a := range arr {
		newArr = append(newArr, strings.TrimSpace(a))
	}
	return newArr
}

func ExtractFloatFromAnyString(s string) float64 {
	if s == "" {
		return 0
	}
	re := regexp.MustCompile(`[0-9]+`)
	tr := strings.TrimSpace(s)
	submatchall := re.FindAllString(tr, -1)
	cleanPriceStr := ""
	for _, element := range submatchall {
		cleanPriceStr += strings.TrimSpace(element)
	}
	if len(cleanPriceStr) <= 0 {
		return 0
	}

	fn, _ := strconv.ParseFloat(cleanPriceStr, 64)

	return fn
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func Nl2br(str string) string {
	return strings.Replace(str, "\n", "<br />", -1)
}
