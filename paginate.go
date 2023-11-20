package util

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/thoas/go-funk"
)

const (
	DBDefaultLimit uint64 = 10
	DBDefaultMaxLimit uint64 = 150
	PageDefaultMaxLimit uint64 = 10
)

func CheckLastBeforeOffset(totalRecords uint64, last *uint, before *string, limit uint64, offset uint64) (uint64, uint64) {

	//check page position
	if last != nil {
		if before == nil {
			//last doang, before nya kosong

			if totalRecords < limit {
				//TEST CASE last 10 records with totalRecords of 5 (0-4)
				//totalRecords: 10 (0-9) Limit nya 11
				//no limit & no offset
				limit = 0
				offset = 0
			} else {
				//TEST CASE, last 10 with totalRecords of 295 (0-294)
				//TEST CASE, last 5 with totalRecords of 5 ( 0-4 )
				offset = totalRecords - limit
			}

		} else {
			//last dan before

			//TEST CASE, last 10, before 5, with totalRecords of 295 (0-294)
			//TEST CASE, last 10, before 1, with totalRecords of 295 (0-294)
			newOffset := int64(offset - limit)
			if newOffset < 0 {
				limit = limit - (limit - offset)
				offset = 0
			} else {
				//TEST CASE, last 10, before 25 ( 14 - 24 ), with totalRecords of 295 (0-294)
				offset = offset - limit - 1
			}

		}
	}

	return limit, offset
}

func CheckPrevNextPage(totalRecords uint64, limit uint64, offset uint64) (bool, bool) {
	hasNextPage := false
	hasPrevPage := false
	if offset != 0 {
		hasPrevPage = true
	}

	if (offset + limit) < totalRecords {
		hasNextPage = true
	}
	return hasPrevPage, hasNextPage
}

func CheckLimitOffset(totalRecords uint64, first *uint, after *string, last *uint, before *string) (uint64, uint64, string) {
	limit, offset, e := DecodeAllCursors(first, after, last, before)
	if e != "" {
		return 0, 0, e
	}
	limit, offset = CheckLastBeforeOffset(totalRecords, last, before, limit, offset)
	return limit, offset, ""
}

func CheckLimitOffsetNormal(limit *uint, offset *uint) (uint64, uint64) {
	var nLimit uint64
	var nOffset uint64 = 0

	if limit == nil || *limit < 0 {
		nLimit = DBDefaultLimit
	} else {
		if *limit > uint(DBDefaultMaxLimit) {
			nLimit = DBDefaultMaxLimit
		} else {
			nLimit = uint64(*limit)
		}
	}

	if offset != nil && *offset > 0 {
		nOffset = uint64(*offset)
	}

	return nLimit, nOffset
}

//return limitQuery string, offsetQuery string, stansiaErrorEnum string
func DecodeAllCursors(first *uint, after *string, last *uint, before *string) (uint64, uint64, string) {
	var limit uint64
	var offset uint64
	var err string

	if last == nil && first == nil {
		limit = DBDefaultLimit
	} else if last != nil {
		if int64(*last) < 1 {
			return 0, 0, "ST00010"
		}
		limit = uint64(*last)
	} else {
		if int64(*first) < 1 {
			return 0, 0, "ST00010"
		}
		limit = uint64(*first)
	}

	if limit > DBDefaultMaxLimit {
		limit = DBDefaultMaxLimit
	}

	if after == nil && before == nil {
		offset = 0
	} else if before != nil {
		_, _, offset, err = DecodeAndSplitCursor(*before, "|")
		if int64(offset) < 0 || err != "" {
			return limit, 0, err
		}
	} else if after != nil {
		_, _, offset, err = DecodeAndSplitCursor(*after, "|")
		if int64(offset) < 0 || err != "" {
			return limit, 0, err
		}
	}

	return limit, offset, ""
}

func DecodeCursor(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil || len(b) == 0 {
		return "", err
	}
	return string(b), nil
}

//return modelName, row.ID, offset, errorCode
func DecodeAndSplitCursor(d string, sep string) (string, string, uint64, string) {
	decoded, err := DecodeCursor(d)
	if err != nil {
		return "", "", 0, "ST00011"
	}

	arr := strings.Split(decoded, sep)
	if len(arr) != 3 {
		return "", "", 0, "ST00011"
	}

	offset, err := strconv.ParseUint(arr[2], 10, 64)
	if err != nil {
		return "", "", 0, "ST00011"
	}

	return arr[0], arr[1], offset, ""
}

