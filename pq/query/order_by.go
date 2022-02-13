package query

import (
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
)

type orderByWriter struct {
	nsql.ColumnWriter
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
