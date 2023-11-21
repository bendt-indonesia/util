package util

import (
	"fmt"
	"github.com/thoas/go-funk"
	"reflect"
	"strings"

	"github.com/bendt-indonesia/util/enum"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func GoqDefaultOrderBy(do map[string]string) []exp.OrderedExpression {
	var orderByExp []exp.OrderedExpression
	for key, val := range do {
		val = strings.ToUpper(val)
		if val == "ASC" {
			orderByExp = append(orderByExp, goqu.C(key).Asc())
		} else {
			orderByExp = append(orderByExp, goqu.C(key).Desc())
		}
	}
	return orderByExp
}

// return expressions, errorCodeStr ST00016
func GoqOrderBy(s interface{}, do map[string]string) []exp.OrderedExpression {
	var orderByExp []exp.OrderedExpression

	//Check if DefaultOrderBy is exists then call
	if reflect.ValueOf(s).IsNil() {
		return GoqDefaultOrderBy(do)
	}

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
		if od == "ASC" {
			orderByExp = append(orderByExp, goqu.C(key).Asc())
		} else {
			orderByExp = append(orderByExp, goqu.C(key).Desc())
		}
	}

	if len(orderByExp) == 0 {
		return GoqDefaultOrderBy(do)
	}

	return orderByExp
}

// return expressions, errorCodeStr ("ST-00016")
// chain always use AND
func GoqWhereSyntaxQuery(s interface{}, binaryCols []string) ([]exp.Expression, string) {
	var whereExp []exp.Expression
	if reflect.ValueOf(s).IsNil() {
		return whereExp, ""
	}

	var err string
	var orWhereExp []exp.Expression

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
		fd := rt.Field(i).Tag.Get("db")
		tp := rt.Field(i).Type.String()
		if fd == "" || tp != "*model.SearchSyntax" {
			continue
		}

		con := reflect.ValueOf(row.Interface()).Elem().FieldByName("Connective")
		q := reflect.ValueOf(row.Interface()).Elem().FieldByName("Keywords").Elem().String()
		com := reflect.ValueOf(row.Interface()).Elem().FieldByName("Comparator")

		whereExp, orWhereExp, err = GoqAppendWhereExp(whereExp, orWhereExp, &fd, binaryCols, &com, &q, &con)
		if err != "" {
			return whereExp, err
		}
	}

	//OR where group expressions
	whereExp = append(whereExp, goqu.Or(
		orWhereExp...,
	))

	return whereExp, ""
}

func GoqWhereSyntaxQueryV2(s interface{}, whereExp []exp.Expression, orWhereExp []exp.Expression, binaryCols []string) ([]exp.Expression, string) {
	if reflect.ValueOf(s).IsNil() {
		return whereExp, ""
	}

	var err string

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
		fd := rt.Field(i).Tag.Get("db")
		tp := rt.Field(i).Type.String()
		if fd == "" || tp != "*model.SearchSyntax" {
			continue
		}

		con := reflect.ValueOf(row.Interface()).Elem().FieldByName("Connective")
		q := reflect.ValueOf(row.Interface()).Elem().FieldByName("Keywords").Elem().String()
		com := reflect.ValueOf(row.Interface()).Elem().FieldByName("Comparator")

		whereExp, orWhereExp, err = GoqAppendWhereExp(whereExp, orWhereExp, &fd, binaryCols, &com, &q, &con)
		if err != "" {
			return whereExp, err
		}
	}

	//OR where group expressions
	whereExp = append(whereExp, goqu.Or(
		orWhereExp...,
	))

	return whereExp, ""
}

