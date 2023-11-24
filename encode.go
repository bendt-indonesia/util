package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHMAC(signed string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(signed))
	return hex.EncodeToString(h.Sum(nil))
}

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

func ToJSON(v interface{}) string {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonStr)
}

func ToByte(v interface{}) []byte {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return jsonStr
}
