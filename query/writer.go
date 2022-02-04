package query

type SelectFieldWriter interface {
	SelectQuery() string
	SetTableAlias(alias string)
	SetMode(mode ColumnMode)
	GetTableName() string
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
