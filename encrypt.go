package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/linkedin/goavro"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

func GenerateHMAC(signed string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(signed))
	return hex.EncodeToString(h.Sum(nil))
}

func EncryptSHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
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

func UnescapeUnicodeCharactersInJSON(_jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(_jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func DecodeAvro(avroData []byte, schema string) map[string]interface{} {
	// Create an Avro codec for the message
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		return nil
	}

	// Decode Avro data to a map
	decoded, _, err := codec.NativeFromBinary(avroData)
	if err != nil {
		return nil
	}

	// Convert the decoded data to a map
	data, ok := decoded.(map[string]interface{})
	if !ok {
		return nil
	}

	return data
}

func EncodeAvro(messageMap map[string]interface{}, schema string) []byte {
	// Create an Avro codec for the message
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		return nil
	}
	// Encode the message to Avro binary format
	avroData, err := codec.BinaryFromNative(nil, messageMap)
	if err != nil {
		return nil
	}

	return avroData
}
