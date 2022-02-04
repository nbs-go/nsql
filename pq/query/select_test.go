package query

import (
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/schema"
	"testing"
	"time"
)

type Person struct {
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
	Id        int64     `db:"id"`
	FullName  string    `db:"fullName"`
}

var person = schema.New(schema.FromModelRef(Person{}))

func TestSelectBuilder(t *testing.T) {
	// Select All
	testSelectBuilder(t, "Select All",
		Select().
			Columns(person, "*").
			From(person),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`,
	)

	testSelectBuilder(t, "Select Specified Fields",
		Select().
			Columns(person, "id", "fullName", "gender").
			From(person),
		`SELECT "Person"."id", "Person"."fullName" FROM "Person"`,
	)

	testSelectBuilder(t, "Select with Alias Table",
		Select().
			Columns(person, "*").
			From(person, "p"),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`,
	)

	testSelectBuilder(t, "Select with Limited Result",
		Select().
			Columns(person, "*").
			From(person).
			Limit(10),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10`,
	)

	testSelectBuilder(t, "Select with Skipped Result",
		Select().
			Columns(person, "*").
			From(person).
			Skip(1),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" OFFSET 1`,
	)

	testSelectBuilder(t, "Select with Limited and Skipped Result",
		Select().
			Columns(person, "*").
			From(person).
			Limit(10).
			Skip(10),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10 OFFSET 10`,
	)

	testSelectBuilder(t, "Select with Order By",
		Select().
			Columns(person, "*").
			From(person).
			OrderBy(person, "createdAt"),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" ORDER BY "Person"."createdAt"`,
	)

	testSelectBuilder(t, "Select with Order By Aliased",
		Select().
			Columns(person, "*").
			From(person, "p").
			OrderBy(person, "createdAt", query.Descending),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC`,
	)

	testSelectBuilder(t, "Select with Order By, Limit and Skip",
		Select().
			Columns(person, "*").
			From(person, "p").
			OrderBy(person, "createdAt", query.Descending).
			Limit(10).
			Skip(0),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC LIMIT 10 OFFSET 0`,
	)

	testSelectBuilder(t, "Select with Order By using Undeclared column",
		Select().
			Columns(person, "*").
			From(person, "p").
			OrderBy(person, "age", query.Descending),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`,
	)
}

func testSelectBuilder(t *testing.T, name string, b *SelectBuilder, expected string) {
	actual := b.Build()
	if actual != expected {
		t.Errorf("got different generated %s query. Query = %s", name, actual)
	} else {
		t.Logf("%s: PASSED", name)
	}
}
