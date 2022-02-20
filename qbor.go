package qbor

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

const (
	AND   = "AND"
	WHERE = "WHERE"
	OR    = "OR"
)

type IQuery interface {
	Where() (string, []interface{})
	Build(query string) (string, []interface{})
	Order(sort string) *Query
	Limit(limit int32) *Query
	Offset(offset int32) *Query

	Fetch(number int32, row string, args string) *Query //sqlserver
	OffsetRow(number int32, args string) *Query         //sqlserver
	Top(number int32, args string) *Query               //sqlserver
}

type Query struct {
	Filter    map[string]interface{}
	query     bytes.Buffer
	args      []interface{}
	sort      string
	limit     string
	offset    string
	fetch     string
	top       string
	offsetRow string
}

func NewQuery(filter map[string]interface{}) IQuery {
	for k, _ := range filter {
		s := strings.Split(k, "$")
		if len(s) > 1 {
			opr := strings.TrimSuffix(s[1], "!")
			if _, ok := operator[opr]; !ok {
				delete(filter, k)
			}
		} else {
			delete(filter, k)
		}
	}
	return &Query{
		Filter: filter,
		query:  bytes.Buffer{},
	}
}

//Where returns the where clause
func (q *Query) Where() (string, []interface{}) {
	var i int
	for k, v := range q.Filter {
		q.make(i, k, v)
		i++
	}
	//sql
	q.query.WriteString(q.sort)
	q.query.WriteString(q.limit)
	q.query.WriteString(q.offset)

	//sqlserver
	q.query.WriteString(q.offsetRow)
	q.query.WriteString(q.fetch)
	return q.query.String(), q.args
}

func (q *Query) Build(query string) (string, []interface{}) {
	//for add query top
	query = strings.Replace(strings.ToLower(query), `select`, fmt.Sprintf(` select %v `, q.top), 1)

	where, args := q.Where()
	query += where
	return query, args
}

func (q *Query) make(iteration int, k string, v interface{}) {
	arg := AND
	if iteration == 0 {
		arg = WHERE
	}

	fields := strings.Split(k, "$")
	columnName := fields[0]
	operatorName := fields[1]
	mustCompile := func(s string) bool {
		return s[len(s)-1:] == "!"
	}(operatorName)
	opr := q.translateOperator(operatorName)

	notNil, kind := q.isArgsNotNil(v)
	if mustCompile || notNil {

		switch opr {

		case operator["null"], operator["notnull"]:

			q.null(arg, columnName, opr)

		case operator["like"], operator["nlike"]:

			tmpArgs, ok := v.(string)
			if ok {
				q.like(arg, columnName, opr, ``)
				q.args = append(q.args, tmpArgs)
			}

		case operator["in"], operator["ni"]:
			q.in(arg, kind, v, columnName, opr)

		default:

			valueColumn, ok := q.isColumn(v)
			if ok {
				q.append(arg, columnName, opr, valueColumn)
			} else {
				q.append(arg, columnName, opr, "")
				q.args = append(q.args, v)
			}
		}
	} else {
		q.def(arg)
	}
}

func (q *Query) translateOperator(operatorName string) string {
	str := strings.Trim(operatorName, "!")
	return operator[strings.ToLower(str)]
}

func (q *Query) isArgsNotNil(i interface{}) (bool, reflect.Kind) {
	r := reflect.ValueOf(i)

	switch r.Kind() {
	case reflect.Slice:
		return r.Len() != 0, reflect.Slice
	case reflect.String:
		return r.String() != "", reflect.String
	case reflect.Int:
		return r.Int() != 0, reflect.Int
	case reflect.Int32:
		return r.Int() != 0, reflect.Int32
	case reflect.Int64:
		return r.Int() != 0, reflect.Int64
	case reflect.Float32:
		return r.Float() != 0, reflect.Float32
	case reflect.Float64:
		return r.Float() != 0, reflect.Float64
	default:
		return false, reflect.Invalid
	}
}

func (q *Query) isColumn(i interface{}) (string, bool) {
	col, ok := i.(string)
	if ok && strings.Contains(col, ":") {
		split := strings.Split(col, ":")
		if split[0] == "column" {
			return split[1], ok
		}
	}
	return col, false
}

