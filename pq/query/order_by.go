package query

import (
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
)

type orderByWriter struct {
	query.ColumnWriter
	direction op.SortDirection
}

func (o *orderByWriter) OrderByQuery() string {
	var direction string
	if o.direction == op.Descending {
		direction = "DESC"
	} else {
		direction = "ASC"
	}
	return o.ColumnQuery() + " " + direction
}
