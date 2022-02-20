package qbor

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Query(t *testing.T) {
	query := NewQuery(map[string]interface{}{
		"name$eq": "test",
	})
	expectedArgs := []interface{}{"test"}
	expected := ` select   * from test where name = ?`
	q := `select * from test`

	q, args := query.Build(q)
	assert.Equal(t, expectedArgs, args)
	assert.Equal(t, expected, strings.ToLower(q))
}

func TestQuery_Where(t *testing.T) {

	tests := []struct {
		name   string
		fields map[string]interface{}
		want   string
		want1  []interface{}
	}{
		{
			name: "test query equal",
			fields: map[string]interface{}{
				"name$eq": "test",
			},
			want:  " where name = ?",
			want1: []interface{}{"test"},
		},
		{
			name: "test query not equal",
			fields: map[string]interface{}{
				"name$ne": "test",
			},
			want:  " where name != ?",
			want1: []interface{}{"test"},
		},
		{
			name: "test query lower than",
			fields: map[string]interface{}{
				"name$lt": "test",
			},
			want:  " where name < ?",
			want1: []interface{}{"test"},
		},
		{
			name: "test query lower than equal",
			fields: map[string]interface{}{
				"name$le": "test",
			},
			want:  " where name <= ?",
			want1: []interface{}{"test"},
		},
		{
			name: "test query greater than",
			fields: map[string]interface{}{
				"name$gt": "test",
			},
			want:  " where name > ?",
			want1: []interface{}{"test"},
		},
		{
			name: "test query greater than equal",
			fields: map[string]interface{}{
				"name$ge": "test",
			},
			want:  " where name >= ?",
			want1: []interface{}{"test"},
		},
		{
			name: "test query in",
			fields: map[string]interface{}{
				"name$in": []string{"test"},
			},
			want:  " where name in (?)",
			want1: []interface{}{"test"},
		},
		{
			name: "test query not in",
			fields: map[string]interface{}{
				"name$ni": []string{"test"},
			},
			want:  " where name not in (?)",
			want1: []interface{}{"test"},
		},
		{
			name: "test query like",
			fields: map[string]interface{}{
				"name$like": "%test%",
			},
			want:  " where name like ?",
			want1: []interface{}{"%test%"},
		},
		{
			name: "test query not like",
			fields: map[string]interface{}{
				"name$nlike": "%test%",
			},
			want:  " where name not like ?",
			want1: []interface{}{"%test%"},
		},
		{
			name: "test query null mandatory",
			fields: map[string]interface{}{
				"name$null!": "",
			},
			want:  " where name is null",
			want1: nil,
		},
		{
			name: "test query not null mandatory",
			fields: map[string]interface{}{
				"name$notnull!": "",
			},
			want:  " where name is not null",
			want1: nil,
		},
		{
			name: "test query equal not mandatory",
			fields: map[string]interface{}{
				"name$eq": "",
			},
			want:  " where 1 = 1",
			want1: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewQuery(tt.fields)
			got, got1 := q.Where()
			assert.Equalf(t, tt.want, strings.ToLower(got), "Where()")
			assert.Equalf(t, tt.want1, got1, "Where()")
		})
	}
}

func TestQuery_Order(t *testing.T) {
	q := NewQuery(map[string]interface{}{}).Order("name,-age")
	got, _ := q.Where()
	assert.Equal(t, " order by name asc,age desc", strings.ToLower(got))
}

func TestQuery_Limit(t *testing.T) {
	q := NewQuery(map[string]interface{}{}).Limit(10)
	got, _ := q.Where()
	assert.Equal(t, " limit 10", strings.ToLower(got))
}

func TestQuery_Offset(t *testing.T) {
	q := NewQuery(map[string]interface{}{}).Offset(10)
	got, _ := q.Where()
	assert.Equal(t, " offset 10", strings.ToLower(got))
}

func TestQuery_Fetch(t *testing.T) {
	q := NewQuery(map[string]interface{}{}).Fetch(10, "ROW", "NEXT")
	got, _ := q.Where()
	assert.Equal(t, " fetch next 10 row only", strings.ToLower(got))
}

func TestQuery_Top(t *testing.T) {
	expected := ` select  top 10 percent  * from test`
	query := `select * from test`
	q := NewQuery(map[string]interface{}{}).Top(10, "PERCENT")
	got, _ := q.Build(query)
	assert.Equal(t, expected, strings.ToLower(got))
}

func TestQuery_OffsetRow(t *testing.T) {
	q := NewQuery(map[string]interface{}{}).OffsetRow(10, "ROW")
	got, _ := q.Where()
	assert.Equal(t, " offset 10 row", strings.ToLower(got))
}
