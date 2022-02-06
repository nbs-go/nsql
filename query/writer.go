package query

import "github.com/nbs-go/nsql/schema"

// SchemaSetter must be implemented by part of query that may not require defining schema,
// but will be set later. For example, Selected fields can only set without defining schema and will be referred to
// schema that is defined in FROM
type SchemaSetter interface {
	TableGetter
	SetSchema(s *schema.Schema)
}

type AliasSetter interface {
	SetTableAs(as string)
}

type TableGetter interface {
	GetTableName() string
}

// SelectWriter must be implemented by part of query that will generate query in SELECT
type SelectWriter interface {
	SelectQuery() string
	SetFormat(format ColumnFormat)
	IsAllColumns() bool
	AliasSetter
	SchemaSetter
}

// FromWriter must be implemented by part of query that will generate query in FROM
type FromWriter interface {
	FromQuery() string
	TableGetter
}

type WhereWriter interface {
	WhereQuery() string
}

type OrderByWriter interface {
	OrderByQuery() string
	AliasSetter
	SchemaSetter
}

type ComparisonWhereWriter interface {
	AliasSetter
	SchemaSetter
}

type LogicalWhereWriter interface {
	GetConditions() []WhereWriter
	SetConditions(conditions []WhereWriter)
}

type ColumnWriter interface {
	ColumnQuery() string
	SetFormat(format ColumnFormat)
	AliasSetter
	SchemaSetter
}

type Expander interface {
	Expand(args ...interface{}) SelectWriter
}
