<a href="https://zerodha.tech"><img src="https://zerodha.tech/static/images/github-badge.svg" align="right" /></a>

# goyesql

> This package is based on [nleof/goyesql](https://github.com/nleof/goyesql) but is not compatible with it any more. This package introduces support for arbitrary tag types and changes structs and error types.

Parses a file and associate SQL queries to a map. Useful for separating SQL from code logic.

# Installation

```
$ go get -u github.com/knadh/goyesql/v2
```

## Usage

Create a file containing your SQL queries

```sql
-- queries.sql

-- name: list
-- some: param
-- some_other: param
SELECT *
FROM foo;

-- name: get
SELECT *
FROM foo
WHERE bar = $1;
```

And just call them in your code!

```go
queries := goyesql.MustParseFile("queries.sql")
// use queries["list"] with sql/database, sqlx ...
// queries["list"].Query is the parsed SQL query string
// queries["list"].Tags is the list of arbitrary tags (some=param, some_other=param)
```

## Scanning

Often, it's necessary to scan multiple queries from a SQL file, prepare them into \*sql.Stmt and use them throught the application. goyesql comes with a helper function that helps with this. Given a goyesql map of queries, it can turn the queries into prepared statements and scan them into a struct that can be passed around.

```go
type MySQLQueries struct {
	// This will be prepared.
	List *sql.Stmt `query:"list"`

	// This will not be prepared.
	Get  string    `query:"get"`
}

type MySQLxQueries struct {
	// These will be prepared.
	List *sqlx.Stmt `query:"list"`
	NamedList *sqlx.NamedStmt `query:"named_list"`

	// This will not be prepared.
	Get  string    `query:"get"`
}

var (
	q  MySQLQueries
	qx MySQLxQueries
)

// Here, db (*sql.DB) is your live DB connection.
err := goyesql.ScanToStruct(&q, queries, db)
if err != nil {
	log.Fatal(err)
}

// Here, db (*sqlx.DB) is your live DB connection.
err := goyesqlx.ScanToStruct(&qx, queries, db)
if err != nil {
	log.Fatal(err)
}

// Then, q.Exec(), q.QueryRow() etc.

```

## Embedding

You can use [stuffbin](https://github.com/knadh/stuffbin) and `ParseBytes()` for embedding SQL queries in your binary.
