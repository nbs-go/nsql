package query_test

import (
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/parse"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"github.com/nbs-go/nsql/test_utils"
	"strings"
	"testing"
	"time"
)

type Order struct {
	CreatedAt   time.Time `db:"createdAt"`
	UpdatedAt   time.Time `db:"updatedAt"`
	Id          int64     `db:"id"`
	OrderNumber string    `db:"orderNumber"`
	Status      string    `db:"status"`
	Version     int64     `db:"version"`
}

var order = schema.New(schema.FromModelRef(Order{}))

func TestFilter(t *testing.T) {
	// Init filter function mapper
	ff := map[string]nsql.FilterParser{
		"status": func(queryValue string) (nsql.WhereWriter, []interface{}) {
			// Get arguments
			args := parse.IntArgs(queryValue)

			// Create filter
			w := query.In(query.Column("status", option.Schema(order)), len(args))

			return w, args
		},
		"orderNumber": query.LikeFilter("orderNumber", op.LikeSubString, option.Schema(order)),
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

	// Test #1
	test_utils.CompareStringIn(t, "CORRECT QUERY FILTER", q,
		[]string{
			`"Order"."status" IN (?, ?, ?) AND "Order"."orderNumber" ILIKE ?`,
			`"Order"."orderNumber" ILIKE ? AND "Order"."status" IN (?, ?, ?)`,
		})

	// Test #2
	test_utils.CompareInt(t, "CORRECT FILTER ARGUMENT COUNT ", 4, len(b.Args()))

	qs = map[string]string{
		"updatedAt": "1644743059a",
	}
	b = query.NewFilter(qs, ff)
	q = b.Conditions().WhereQuery()

	// Test #3
	test_utils.CompareString(t, "EMPTY FILTER", q, "")

	// Test #4
	qs = map[string]string{
		"orderNumber": "001",
	}
	b = query.NewFilter(qs, ff)
	test_utils.CompareInterfaceArray(t, "CORRECT LIKE ARGUMENT", b.Args(), []interface{}{"%001%"})
}

func TestLikeFilter_Prefix(t *testing.T) {
	// Prepare parser
	ff := map[string]nsql.FilterParser{
		"orderNumber": query.LikeFilter("orderNumber", op.LikePrefix, option.Schema(order)),
	}

	qs := map[string]string{
		"orderNumber": "001",
	}
	b := query.NewFilter(qs, ff)
	test_utils.CompareInterfaceArray(t, "CORRECT LIKE ARGUMENT", b.Args(), []interface{}{"%001"})
}

func TestLikeFilter_Suffix(t *testing.T) {
	// Prepare parser
	ff := map[string]nsql.FilterParser{
		"orderNumber": query.LikeFilter("orderNumber", op.LikeSuffix, option.Schema(order)),
	}

	qs := map[string]string{
		"orderNumber": "001",
	}
	b := query.NewFilter(qs, ff)
	test_utils.CompareInterfaceArray(t, "CORRECT LIKE ARGUMENT", b.Args(), []interface{}{"001%"})
}

func TestLikeFilter_Exact(t *testing.T) {
	// Prepare parser
	ff := map[string]nsql.FilterParser{
		"orderNumber": query.LikeFilter("orderNumber", op.LikeExact, option.Schema(order)),
	}

	qs := map[string]string{
		"orderNumber": "001",
	}
	b := query.NewFilter(qs, ff)
	test_utils.CompareInterfaceArray(t, "CORRECT LIKE ARGUMENT", b.Args(), []interface{}{"001"})
}

func TestEqualFilter_Args(t *testing.T) {
	// Prepare parser
	s := schema.New(schema.FromModelRef(new(Person)))
	ff := map[string]nsql.FilterParser{
		"id": query.EqualFilter(s, "id"),
	}

	qs := map[string]string{
		"id": "1234",
	}
	b := query.NewFilter(qs, ff)

	// Assert
	actual := b.Args()
	expected := []interface{}{"1234"}

	if len(actual) != len(expected) {
		t.Errorf("FAILED\n  > got different values: %v, expected: %v", actual, expected)
		return
	}

	for i, v := range actual {
		if expected[i] != v {
			t.Errorf("FAILED\n  > got different values: %v, expected: %v", actual, expected)
			return
		}
	}
}

func TestTimeBetweenFilter(t *testing.T) {
	// Prepare parser
	s := schema.New(schema.FromModelRef(new(Person)))
	ff := map[string]nsql.FilterParser{
		"fromCreatedAt":        query.TimeGreaterThanEqualFilter(s, "createdAt"),
		"toCreatedAt":          query.TimeLessThanEqualFilter(s, "createdAt"),
		"invalidFromCreatedAt": query.TimeGreaterThanEqualFilter(s, "createdAt"),
		"invalidToCreatedAt":   query.TimeLessThanEqualFilter(s, "createdAt"),
	}

	qs := map[string]string{
		"fromCreatedAt":        "2022-01-01T00:00:00+07:00",
		"toCreatedAt":          "2022-01-31T00:00:00+07:00",
		"invalidFromCreatedAt": "NOT A TIME",
		"invalidToCreatedAt":   "NOT A TIME",
	}
	b := query.NewFilter(qs, ff)

	// Assert
	actual := b.Args()
	expected := map[interface{}]bool{
		time.Unix(1640970000, 0).UTC(): true,
		time.Unix(1643562000, 0).UTC(): true,
	}

	if len(actual) != len(expected) {
		t.Errorf("FAILED\n  > got different values: %v, expected: %v", actual, expected)
		return
	}

	for _, v := range actual {
		if _, ok := expected[v]; !ok {
			t.Errorf("FAILED\n  > got different values: %v, expected: %v", actual, expected)
			return
		}
	}
}

func TestIntBetweenFilter(t *testing.T) {
	// Define table
	type OrderItem struct {
		CreatedAt time.Time `db:"createdAt"`
		UpdatedAt time.Time `db:"updatedAt"`
		Id        int64     `db:"id"`
		Qty       int64     `db:"qty"`
	}

	// Prepare parser
	s := schema.New(schema.FromModelRef(new(OrderItem)))
	ff := map[string]nsql.FilterParser{
		"moreThanQty":        query.IntGreaterThanEqualFilter(s, "qty"),
		"lessThanQty":        query.IntLessThanEqualFilter(s, "qty"),
		"invalidMoreThanQty": query.IntGreaterThanEqualFilter(s, "qty"),
		"invalidLessThanQty": query.IntLessThanEqualFilter(s, "qty"),
	}

	qs := map[string]string{
		"moreThanQty":        "10",
		"lessThanQty":        "20",
		"invalidMoreThanQty": "A",
		"invalidLessThanQty": "B",
	}
	b := query.NewFilter(qs, ff)

	// Assert
	actual := b.Args()
	expected := map[interface{}]bool{
		int64(10): true,
		int64(20): true,
	}

	if len(actual) != len(expected) {
		t.Errorf("FAILED\n  > got different values: %v, expected: %v", actual, expected)
		return
	}

	for _, v := range actual {
		if _, ok := expected[v]; !ok {
			t.Errorf("FAILED\n  > got different values: %v, expected: %v", actual, expected)
			return
		}
	}
}

func TestFloatBetweenFilter(t *testing.T) {
	// Define table
	type Task struct {
		CreatedAt time.Time `db:"createdAt"`
		UpdatedAt time.Time `db:"updatedAt"`
		Id        int64     `db:"id"`
		Progress  float64   `db:"progress"`
	}

	// Prepare parser
	s := schema.New(schema.FromModelRef(new(Task)))
	ff := map[string]nsql.FilterParser{
		"moreThanProgress":        query.FloatGreaterThanEqualFilter(s, "progress"),
		"lessThanProgress":        query.FloatLessThanEqualFilter(s, "progress"),
		"invalidMoreThanProgress": query.FloatGreaterThanEqualFilter(s, "progress"),
		"invalidLessThanProgress": query.FloatLessThanEqualFilter(s, "progress"),
	}

	qs := map[string]string{
		"moreThanProgress":        "25.5",
		"lessThanProgress":        "50.7",
		"invalidMoreThanProgress": "A",
		"invalidLessThanProgress": "B",
	}
	b := query.NewFilter(qs, ff)

	// Assert
	actual := b.Args()
	expected := map[interface{}]bool{25.5: true, 50.7: true}

	if len(actual) != len(expected) {
		t.Errorf("FAILED\n  > got different values: %v, expected: %v", actual, expected)
		return
	}

	for _, v := range actual {
		if _, ok := expected[v]; !ok {
			t.Errorf("FAILED\n  > got different values: %v, expected: %v", actual, expected)
			return
		}
	}
}

func TestStaticArgFilter(t *testing.T) {
	// Init filter function mapper
	ff := map[string]nsql.FilterParser{
		"isActive": func(queryValue string) (nsql.WhereWriter, []interface{}) {
			if strings.ToLower(queryValue) != "true" {
				return nil, nil
			}

			// Create filter
			w := query.IsNotNull(query.Column("statusId", option.Schema(order)))

			return w, nil
		},
	}

	// Get actual value
	qs := map[string]string{
		"isActive": "true",
	}
	b := query.NewFilter(qs, ff)
	actual := b.Conditions().WhereQuery()

	// Assert
	expected := `"Order"."statusId" IS NOT NULL`
	if actual != expected {
		t.Errorf("FAILED\n  > got different values: %v, expected: %v", actual, expected)
		return
	}
}
