# Changelog

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
