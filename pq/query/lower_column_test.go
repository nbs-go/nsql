package query_test

import (
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"github.com/nbs-go/nsql/test_utils"
	"testing"
	"time"
)

type BlogPost struct {
	Id        int64     `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"createdAt"`
}

func TestLowerSelect(t *testing.T) {
	s := schema.New(schema.FromModelRef(BlogPost{}))
	actual := query.Select(
		query.Lower(
			query.Column("title"),
		),
	).
		From(s).
		Build()
	expected := `SELECT LOWER("BlogPost"."title") FROM "BlogPost"`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestLowerAsSelect(t *testing.T) {
	s := schema.New(schema.FromModelRef(BlogPost{}))
	actual := query.Select(
		query.Lower(
			query.Column("title"),
			option.As("lowerTitle"),
		),
	).
		From(s).
		Build()
	expected := `SELECT LOWER("BlogPost"."title") AS "lowerTitle" FROM "BlogPost"`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestLowerOrder(t *testing.T) {
	s := schema.New(schema.FromModelRef(BlogPost{}))
	actual := query.Select(
		query.Column(query.AllColumns),
	).
		From(s).
		OrderByColumn(query.Lower(query.Column("title"))).
		Build()
	expected := `SELECT "BlogPost"."id", "BlogPost"."title", "BlogPost"."content", "BlogPost"."createdAt" FROM "BlogPost" ORDER BY LOWER("BlogPost"."title") ASC`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestLowerAsOrder(t *testing.T) {
	s := schema.New(schema.FromModelRef(BlogPost{}))
	lowerColumn := query.Lower(
		query.Column("title"),
		option.As("lowerTitle"),
	)
	actual := query.Select(
		query.Column(query.AllColumns),
		lowerColumn,
	).
		From(s).
		OrderByColumn(lowerColumn).
		Build()

	expected := `SELECT "BlogPost"."id", "BlogPost"."title", "BlogPost"."content", "BlogPost"."createdAt", LOWER("BlogPost"."title") AS "lowerTitle" FROM "BlogPost" ORDER BY "lowerTitle" ASC`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestLowerPanicAllColumn(t *testing.T) {
	defer test_utils.RecoverPanic(t, "Panic when ColumnWriter is using All Columns (*)", "nsql: all column (*) is not supported")()
	_ = query.Lower(query.Column(query.AllColumns))
}

func TestLowerPanicNil(t *testing.T) {
	defer test_utils.RecoverPanic(t, "Panic when ColumnWriter is nil", "nsql: column cannot be nil")()
	_ = query.Lower(nil)
}
