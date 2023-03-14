# nsql

[![Go Report Card](https://goreportcard.com/badge/github.com/nbs-go/nsql)](https://goreportcard.com/report/github.com/nbs-go/nsql)
[![GitHub license](https://img.shields.io/github/license/nbs-go/nsql)](https://github.com/nbs-go/nsql/blob/master/LICENSE)

SQL Utility for Golang. Compatible with [`jmoiron/sqlx`](https://github.com/jmoiron/sqlx). Features:
- Query Builder (PostgreSQL / MySQL)

## Installing

> **WARNING**
>
> API is not yet stable and we might introduce breaking changes until we reached version v1.0.0. See [Breaking Changes](#breaking-changes) section for deprecation notes.

```shell
go get -u github.com/nbs-go/nsql
```

## Features

### DBMS Compatible:
 1. PostgreSQL
 2. MySQL


## Usage

nsql is SQL Utility to generate query as you want, this utility compatible with compatible with [`jmoiron/sqlx`](https://github.com/jmoiron/sqlx) as executor, and here several step to go.

1. Create your own struct:
```go
type Person struct {
    CreatedAt time.Time `db:"createdAt"`
    UpdatedAt time.Time `db:"updatedAt"`
    Id        int64     `db:"id"`
    FullName  string    `db:"fullName"`
}
```

2. Generated your schema:
```go
var person = schema.New(schema.FromModelRef(Person{}))
```

3. Get your query:
```go
var q = query.Select(query.Column("*")).From(person).Build()
```

Result query generated should be:
```postgresql
SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"
```


### Quick Start :

Example using postgres query:
```go
package main

import (
    "fmt"
    "time"
    "github.com/nbs-go/nsql/pq/query"
    "github.com/nbs-go/nsql/schema"
)

func main() {
    type Person struct {
        CreatedAt time.Time `db:"createdAt"`
        UpdatedAt time.Time `db:"updatedAt"`
        Id        int64     `db:"id"`
        FullName  string    `db:"fullName"`
    }

var person = schema.New(schema.FromModelRef(Person{}))

// SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"
var q = query.Select(query.Column("*")).From(person).Build()

    fmt.Println(q)

// SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" WHERE "Person"."fullName" = ?
q = query.Select(query.Column("*")).
    From(person).
    Where(query.Equal(query.Column("fullName"))).
    Build()

    fmt.Println(q)

// UPDATE "Person" SET "fullName" = :fullName WHERE "id" = :id
q = query.Update(person, "fullName").
    Where(query.Equal(query.Column("id"))).
    Build()

    fmt.Println(q)

// DELETE FROM "Person" WHERE "id" = ?
q = query.Delete(person).Build()

    fmt.Println(q)
}
```

## Breaking Changes

### v0.13.0

- `From()` and `Join()` option setter `option.As()` is deprecated. Use `schema.New()` option setter `schema.As()` instead.
    ```
    // <= v0.12.0
    person := schema.New(schema.FromModelRef(Person{}))
    query.Select(query.Column("*")).
        From(person, option.As("p")) // <- Set table alias with option.As()
  
    // >= v0.13.0
    person := schema.New(schema.FromModelRef(Person{}), schema.As("p")) // <- Set table alias with schema.As()
    query.Select(query.Column("*")).From(person)
    ```
## Contributors

<a href="https://github.com/nbs-go/nsql/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=nbs-go/nsql" alt="contributors" />
</a>
