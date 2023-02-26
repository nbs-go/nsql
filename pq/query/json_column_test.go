package query_test

import (
	"encoding/json"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"testing"
	"time"
)

type Content struct {
	CreatedAt time.Time       `db:"createdAt"`
	UpdatedAt time.Time       `db:"updatedAt"`
	Id        int64           `db:"id"`
	Content   json.RawMessage `db:"content"`
}

type ContentTag struct {
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
	Id        int64     `db:"id"`
	ContentId int64     `db:"contentId"`
	Tag       string    `db:"tag"`
}

var content = schema.New(schema.FromModelRef(Content{}))
var contentTag = schema.New(schema.FromModelRef(ContentTag{}))

func TestSelectJsonColumn1(t *testing.T) {
	// Process actual value
	actual := query.Select(query.JsonColumn("content.title")).
		From(content).
		Build()

	// Compare
	expected := `SELECT "Content"."content"->>'title' FROM "Content"`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectJsonColumn1_As(t *testing.T) {
	// Process actual value
	actual := query.Select(query.JsonColumn("content.title", option.As("title"))).
		From(content, option.As("c")).
		Build()

	// Compare
	expected := `SELECT ("c"."content"->>'title') AS "title" FROM "Content" AS "c"`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectJsonColumn1_WhereLike(t *testing.T) {
	// Process actual value
	actual := query.Select(query.JsonColumn("content.title")).
		From(content).
		Where(query.Like(query.JsonColumn("content.title"))).
		Build()

	// Compare
	expected := `SELECT "Content"."content"->>'title' FROM "Content" WHERE "Content"."content"->>'title' LIKE ?`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectJsonColumn2(t *testing.T) {
	// Process actual value
	actual := query.Select(query.JsonColumn("content.image.fileName", option.Schema(content))).
		From(content).
		Build()

	// Compare
	expected := `SELECT "Content"."content"->'image'->>'fileName' FROM "Content"`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectJsonColumn2_As(t *testing.T) {
	// Process actual value
	actual := query.Select(query.JsonColumn("content.image.fileName", option.As("fileName"))).
		From(content).
		Build()

	// Compare
	expected := `SELECT ("Content"."content"->'image'->>'fileName') AS "fileName" FROM "Content"`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectJsonColumn_FlagSkip(t *testing.T) {
	// Process actual value
	actual := query.Select(
		query.Column("content"),
		query.JsonColumn("noContent.title"),
	).
		From(content).
		Build()

	// Compare
	expected := `SELECT "Content"."content" FROM "Content"`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectJsonColumn_AsWhereLike(t *testing.T) {
	// Process actual value
	titleCol := query.JsonColumn("content.title", option.As("title"))
	actual := query.Select(titleCol).
		From(content).
		Where(query.Like(titleCol)).
		Build()

	// Compare
	expected := `SELECT ("Content"."content"->>'title') AS "title" FROM "Content" WHERE "Content"."content"->>'title' LIKE ?`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectJsonColumn_Join(t *testing.T) {
	// Process actual value
	actual := query.Select(query.JsonColumn("content.title", option.Schema(content), option.As("title"))).
		From(contentTag, option.As("ct")).
		Join(content,
			query.Equal(query.Column("contentId"), query.On("id")),
			option.As("c"),
		).
		Where(
			query.Equal(query.Column("tag", option.Schema(contentTag))),
		).
		Build()

	// Compare
	expected := `SELECT ("c"."content"->>'title') AS "title" FROM "ContentTag" AS "ct" INNER JOIN "Content" AS "c" ON "ct"."contentId" = "c"."id" WHERE "ct"."tag" = ?`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectJsonColumn_WhereNotNull(t *testing.T) {
	// Process actual value
	actual := query.Select(query.Column("*")).
		From(content).
		Where(query.IsNotNull(query.JsonColumn("content.title"))).
		Build()

	// Compare
	expected := `SELECT "Content"."createdAt", "Content"."updatedAt", "Content"."id", "Content"."content" FROM "Content" WHERE "Content"."content"->>'title' IS NOT NULL`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectJsonColumn_GetColumn(t *testing.T) {
	// Process actual value
	actual := query.JsonColumn("content.title").GetColumn()

	// Compare
	expected := `content`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectJsonColumn_InvalidColumn1(t *testing.T) {
	defer func() {
		errStr := "nsql: invalid JsonColumn value, attributes is not defined"
		r := recover()
		if r == nil {
			t.Errorf("FAILED\n  > code did not panic")
			return
		}

		err, ok := r.(error)
		if !ok {
			t.Errorf("%s: FAILED\n  > unknown recovered value: %v", errStr, r)
			return
		}

		if err.Error() != errStr {
			t.Errorf("%s: FAILED\n  > got different error: %v", errStr, err)
			return
		}
	}()
	query.JsonColumn("content.")
}

func TestSelectJsonColumn_InvalidColumn2(t *testing.T) {
	defer func() {
		errStr := "nsql: invalid JsonColumn value, attributes is not defined"
		r := recover()
		if r == nil {
			t.Errorf("FAILED\n  > code did not panic")
			return
		}

		err, ok := r.(error)
		if !ok {
			t.Errorf("%s: FAILED\n  > unknown recovered value: %v", errStr, r)
			return
		}

		if err.Error() != errStr {
			t.Errorf("%s: FAILED\n  > got different error: %v", errStr, err)
			return
		}
	}()
	query.JsonColumn("content")
}
