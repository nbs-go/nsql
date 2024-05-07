# Changelog

## v0.17.0

- feat(pg): Add Lower to generate LOWER() function query
- feat(pg): Add OrderByColumn to use ColumnWriter as input argument

## v0.16.0

- feat(pg/types): Add Postgis Point and NullPoint data type

## v0.15.0

- feat(json): Add helper function to scan JSON column from RDBMS
- feat(dsn): URL Encode password when generating DSN
- feat(dsn): Add PostgreSQL and MySQL dsn formatter

## v0.14.0

- feat(pg): Add BoolVar for writing boolean value in WHERE clause

## v0.13.0

- fix(mysql): Log option.As() deprecation upon function called in MySQL Query Builder
- fix(pg): Log option.As() deprecation upon function called in PostgreSQL Query Builder
- BREAKING CHANGE(mysql): Implement SchemaReference in MySQL
- BREAKING CHANGE(pg): Implement SchemaReference for PostgreSQL
- BREAKING CHANGE: Replace Table with SchemaReference as table reference for building query

## v0.12.0

- feat(pg): Add JSON Column query writer

## v0.11.0

- feat(mysql): Allow parsing Filters that does not have bound arguments to pass (e.g. is null / is not null)
- feat(pg): Allow parsing Filters that does not have bound arguments to pass (e.g. is null / is not null)

## v0.10.0

- feat(mysql): Add IsNull and IsNotNull conditional operator
- feat(pq): Add IsNull and IsNotNull conditional operator

## v0.9.0

- feat(mysql): Add select column as query
- feat(pq): Add select column as query

## v0.8.0

- feat(pq): Add utility to Parse postgres error

## v0.7.1

- fix(pq): Add RETURNING pk on generated Insert query

## 0.7.0

- feat(mysql): Add Int and Float comparison FilterParser
- feat(pq): Add Int and Float comparison FilterParser
- feat(mysql): Add Equal and Time comparison FilterParser
- feat(pq): Add Equal and Time comparison FilterParser

## 0.6.0

- feat(mysql): Add query builder for MySQL / MariaDB

## 0.5.1

- [CHANGED] Change Limit and Skip data type to int64

## 0.5.0

- [FIXED] Fix panic on skipped Join condition
- [ADDED] Add BindVar for dynamic compare value

## 0.4.0

- [ADDED] Add ResetSkip function
- [ADDED] Add ResetLimit function

## 0.3.2

- [ADDED] Separate tableAs in Column writers to fix column resolver after query has been build

## 0.3.1

- [FIXED] Fix passing arguments in From builder constructor

## 0.3.0

- [ADDED] Add ResetOrderBy to clear orderBy queries
- [ADDED] Add Select builder function
- [ADDED] Add From query builder constructor

## 0.2.0

- [ADDED] Add LikeFilter helper
- [ADDED] Add FilterBuilder for parsing query string into WHERE conditions
- [CHANGED] Allow empty condition on AND / OR logic
- [CHANGED] Re-structure query interface and operator to fix ambiguous package naming 
- [FIXED] Generate bind var for IN / NOT IN based on argument count

## 0.1.1

- [FIXED] Fix isExists query in Schema builder

## 0.1.0

- [ADDED] Add build IsExists query for schema
- [ADDED] Implement query.SelectWriter for whereCompareWriter
- [CHANGED] Change Select arguments with query.SelectWriter
- [CHANGED] Rename FilterColumns to Filter
- [ADDED] Add filter columns to ensure input columns are declared in schema
- [ADDED] Add Schema-based query builder
- [FIXED] Use FROM table if schema is not defined in Select Count
- [ADDED] Add Delete query builder
- [FIXED] Treat empty columns as All Columns on Select Builder
- [ADDED] Add Update query builder
- [ADDED] Insert query builder
- [FIXED] Fix unsorted join
- [ADDED] Add getter Insert and Update columns
- [CHANGED] Make schema fields as Read-only
- [ADDED] Add Join Query Builder
- [ADDED] Add Count query builder
- [ADDED] Add Where query builder
- [ADDED] Add Select Query Builder without WHERE condition
- [ADDED] Add Schema to represents Table in database
