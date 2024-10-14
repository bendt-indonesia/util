package util

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/mail"
)

func ValidateStruct(i interface{}) error {
	if i == nil {
		return fmt.Errorf("No input to validate.")
	}
	v := validator.New()
	err := v.Struct(i)
	if err != nil {
		var msg string
		for _, err2 := range err.(validator.ValidationErrors) {
			msg += fmt.Sprintf("[%s] %s. ", err2.Field(), err2.Tag())
		}
		return fmt.Errorf(msg)
	}
	return nil
}

func ValidateStructMapped(i interface{}) map[string][]string {
	results := make(map[string][]string)
	if i == nil {
		results["EMPTY"] = []string{"No input to validate."}
		return results
	}
	v := validator.New()
	err := v.Struct(i)
	if err != nil {
		for _, err2 := range err.(validator.ValidationErrors) {
			if _, ex := results[err2.Field()]; !ex {
				results[err2.Field()] = []string{}
			}
			results[err2.Field()] = append(results[err2.Field()], fmt.Sprintf("[%s] %s", err2.Tag(), err2.Error()))
		}
	}
	return results
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
