package scan

import "fmt"

func ScanUint(label string, min, max *uint) uint {
	var scan uint
	promp := "Enter to skip"
	if min != nil || max != nil {
		if min != nil {
			promp = fmt.Sprintf("%d..", *min)
		} else {
			promp = "0.."
		}

		if max != nil {
			promp += fmt.Sprintf("%d", *max)
		} else {
			promp += "any"
		}
	}
	for true {
		fmt.Printf("%s [%s]: ", label, promp)
		fmt.Scanln(&scan)

		if min != nil || max != nil {
			if min != nil && scan < *min {
				fmt.Println("Minimum value is ", *min)
				continue
			}
			if max != nil && scan > *max {
				fmt.Println("Maximum value is ", *max)
				continue
			}
		}

		return scan
	}

	return scan
}

func ScanString(promptMsg string) string {
	var scan string

	for true {
		fmt.Printf(promptMsg)
		fmt.Scanln(&scan)
		break
	}

	return scan
}

func ScanBool(label string) bool {
	var scan string
	if label == "" {
		label = "Would you like to proceed [1/0]? "
	}
	if scan != "1" && scan != "0" {
		fmt.Printf(label)
		fmt.Scanln(&scan)
		if scan == "1" {
			return true
		} else if scan == "0" {
			return false
		}
	}
	return false
}
