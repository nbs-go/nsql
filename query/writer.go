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
