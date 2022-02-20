## QBOR

### What is QBOR?

QBOR is a tool for simple query building with sqlx

### installation
```
 go get github.com/luqman-v1/qbor
```

### Operator 
| code | operator | description        |
|:--- |:---------|:-------------------|
|"eq" |     "=" | equal              |
|"ne" |     "!=" | not equal          |
|"lt" |     "<" | lower than         |
|"le" |     "<=" | lower than equal   |
|"gt" |     ">" | greater than       |
|"ge" |     ">=" | greater than equal |
|"in" |     "IN" | in                 |
|"ni" |     "NOT IN" | not in             |
|"like" |   "LIKE"   | like               |
|"nlike" |  "NOT LIKE"    | not like           |
|"null" |   "is null"   | null               |
|"notnull" |"is not null"      | not null           |
| ! |  "" | for mandatory field    |

### Usage
```go
package main

import (
    "github.com/jmoiron/sqlx"
    "github.com/luqman-v1/qbor"
)

var (
    Db *sqlx.DB
)

func main() {
    query := qbor.NewQuery(map[string]interface{}{
        "name$eq": "luqman",
        "age$eq": "", // because empty string, it will be treated as 1 = 1
        "hobby$eq!": "cooking", // because !, it will be treated as mandatory field
    }).Order("name,-age")

    q := `SELECT name FROM USERS`

    q, args := query.Build(q)

    //q = "SELECT name FROM USERS WHERE name = ? AND hobby = ? AND 1 = 1 order by name asc, age desc"
    //args = []interface{}{"luqman","cooking"}

    var data []Order
    Db.Select(&data, Db.Rebind(q), args...)
}

```

### API

- Order(sort string)
- Limit(limit int32) 
- Offset(offset int32)
- Fetch(number int32, row string, args string) //sqlserver
- OffsetRow(number int32, args string)  //sqlserver
- Top(number int32, args string) //sqlserver

