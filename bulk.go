package qbor

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type IBulkQuery interface {
	BulkQuery(model interface{}, query string) (string, []interface{}, error)
}

type BulkQuery struct {
}

func NewBulkQuery() IBulkQuery {
	return &BulkQuery{}
}

// BulkQuery returns a query string and a list of parameters
// for the given model and query.
// model struct must be same order as query string
// and must have same field names as query string
func (b *BulkQuery) BulkQuery(model interface{}, query string) (string, []interface{}, error) {
	result := reflect.ValueOf(model)
	if result.Len() <= 0 {
		return "", nil, errors.New("data struct is empty")
	}
	buffParam := `(%v),`
	queryBuff := new(bytes.Buffer)
	var args []interface{}
	paramValue := ``
	for i := 0; i < result.Len(); i++ {
		var strQuery string
		for j := 0; j < result.Index(i).NumField(); j++ {
			args = append(args, result.Index(i).Field(j).Interface())
			strQuery += ` ?`
			if j < result.Index(i).NumField()-1 {
				strQuery += `,`
			}
		}
		paramValue += fmt.Sprintf(buffParam, strQuery)
	}
	queryBuff.WriteString(paramValue)
	queryString := queryBuff.String()
	return strings.ReplaceAll(query, "{{values}}", queryString[:len(queryString)-1]), args, nil
}
