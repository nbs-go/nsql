package query

import "github.com/nbs-go/nsql/pq/op"

type OrderBy struct {
	TableName string
	Column    string
	Direction op.SortDirection
}