func GoqWhereQuotedSyntaxQuery(s interface{}, binaryCols []string, quotedCols []string) ([]exp.Expression, string) {

	var whereExp []exp.Expression
	if reflect.ValueOf(s).IsNil() {
		return whereExp, ""
	}

	var err string
	var orWhereExp []exp.Expression

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
		fd := rt.Field(i).Tag.Get("db")
		tp := rt.Field(i).Type.String()
		if fd == "" || tp != "*model.SearchSyntax" {
			continue
		}

		if funk.ContainsString(quotedCols, fd) {
			fd = "`" + fd + "`"
		}

		con := reflect.ValueOf(row.Interface()).Elem().FieldByName("Connective")
		q := reflect.ValueOf(row.Interface()).Elem().FieldByName("Keywords").Elem().String()
		com := reflect.ValueOf(row.Interface()).Elem().FieldByName("Comparator")

		whereExp, orWhereExp, err = GoqAppendWhereExp(whereExp, orWhereExp, &fd, binaryCols, &com, &q, &con)
		if err != "" {
			return whereExp, err
		}
	}

	//OR where group expressions
	whereExp = append(whereExp, goqu.Or(
		orWhereExp...,
	))

	return whereExp, ""
}

func GoqAndGroupWhereSyntax(s interface{}) (exp.Expression, string) {
	var whereExp exp.Expression
	if reflect.ValueOf(s).IsNil() {
		return whereExp, ""
	}

	whereExps, ers := GoqWhereSyntaxQuery(s, []string{})
	if ers != "" {
		return whereExp, ers
	}

	whereExp = goqu.And(
		whereExps...,
	)

	return whereExp, ""
}

