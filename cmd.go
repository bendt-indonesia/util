package util

import (
	"fmt"
	"github.com/thoas/go-funk"
)

func ConfirmCmd(taskNo string, infoTxt string) bool {
	var confirmInput string
	fmt.Printf(taskNo + "[CONFIRMATION] Do you wish to perform updates on " + infoTxt + " [1/0]: ")
	for confirmInput == "" {
		fmt.Scanln(&confirmInput)
		if funk.ContainsString([]string{"1", "0"}, confirmInput) {
			if confirmInput == "1" {
				return true
			}
			return false
		} else {
			confirmInput = ""
		}
	}

	return false
}
