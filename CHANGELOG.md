# Changelog

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
