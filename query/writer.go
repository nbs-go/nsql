package query

type SelectFieldWriter interface {
	GetTableName() string
	SelectQuery() string
	SetMode(mode ColumnMode)
	SetTableAlias(alias string)
}

type FromWriter interface {
	FromQuery() string
}

type WhereWriter interface {
	WhereQuery() string
}

type ComparisonWhereWriter interface {
	SetTableAlias(alias string)
	GetTableName() string
}

type LogicalWhereWriter interface {
	GetConditions() []WhereWriter
}
