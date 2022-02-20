package qbor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string
	Age  int
}

func TestBulkQuery_BulkQuery(t *testing.T) {
	q := `INSERT INTO USERS (name, age) VALUES {{values}}`
	b := NewBulkQuery()
	testStructs := []testStruct{
		{"John", 20},
	}
	q, args, err := b.BulkQuery(testStructs, q)

	assert.Nil(t, err)
	assert.Equal(t, `INSERT INTO USERS (name, age) VALUES ( ?, ?)`, q)
	assert.Equal(t, []interface{}{"John", 20}, args)
}

func TestBulkQuery_structEmpty(t *testing.T) {
	q := `INSERT INTO USERS (name, age) VALUES {{values}}`
	b := NewBulkQuery()
	var testStructs []testStruct
	q, args, err := b.BulkQuery(testStructs, q)

	assert.NotNil(t, err)
	assert.Equal(t, ``, q)
	assert.Nil(t, args)
}