func (q *Query) append(args, columnName, operator, valueColumn string) {
	if valueColumn == "" {
		valueColumn = q.isSliceOpr(operator)
	}
	q.query.WriteString(fmt.Sprintf(` %v %v %v %v`, args, columnName, operator, valueColumn))
}

func (q *Query) null(args, columnName, operator string) {
	q.query.WriteString(fmt.Sprintf(` %v %v %v`, args, columnName, operator))
}

//default value
func (q *Query) def(args string) {
	q.query.WriteString(fmt.Sprintf(` %v 1 = 1`, args))
}

func (q *Query) isSliceOpr(operator string) string {
	_, ok := isSliceOperator[operator]
	if ok {
		return `( ? )`
	}
	return `?`
}

func (q *Query) like(args, columnName, operator, valueColumn string) {
	valueColumn = `?`
	q.query.WriteString(fmt.Sprintf(` %v %v %v %v`, args, columnName, operator, valueColumn))
}

func (q *Query) in(args string, kind reflect.Kind, v interface{}, columnName, operator string) {
	if kind == reflect.Slice {
		s := reflect.ValueOf(v)
		if s.Len() > 0 {
			var smt string
			for j := 0; j < s.Len(); j++ {
				smt += `?,`
				q.kindIn(s, j)
			}
			q.append(args, columnName, operator, fmt.Sprintf(`(%v)`, smt[:len(smt)-1]))
		} else {
			q.append(args, columnName, operator, "")
			q.args = append(q.args, nil)
		}
	} else if kind == reflect.Invalid {
		q.append(args, columnName, operator, "")
		q.args = append(q.args, nil)
	} else {
		q.append(args, columnName, operator, "")
		q.args = append(q.args, v)
	}
}

func (q *Query) kindIn(s reflect.Value, index int) {
	switch s.Index(index).Kind() {
	case reflect.Int:
		q.args = append(q.args, s.Index(index).Int())
	case reflect.Int32:
		q.args = append(q.args, s.Index(index).Int())
	case reflect.Int64:
		q.args = append(q.args, s.Index(index).Int())
	case reflect.Float32:
		q.args = append(q.args, s.Index(index).Float())
	case reflect.Float64:
		q.args = append(q.args, s.Index(index).Float())
	default:
		q.args = append(q.args, s.Index(index).String())
	}
}

//Order ordering function
//multiple sort separators by comma
//for order by desc add prefix "-"
//example: Order("-column1,column2")
//result: order by column1 desc, column2 asc
func (q *Query) Order(sort string) *Query {
	if len(sort) > 0 {
		field := strings.Split(sort, ",")
		qSort := ` ORDER BY `
		for _, v := range field {
			if strings.HasSuffix(v, "-") {
				return q
			}
			sortType := func(str string) string {
				if strings.HasPrefix(str, "-") {
					return `desc`
				}
				return `asc`
			}
			qSort += strings.TrimPrefix(v, "-") + ` ` + sortType(v) + `,`
		}
		q.sort = qSort[:len(qSort)-1]
		return q
	}
	return q
}

//Limit limit function
func (q *Query) Limit(limit int32) *Query {
	qLimit := fmt.Sprintf(` limit %v`, limit)
	q.limit = qLimit
	return q
}

//Offset offset function
func (q *Query) Offset(offset int32) *Query {
	qLimit := fmt.Sprintf(` offset %v`, offset)
	q.offset = qLimit
	return q
}

//Fetch for fetching sql server
// example:Fetch(10,"FIRST","ROWS")
// default row is ROWS
// default fetch type is FIRST
func (q *Query) Fetch(number int32, row string, args string) *Query {
	if row == "" {
		row = `ROWS`
	}

	if args == "" {
		args = `FIRST`
	}

	qFetch := fmt.Sprintf(` FETCH %v %v %v ONLY`, args, number, row)
	q.fetch = qFetch
	return q
}

//Top for limit data sql server
//example: Top(10,PERCENT)
func (q *Query) Top(number int32, args string) *Query {
	qTop := fmt.Sprintf(` TOP %v %v`, number, args)
	q.top = qTop
	return q
}

//OffsetRow for offset data sql server
//example: OffsetRow(10)
//default offset type is ROWS
func (q *Query) OffsetRow(number int32, args string) *Query {
	if args == "" {
		args = `ROWS`
	}
	qOffset := fmt.Sprintf(` OFFSET %v %v`, number, args)
	q.offset = qOffset
	return q
}
