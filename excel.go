package util

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ReadExcelFile(path string) (*excelize.File, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	return f, nil
}

func ReadSheets(f *excelize.File) []string {
	var sheets []string
	if f == nil {
		return sheets
	}

	for _, name := range f.GetSheetMap() {
		sheets = append(sheets, name)
	}

	return sheets
}
