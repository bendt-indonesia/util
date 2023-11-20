package util

import (
	"fmt"
	"encoding/json"
	"math/rand"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"golang.org/x/crypto/bcrypt"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")

func NewUuid() string {
	uid := uuid.New().String()
	return FormatUid(uid)
}

func FormatUid(uid string) string {
	uid = strings.ReplaceAll(uid,"-","")
	uid = strings.ToUpper(uid)

	return uid
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

//HashPassword hashes given password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 11)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"file":  "util",
			"func":  "HashPassword",
			"error": err,
		}).Error("Unable to Hash Password")
	}

	return string(bytes)
}

//Check if phone begin with zero, then add country code
func PhonePrefix(phone string) string {
	if phone[:1] == "+" {
		phone = phone[1:]
	}
	if phone[:1] == "0" {
		phone = "62"+phone[1:]
	}
	return phone
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake  = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func RandomTimestamp() time.Time {
	randomTime := rand.Int63n(time.Now().Unix() - 94608000) + 94608000

	randomNow := time.Unix(randomTime, 0)

	return randomNow
}

func RandomTimestampStr() string {
	unix := time.Now().Unix()
	return strconv.FormatInt(unix, 10)
}

func RandomLetterString(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
	}
	return string(b)
}

func RandomNumberString(length int) string {
	const letterBytes = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
	}
	return string(b)
}


type T struct {
	A int
	B string
	C float64
}
//Still Bug
func FindAvailableSlug(slug string, existing []string) string {
	i := 0
	loop := true
	for loop {
		i++
		key := slug+"-"+fmt.Sprintf("%d",i)
		loop = funk.Contains(existing, key)
	}
	return slug+"-"+fmt.Sprintf("%d",i)
}

func GetFileNameWithoutExt(n string) string {
	return strings.TrimSuffix(n, filepath.Ext(n))
}

func GetFileExt(n string) string {
	return filepath.Ext(n)
}

func ComposeUploadFileName(n string) string {
	fileName := GetFileNameWithoutExt(n)
	ext := GetFileExt(n)
	ts := RandomTimestampStr()

	return fileName+"-"+ts+ext
}

func ComposeUploadFileNameV2(n string, aliasName *string) string {
	fileName := GetFileNameWithoutExt(n)
	if aliasName != nil {
		fileName = *aliasName
	}

	ext := GetFileExt(n)
	ts := RandomTimestampStr()

	return fileName+"-"+ts+ext
}

func ToJSON(v interface{}) string {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonStr)
}
