package query

import (
	"github.com/nbs-go/nsql/query/op"
	opt "github.com/nbs-go/nsql/query/option"
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

func TestSelectBasicQuery(t *testing.T) {
	// Select All
	testSelectBuilder(t, "SELECT ALL",
		Select(opt.Columns("*")).
			From(person),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`,
	)

	testSelectBuilder(t, "SELECT SPECIFIED FIELDS",
		Select(opt.Columns("id", "fullName", "gender")).
			From(person),
		`SELECT "Person"."id", "Person"."fullName" FROM "Person"`,
	)

	testSelectBuilder(t, "SELECT WITH ALIAS TABLE",
		Select(opt.Columns("*")).
			From(person, opt.As("p")),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`,
	)

	testSelectBuilder(t, "SELECT WITH LIMITED RESULT",
		Select(opt.Columns(person, "*")).
			From(person).
			Limit(10),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10`,
	)

	testSelectBuilder(t, "SELECT WITH SKIPPED RESULT",
		Select(opt.Columns("*")).
			From(person).
			Skip(1),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" OFFSET 1`,
	)

	testSelectBuilder(t, "SELECT WITH LIMITED AND SKIPPED RESULT",
		Select(opt.Columns("*")).
			From(person).
			Limit(10).
			Skip(10),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10 OFFSET 10`,
	)

	testSelectBuilder(t, "SELECT WITH ORDER BY",
		Select(opt.Columns("*")).
			From(person).
			OrderBy("createdAt"),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" ORDER BY "Person"."createdAt" ASC`,
	)

	testSelectBuilder(t, "SELECT WITH ORDER BY WITH ALIAS TABLE",
		Select(opt.Columns("*")).
			From(person, opt.As("p")).
			OrderBy("createdAt", opt.SortDirection(op.Descending), opt.Schema(person)),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC`,
	)

	testSelectBuilder(t, "SELECT WITH ORDER BY, LIMIT AND SKIP",
		Select(opt.Columns("*")).
			From(person, opt.As("p")).
			OrderBy("createdAt", opt.SortDirection(op.Descending)).
			Limit(10).
			Skip(0),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC LIMIT 10 OFFSET 0`,
	)

	testSelectBuilder(t, "SELECT WITH ORDER BY USING UNDECLARED COLUMN",
		Select(opt.Columns("*")).
			From(person, opt.As("p")).
			OrderBy("age", opt.SortDirection(op.Descending)),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`,
	)

	testSelectBuilder(t, "SELECT WITH WHERE BY PK",
		Select(opt.Columns("*")).
			From(person).
			Where(Equal(person, person.PrimaryKey)),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" WHERE "Person"."id" = ?`,
	)

	testSelectBuilder(t, "SELECT WITH WHERE AND",
		Select(opt.Columns("*")).
			From(person, opt.As("p")).
			Where(Equal(person, person.PrimaryKey), Equal(person, "fullName")),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" WHERE "p"."id" = ? AND "p"."fullName" = ?`,
	)

	testSelectBuilder(t, "SELECT WITH WHERE AND OR NESTED",
		Select(opt.Columns("*")).
			From(person, opt.As("p")).
			Where(
				Or(
					Like(person, "fullName"),
					NotLike(person, "fullName"),
					ILike(person, "fullName"),
					NotILike(person, "fullName"),
					And(
						In(person, "id"),
						NotIn(person, "id"),
					),
				),
				And(
					Equal(person, "id"),
					NotEqual(person, "id"),
					LessThan(person, "id"),
					LessThanEqual(person, "id"),
					GreaterThan(person, "id"),
					GreaterThanEqual(person, "id"),
					Or(
						Between(person, "createdAt"),
						NotBetween(person, "updatedAt"),
					),
				),
			),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" WHERE ("p"."fullName" LIKE ? OR "p"."fullName" NOT LIKE ? OR "p"."fullName" ILIKE ? OR "p"."fullName" NOT ILIKE ? OR ("p"."id" IN (?) AND "p"."id" NOT IN (?))) AND ("p"."id" = ? AND "p"."id" != ? AND "p"."id" < ? AND "p"."id" <= ? AND "p"."id" > ? AND "p"."id" >= ? AND ("p"."createdAt" BETWEEN ? AND ? OR "p"."updatedAt" NOT BETWEEN ? AND ?))`,
	)
}

func TestSelectCount(t *testing.T) {
	testSelectBuilder(t, "COUNT ALL",
		Select(opt.Count("*")).
			From(person),
		`SELECT COUNT(*) FROM "Person"`,
	)

	testSelectBuilder(t, "COUNT BY ID",
		Select(opt.Count("id", opt.Schema(person))).
			From(person),
		`SELECT COUNT("Person"."id") FROM "Person"`,
	)

	testSelectBuilder(t, "COUNT BY ID WITH ALIAS TABLE",
		Select(opt.Count("id", opt.Schema(person))).
			From(person, opt.As("p")),
		`SELECT COUNT("p"."id") FROM "Person" AS "p"`,
	)

	testSelectBuilder(t, "COUNT BY ID WITH ALIAS FIELD",
		Select(opt.Count("id", opt.Schema(person), opt.As("count"))).
			From(person),
		`SELECT COUNT("Person"."id") AS "count" FROM "Person"`,
	)
}

func testSelectBuilder(t *testing.T, name string, b *SelectBuilder, expected string) {
	actual := b.Build()
	if actual != expected {
		t.Errorf("%s: FAILED\n  > got different generated query. Query = %s", name, actual)
	} else {
		t.Logf("%s: PASSED", name)
	}
}