// return andWhereExpressions, orWhereExpressions, errorCodeStr ("ST-00016")
func GoqAppendWhereExp(andWhereExpression []exp.Expression, orExpressionLists []exp.Expression, fd *string, binaryCols []string, com *reflect.Value, q *string, con *reflect.Value) ([]exp.Expression, []exp.Expression, string) {

	//fmt.Println("=====================t=====================", con.IsNil())
	conn := enum.ConnectiveAnd
	if con.IsNil() || !enum.Connective(con.Elem().String()).IsValid() {
		conn = enum.ConnectiveAnd
	} else {
		conn = enum.Connective(con.Elem().String())
	}

	comparator := enum.ComparatorAnyWith
	//Query should be exist
	if com.IsNil() || !enum.Comparator(com.Elem().String()).IsValid() {
		comparator = enum.ComparatorAnyWith
	}

	field := *fd
	field = HexColName(field, binaryCols)

	keywords := *q
	keywords2 := ""
	comparator = enum.Comparator(com.Elem().String())

	switch comparator {
	case enum.ComparatorAnyWith:
		keywords = "%" + keywords + "%"
		break
	case enum.ComparatorBeginWith:
		keywords = keywords + "%"
		break
	case enum.ComparatorEndWith:
		keywords = "%" + keywords
		break
	case enum.ComparatorBetween, enum.ComparatorNotBetween:
		//fmt.Println("=======================between===================", keywords)
		params := strings.Split(keywords, ",")
		//fmt.Printf("%+v\n\n", params)
		if len(params) != 2 {
			return andWhereExpression, orExpressionLists, "ST00016"
		}
		keywords = params[0]
		keywords2 = params[1]
		break
	}
	//
	//fmt.Println(field, keywords, conn, comparator)
	//fmt.Println("=======================t===================", comparator)

	switch comparator {
	case enum.ComparatorAnyWith, enum.ComparatorBeginWith, enum.ComparatorEndWith:
		if conn == enum.ConnectiveAnd {
			andWhereExpression = append(andWhereExpression,
				goqu.And(
					goqu.L(field).ILike(keywords),
				),
			)
		} else {
			orExpressionLists = append(orExpressionLists,
				goqu.Or(
					goqu.L(field).ILike(keywords),
				),
			)
		}

		return andWhereExpression, orExpressionLists, ""
	case enum.ComparatorEqualsTo:
		exp := goqu.L(field+" = ?", keywords)
		if strings.ToLower(keywords) == "null" {
			exp = goqu.L(field + " IS NULL")
		}

		if conn == enum.ConnectiveAnd {
			andWhereExpression = append(andWhereExpression,
				goqu.And(exp),
			)
		} else {
			orExpressionLists = append(orExpressionLists,
				goqu.Or(exp),
			)
		}

		return andWhereExpression, orExpressionLists, ""
	case enum.ComparatorGreaterThan:
		if conn == enum.ConnectiveAnd {
			andWhereExpression = append(andWhereExpression,
				goqu.And(
					goqu.L(field+" = ?", keywords),
				),
			)
		} else {
			orExpressionLists = append(orExpressionLists,
				goqu.Or(
					goqu.L(field+" > ?", keywords),
				),
			)
		}
		return andWhereExpression, orExpressionLists, ""
	case enum.ComparatorGreaterThanEqualsTo:
		if conn == enum.ConnectiveAnd {
			andWhereExpression = append(andWhereExpression,
				goqu.And(
					goqu.L(field+" >= ?", keywords),
				),
			)
		} else {
			orExpressionLists = append(orExpressionLists,
				goqu.Or(
					goqu.L(field+" >= ?", keywords),
				),
			)
		}
		return andWhereExpression, orExpressionLists, ""
	case enum.ComparatorLessThan:
		if conn == enum.ConnectiveAnd {
			andWhereExpression = append(andWhereExpression,
				goqu.And(
					goqu.L(field+" < ?", keywords),
				),
			)
		} else {
			orExpressionLists = append(orExpressionLists,
				goqu.Or(
					goqu.L(field+" < ?", keywords),
				),
			)
		}
		return andWhereExpression, orExpressionLists, ""
	case enum.ComparatorLessThanEqualsTo:
		if conn == enum.ConnectiveAnd {
			andWhereExpression = append(andWhereExpression,
				goqu.And(
					goqu.L(field+" <= ?", keywords),
				),
			)
		} else {
			orExpressionLists = append(orExpressionLists,
				goqu.Or(
					goqu.L(field+" <= ?", keywords),
				),
			)
		}
		return andWhereExpression, orExpressionLists, ""

	case enum.ComparatorInclude:
		qs := strings.Split(keywords, ",")
		if conn == enum.ConnectiveAnd {
			andWhereExpression = append(andWhereExpression,
				goqu.And(
					goqu.L(field+" IN ?", qs),
				),
			)
		} else {
			orExpressionLists = append(orExpressionLists,
				goqu.Or(
					goqu.L(field+" IN ?", qs),
				),
			)
		}
		return andWhereExpression, orExpressionLists, ""

	case enum.ComparatorExclude:
		qs := strings.Split(keywords, ",")
		if conn == enum.ConnectiveAnd {
			andWhereExpression = append(andWhereExpression,
				goqu.And(
					goqu.L(field+" NOT IN ?", qs),
				),
			)
		} else {
			orExpressionLists = append(orExpressionLists,
				goqu.Or(
					goqu.L(field+" NOT IN ?", qs),
				),
			)
		}
		return andWhereExpression, orExpressionLists, ""

	case enum.ComparatorBetween:
		if conn == enum.ConnectiveAnd {
			andWhereExpression = append(andWhereExpression,
				goqu.And(
					goqu.L(field+" >= ?", keywords),
					goqu.L(field+" <= ?", keywords2),
				),
			)
		} else {
			orExpressionLists = append(orExpressionLists,
				goqu.Or(
					goqu.L(field+" >= ?", keywords),
					goqu.L(field+" <= ?", keywords2),
				),
			)
		}
		return andWhereExpression, orExpressionLists, ""

	case enum.ComparatorNotBetween:
		if conn == enum.ConnectiveAnd {
			andWhereExpression = append(andWhereExpression,
				goqu.And(
					goqu.L(field+" < ?", keywords),
					goqu.L(field+" > ?", keywords2),
				),
			)
		} else {
			orExpressionLists = append(orExpressionLists,
				goqu.Or(
					goqu.L(field+" < ?", keywords),
					goqu.L(field+" > ?", keywords2),
				),
			)
		}
		return andWhereExpression, orExpressionLists, ""
	}

	return andWhereExpression, orExpressionLists, ""
}

func GoqMergeWhereExpression(ex []exp.Expression, args ...exp.Expression) []exp.Expression {
	return append(ex, args...)
}
func GoqMergeWhereDeletedIsNull(ex []exp.Expression) []exp.Expression {
	return GoqMergeWhereExpression(ex,
		goqu.Ex{
			"deleted_at": nil,
		},
	)
}
func GoqMergeOrderByExpression(ex []exp.OrderedExpression, args ...exp.OrderedExpression) []exp.OrderedExpression {
	return append(ex, args...)
}