func EncodeCursor(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func EncodeNodeCursor(m string, id string, offset uint64) string {
	offStr := strconv.FormatUint(offset, 10)
	compose := m + "|" + id + "|" + offStr

	return base64.StdEncoding.EncodeToString([]byte(compose))
}

func FilterFieldCtx(cp []string, mFieldCtx []string) string {
	f := funk.IntersectString(cp, mFieldCtx)

	return strings.Join(f, ", ")
}

func FilterFieldCtxWithHex(cp []string, mFieldCtx []string, hexedCols []string) string {
	f := funk.IntersectString(cp, mFieldCtx)

	res := ""
	for _, r := range f {
		if funk.ContainsString(hexedCols, r){
			res += "HEX(`"+r+"`) as "+r+", "
		} else {
			res += r+", "
		}
	}

	res = TrimSuffix(res, ", ")
	return res
}

func FilterFieldCtxWithHexAndQuoted(cp []string, mFieldCtx []string, hexedCols []string, quotedCols []string) string {
	f := funk.IntersectString(cp, mFieldCtx)

	res := ""
	for _, r := range f {
		if funk.ContainsString(hexedCols, r){
			res += "HEX(`"+r+"`) as "+r+", "
		} else if funk.ContainsString(quotedCols, r) {
			res += "`"+r+"`, "
		} else {
			res += r+", "
		}
	}

	res = TrimSuffix(res, ", ")
	return res
}

func AppendSelectQuery(query string, cols []string) string {
	for _, c := range cols {
		query += ", "+c
	}
	query += " "
	return query
}

func InQuery(s interface{}) string {
	rv := reflect.ValueOf(s)

	if rv.Kind() != reflect.Slice || rv.IsNil() {
		return ""
	}

	q := "("
	for i := 0; i < rv.Len(); i++ {
		q += fmt.Sprintf("%v,", rv.Index(i).Interface())
	}
	q = TrimSuffix(q, ",")
	q += ")"

	return q
}

func InStrQuery(s interface{}) string {
	rv := reflect.ValueOf(s)
	if rv.Kind() != reflect.Slice || rv.IsNil() {
		return ""
	}

	q := "("
	for i := 0; i < rv.Len(); i++ {
		q += fmt.Sprintf("\"%v\",", rv.Index(i).Interface())
	}
	q = TrimSuffix(q, ",")
	q += ")"

	return q
}

func ShopInQuery(s interface{}) string {
	rv := reflect.ValueOf(s)
	if rv.Kind() != reflect.Slice || rv.IsNil() {
		return ""
	}

	q := "("
	for i := 0; i < rv.Len(); i++ {
		shopId := fmt.Sprintf("\"%v\",", rv.Index(i).Interface())
		q += strings.ReplaceAll(shopId, "-", "")
	}
	q = TrimSuffix(q, ",")
	q += ")"

	return q
}

func LimitOffsetQuery(limit uint64, offset uint64) string {
	var q string
	if limit > 0 {
		q += " LIMIT " + fmt.Sprintf("%d", limit)
	}
	if offset > 0 && offset < 100000 {
		q += " OFFSET " + fmt.Sprintf("%d", offset)
	}
	return q
}

func GetDefaultOrderBy(s interface{}) string {
	cme := reflect.ValueOf(s).MethodByName("DefaultOrderBy")
	if cme.IsValid() {
		res := cme.Call([]reflect.Value{})
		return res[0].String()
	}
	return ""
}

func OrderByQuery(s interface{}) string {
	//Check if DefaultOrderBy is exists then call
	if reflect.ValueOf(s).IsNil() {
		return GetDefaultOrderBy(s)
	}

	var orderBy string

	rt := reflect.TypeOf(s).Elem()
	rv := reflect.ValueOf(s).Elem()
	num := rv.NumField()

	for i := 0; i < num; i++ {
		row := rv.Field(i)

		//Check if interface pointer is nil then continue
		//It means no sort present for that field
		if row.IsNil() {
			continue
		}

		//Check where name by using db tag
		key := rt.Field(i).Tag.Get("db")
		if key == "" {
			continue
		}

		od := strings.ToUpper(fmt.Sprintf("%v", row.Interface()))
		orderBy += " " + key + " " + od + ","
	}

	if len(orderBy) == 0 {
		return GetDefaultOrderBy(s)
	}

	orderBy = TrimSuffix(" ORDER BY"+orderBy, ",")
	return orderBy
}

func SelectQuery(s []string, t interface{}) string {
	tags := StructDBTags(t)
	f := funk.IntersectString(tags, s)

	return strings.Join(f, ", ")
}

func StoreQuery(tbn string, f []string) string {
	if len(f) == 0 {
		return ""
	}

	q := "INSERT INTO `" + tbn + "` ("
	//Columns
	for _, r := range f {
		if r == "id" {
			continue
		}
		q += "`" + r + "`,"
	}
	q = TrimSuffix(q, ",")

	q += ") VALUES("
	//Values with prefix :, for sqlx db tagging
	for _, r := range f {
		if r == "id" {
			continue
		}
		//UNHEX(REPLACE("3f06af63-a93c-11e4-9797-00505690773f", "-",""))
		if r == "shop_id" {
			q += "UNHEX(REPLACE("
		}
		q += ":" + r
		if r == "shop_id" {
			q += ", '-',''))"
		}
		q += ","
	}
	q = TrimSuffix(q, ",")

	q += ")"

	return q
}

func StoreQueryWithHexId(tbn string, f []string, hexedCols []string, skipCols []string) string {
	if len(f) == 0 {
		return ""
	}

	q := "INSERT INTO `" + tbn + "` ("
	//Columns
	for _, r := range f {
		if funk.ContainsString(skipCols, r){
			continue
		}

		q += "`" + r + "`,"
	}
	q = TrimSuffix(q, ",")

	q += ") VALUES("

	//Values with prefix :, for sqlx db tagging
	for _, r := range f {
		if funk.ContainsString(skipCols, r){
			continue
		}

		if funk.ContainsString(hexedCols, r) {
			q += "UNHEX(REPLACE("
		}
		q += ":" + r
		if funk.ContainsString(hexedCols, r) {
			q += ", '-',''))"
		}
		q += ","
	}
	q = TrimSuffix(q, ",")

	q += ")"

	return q
}

func StructDBTags(s interface{}) []string {
	var fields []string

	rt := reflect.TypeOf(s).Elem()
	num := rt.NumField()

	for i := 0; i < num; i++ {
		//Check where name by using db tag
		key := rt.Field(i).Tag.Get("db")
		if key == "" {
			continue
		}
		fields = append(fields, key)
	}
	return fields
}

func UpdateQuery(tbn string, s interface{}) (string, string) {

	//Check if DefaultOrderBy is exists then call
	if reflect.ValueOf(s).IsNil() {
		return "", ""
	}

	rt := reflect.TypeOf(s).Elem()
	rv := reflect.ValueOf(s).Elem()
	num := rv.NumField()

	q := "UPDATE `" + tbn + "` SET "
	c := 0
	for i := 0; i < num; i++ {
		row := rv.Field(i)
		col := rt.Field(i).Tag.Get("db")

		//Check where name by using db tag
		//It means no sort present for that field
		if col == "" || col == "id" {
			continue
		}

		//Check if interface pointer is nil then continue
		if row.Kind() == reflect.Ptr {
			if row.IsNil() {
				continue
			}
		}

		q += "`" + col + "` = :" + col + " ,"
		c++
	}
	q = TrimSuffix(q, ",")
	q += " WHERE id = :id"

	//1 for minimum updates: updated_by_id and updated_at
	if c <= 1 {
		return "", "ST00013"
	}
	return q, ""
}

func UpdateQueryWithHexCols(tbn string, s interface{}, hexedCols []string, skipCols []string, nullCols []string) (string, string) {

	//Check if DefaultOrderBy is exists then call
	if reflect.ValueOf(s).IsNil() {
		return "", ""
	}

	rt := reflect.TypeOf(s).Elem()
	rv := reflect.ValueOf(s).Elem()
	num := rv.NumField()

	q := "UPDATE `" + tbn + "` SET "
	c := 0
	for i := 0; i < num; i++ {
		row := rv.Field(i)
		col := rt.Field(i).Tag.Get("db")

		//Check where name by using db tag
		//It means no sort present for that field
		if col == "" || col == "id" || funk.ContainsString(skipCols, col) {
			continue
		}

		//Check if interface pointer is nil then continue
		if row.Kind() == reflect.Ptr {
			if row.IsNil() && !funk.ContainsString(nullCols, col){
				continue
			}
		}

		q += "`" + col + "` = "
		if funk.ContainsString(hexedCols, col) {
			q += "UNHEX(REPLACE("
		}
		q += ":" + col
		if funk.ContainsString(hexedCols, col) {
			q += ", '-',''))"
		}
		q += " ,"

		c++
	}

	q = TrimSuffix(q, ",")
	q += " WHERE id = "

	if funk.ContainsString(hexedCols, "id") {
		q += "UNHEX(REPLACE(:id, '-','')) "
	} else {
		q += ":id "
	}

	//1 for minimum updates: updated_by_id and updated_at
	if c <= 1 {
		return q, "ST00013"
	}
	return q, ""
}

func UpdateFieldsWithHexCols(s interface{}, hexedCols []string, skipCols []string, nullCols []string) (map[string]interface{}, string) {

	//Check if DefaultOrderBy is exists then call
	if reflect.ValueOf(s).IsNil() {
		return map[string]interface{}{}, ""
	}

	rt := reflect.TypeOf(s).Elem()
	rv := reflect.ValueOf(s).Elem()
	num := rv.NumField()

	cols := map[string]interface{}{}

	for i := 0; i < num; i++ {
		row := rv.Field(i)
		col := rt.Field(i).Tag.Get("db")

		//Check where name by using db tag
		//It means no sort present for that field
		if col == "" || col == "id" || funk.ContainsString(skipCols, col) {
			continue
		}

		//Check if interface pointer is nil then continue
		if row.Kind() == reflect.Ptr {
			if row.IsNil() && !funk.ContainsString(nullCols, col){
				continue
			}
		}

		if funk.ContainsString(hexedCols, col) {
			col = HexColName(col, hexedCols)
		}

		cols[col] = row.Elem().Interface()
	}

	return cols, ""
}

func WhereQuery(s interface{}) string {

	if reflect.ValueOf(s).IsNil() {
		return ""
	}

	var whereQuery string

	rt := reflect.TypeOf(s).Elem()
	rv := reflect.ValueOf(s).Elem()

	num := rv.NumField()

	for i := 0; i < num; i++ {
		row := rv.Field(i)

		//Check if interface pointer is nil then continue
		//It means no sort present for that field
		if row.IsNil() {
			continue
		}

		//Check where name by using db tag
		key := rt.Field(i).Tag.Get("db")
		if key == "" {
			continue
		}

		tp := rt.Field(i).Type.String()

		q := strings.ToUpper(fmt.Sprintf("%v", row.Elem().Interface()))
		if tp == "*string" {
			whereQuery += " AND " + key + " LIKE '%" + q + "%' "
		} else if tp == "*int" || tp == "*uint" || tp == "*int8" || tp == "*int64" || tp == "*int32" || tp == "*uint64" || tp == "*uint32" {
			whereQuery += " AND " + key + " = " + q + " "
		} else if tp == "*bool" {
			fmt.Printf("BOOL %v \n", q)
			whereQuery += " AND " + key + " = " + q + " "
		}

	}

	return whereQuery
}

func WhereQueryExcept(s interface{}, exceptCols []string) string {

	if reflect.ValueOf(s).IsNil() {
		return ""
	}

	var whereQuery string

	rt := reflect.TypeOf(s).Elem()
	rv := reflect.ValueOf(s).Elem()

	num := rv.NumField()

	for i := 0; i < num; i++ {
		row := rv.Field(i)

		//Check if interface pointer is nil then continue
		//It means no sort present for that field
		if row.IsNil() {
			continue
		}

		//Check where name by using db tag
		key := rt.Field(i).Tag.Get("db")
		if key == "" || funk.ContainsString(exceptCols, key) {
			continue
		}

		tp := rt.Field(i).Type.String()

		q := strings.ToUpper(fmt.Sprintf("%v", row.Elem().Interface()))
		if tp == "*string" {
			whereQuery += " AND " + key + " LIKE '%" + q + "%' "
		} else if tp == "*int" || tp == "*uint" || tp == "*int8" || tp == "*int64" || tp == "*int32" || tp == "*uint64" || tp == "*uint32" {
			whereQuery += " AND " + key + " = " + q + " "
		} else if tp == "*bool" {
			fmt.Printf("BOOL %v \n", q)
			whereQuery += " AND " + key + " = " + q + " "
		}

	}

	return whereQuery
}

func WhereShopQuery(shopId string, s interface{}) string {
	q := WhereQuery(s)
	q += " AND shop_id = UNHEX(\"" + shopId + "\") "

	return q
}

func AndWhere(whereQuery string, col string, operator string, value string) string {
	whereQuery += " AND " + col + " " + operator + " " + value
	return whereQuery
}

func AndVariantValueIdsQuery(whereQuery string, valueIds []*uint) string {
	if len(valueIds) == 0 {
		return whereQuery
	}

	inQuery := "("
	for _, v := range valueIds {
		inQuery += strconv.FormatUint(uint64(*v), 10) + ","
	}
	inQuery = TrimSuffix(inQuery, ",")
	inQuery += ")"

	whereQuery += "AND ( value1_id IN " + inQuery + " OR value2_id IN " + inQuery + ")"

	return whereQuery
}
