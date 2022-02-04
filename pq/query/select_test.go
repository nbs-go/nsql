package query

import (
	"github.com/nbs-go/nsql/pq/op"
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
			OrderBy(person, "createdAt", op.Descending),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC`,
	)

	testSelectBuilder(t, "Select with Order By, Limit and Skip",
		Select().
			Columns(person, "*").
			From(person, "p").
			OrderBy(person, "createdAt", op.Descending).
			Limit(10).
			Skip(0),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC LIMIT 10 OFFSET 0`,
	)

	testSelectBuilder(t, "Select with Order By using Undeclared column",
		Select().
			Columns(person, "*").
			From(person, "p").
			OrderBy(person, "age", op.Descending),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`,
	)

	testSelectBuilder(t, "Select with Where by PK",
		Select().
			Columns(person, "*").
			From(person).
			Where(Equal(person, person.PrimaryKey)),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" WHERE "Person"."id" = ?`,
	)

	testSelectBuilder(t, "Select with Where And",
		Select().
			Columns(person, "*").
			From(person, "p").
			Where(Equal(person, person.PrimaryKey), Equal(person, "fullName")),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" WHERE "p"."id" = ? AND "p"."fullName" = ?`,
	)

	testSelectBuilder(t, "Select with Where And Or Nested",
		Select().
			Columns(person, "*").
			From(person, "p").
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

func testSelectBuilder(t *testing.T, name string, b *SelectBuilder, expected string) {
	actual := b.Build()
	if actual != expected {
		t.Errorf("got different generated %s query. Query = %s", name, actual)
	} else {
		t.Logf("%s: PASSED", name)
	}
}
