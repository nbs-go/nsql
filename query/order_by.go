package query

import "github.com/nbs-go/nsql/query/op"

type OrderBy struct {
	TableName string
	Column    string
	Direction op.SortDirection
}
