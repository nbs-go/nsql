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

## Usage

> TODO

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
