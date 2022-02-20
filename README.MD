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

* query builder

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

* BulkQuery

```go
package main

import (
    "github.com/jmoiron/sqlx"
    "github.com/luqman-v1/qbor"
    "log"
)

type testStruct struct {
    Name string
    Age  int
}

var (
    Db *sqlx.DB
)

func main(){
    //model struct must be same order as query string
    q := `INSERT INTO USERS (name, age) VALUES {{values}}`
    bulk := qbor.NewBulkQuery()
    testStructs := []testStruct{
        {"John", 20},
        {"Ana", 21},
    }
    q, args, err := bulk.BulkQuery(testStructs, q)
    
    //q  = INSERT INTO USERS (name, age) VALUES (?,?), (?,?)
    //args = []interface{}{"John",20,"Ana",21}
    
    if err != nil {
        log.Println("Error: ", err)
    }
    Db.ExecContext(ctx, Db.Rebind(q), args...)
}


```


### API

- Order(sort string)
- Limit(limit int32) 
- Offset(offset int32)
- Fetch(number int32, row string, args string) //sqlserver
- OffsetRow(number int32, args string)  //sqlserver
- Top(number int32, args string) //sqlserver
- BulkQuery(model interface{}, query string) (string, []interface{}, error)



License
----

MIT

