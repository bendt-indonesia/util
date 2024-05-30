package util

import (
	"testing"
)

func TestValidatePhoneNo(t *testing.T) {

	var test []interface{}
	expectedResults := make([]bool, len(test))
	test = append(test, "62817179611")
	expectedResults = append(expectedResults, false)
	test = append(test, "6224548787845")
	expectedResults = append(expectedResults, false)
	test = append(test, "6287171874141")
	expectedResults = append(expectedResults, false)
	test = append(test, "01841959414")
	expectedResults = append(expectedResults, false)
	test = append(test, nil)
	expectedResults = append(expectedResults, false)
	test = append(test, "")
	expectedResults = append(expectedResults, false)
	test = append(test, "141")
	expectedResults = append(expectedResults, false)
	test = append(test, "gwgwgg")
	expectedResults = append(expectedResults, false)
	test = append(test, "(0817) 178611")
	expectedResults = append(expectedResults, true)
	test = append(test, "081748 4199419")
	expectedResults = append(expectedResults, true)

	for idx, r := range test {
		v := ValidatePhoneNo(r, []string{"62"})
		if expectedResults[idx] && v != nil {
			t.Errorf("Test %+v, expected %v, got %s", r, expectedResults[idx], v.Error())
		} else if expectedResults[idx] == false && v == nil {
			t.Errorf("Test %+v, expected %v (INVALID), got VALID PHONE!", r, expectedResults[idx])
		}
	}

}
