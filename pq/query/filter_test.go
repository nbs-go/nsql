package query_test

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/parse"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"github.com/nbs-go/nsql/test_utils"
	"testing"
	"time"
)

func TestFilter(t *testing.T) {
	type Order struct {
		CreatedAt   time.Time `db:"createdAt"`
		UpdatedAt   time.Time `db:"updatedAt"`
		Id          int64     `db:"id"`
		OrderNumber string    `db:"orderNumber"`
		Status      string    `db:"status"`
		Version     int64     `db:"version"`
	}
	order := schema.New(schema.FromModelRef(Order{}))

	// Init filter function mapper
	ff := map[string]nsql.FilterParser{
		"status": func(queryValue string) (nsql.WhereWriter, []interface{}) {
			// Get arguments
			args := parse.IntArgs(queryValue)

			// Create filter
			w := query.In(query.Column("status", option.Schema(order)), len(args))

			return w, args
		},
		"orderNumber": func(queryValue string) (nsql.WhereWriter, []interface{}) {
			// args
			args := []interface{}{fmt.Sprintf(`%%%s%%`, queryValue)}

			// Create filter
			w := query.ILike(query.Column("orderNumber", option.Schema(order)))
			return w, args
		},
		"updatedAt": func(queryValue string) (nsql.WhereWriter, []interface{}) {
			args := parse.TimeArg(queryValue, "")
			if len(args) == 0 {
				return nil, nil
			}

			w := query.GreaterThanEqual(query.Column("updatedAt", option.Schema(order)))

			return w, args
		},
	}

	qs := map[string]string{
		"status":       "1,2,3",
		"orderNumber":  "1234",
		"customerName": "john",
		"createdAt":    "",
		"updatedAt":    "1644743059a",
	}

	// Create query
	b := query.NewFilter(qs, ff)
	q := b.Conditions().WhereQuery()

	// Compare string
	test_utils.CompareStringIn(t, "CORRECT QUERY FILTER", q,
		[]string{
			`"Order"."status" IN (?, ?, ?) AND "Order"."orderNumber" ILIKE ?`,
			`"Order"."orderNumber" ILIKE ? AND "Order"."status" IN (?, ?, ?)`,
		})

	test_utils.CompareInt(t, "CORRECT FILTER ARGUMENT COUNT ", 4, len(b.Args()))

	qs = map[string]string{
		"updatedAt": "1644743059a",
	}
	b = query.NewFilter(qs, ff)
	q = b.Conditions().WhereQuery()

	test_utils.CompareString(t, "EMPTY FILTER", q, "")
}
