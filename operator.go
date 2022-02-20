package qbor

var operator = map[string]string{
	"eq":      "=",
	"ne":      "!=",
	"lt":      "<",
	"le":      "<=",
	"gt":      ">",
	"ge":      ">=",
	"in":      "IN",
	"ni":      "NOT IN",
	"like":    "LIKE",
	"nlike":   "NOT LIKE",
	"null":    "is null",
	"notnull": "is not null",
}

var isSliceOperator = map[string]string{
	"in": "in",
	"ni": "ni",
}
