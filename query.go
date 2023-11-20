package util

func HexColName(col string, binCols []string) string {
	for _, check := range binCols {
		if check == col {
			return "HEX(`" + col + "`)"
		}
	}
	return col
}

func HexColAsName(col string, binCols []string) string {
	for _, check := range binCols {
		if check == col {
			return "HEX(`" + col + "`) as `"+col+"`"
		}
	}
	return "`" + col + "`"
}

func HexValues(values []string) []string {

	for idx, val := range values {
		values[idx] = "HEX('" + val + "')"
	}
	return values
}
