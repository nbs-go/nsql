package nsql

import (
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/schema"
)

// SchemaReference must be implemented by part of query that may not require defining schema,
// but will be set later. For example, Selected fields can only set without defining schema and will be referred to
// schema that is defined in FROM
type SchemaReference interface {
	TableGetter
	SchemaRefGetter
	SetSchema(s *schema.Schema)
}

type SchemaRefGetter interface {
	GetSchemaRef() schema.Reference
}

type AliasSetter interface {
	SetTableAs(as string)
}

type TableGetter interface {
	GetTableName() string
}

type ColumnGetter interface {
	GetColumn() string
}

// SelectWriter must be implemented by part of query that will generate query in SELECT
type SelectWriter interface {
	SelectQuery() string
	SetFormat(format op.ColumnFormat)
	IsAllColumns() bool
	AliasSetter
	SchemaReference
}

// FromWriter must be implemented by part of query that will generate query in FROM
type FromWriter interface {
	FromQuery() string
	Join(j JoinWriter)
	SchemaRefGetter
}

type WhereWriter interface {
	WhereQuery() string
}

type OrderByWriter interface {
	OrderByQuery() string
	AliasSetter
	SchemaReference
}

type WhereCompareWriter interface {
	GetVariable() VariableWriter
	SetVariable(v VariableWriter)
	ColumnGetter
	AliasSetter
	SchemaReference
}

type WhereLogicWriter interface {
	GetConditions() []WhereWriter
	SetConditions(conditions []WhereWriter)
}

type ColumnWriter interface {
	ColumnQuery() string
	SetFormat(format op.ColumnFormat)
	ColumnGetter
	AliasSetter
	SchemaReference
}

type Expander interface {
	Expand(args ...interface{}) SelectWriter
}

type JoinWriter interface {
	JoinQuery() string
	GetTableName() string
	GetIndex() int
	SetIndex(s int)
	SchemaRefGetter
}

type VariableWriter interface {
	VariableQuery() string
}
