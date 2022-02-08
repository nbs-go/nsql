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

type Vehicle struct {
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	Category  string    `db:"category"`
}

type VehicleOwnership struct {
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
	Id        int64     `db:"id"`
	PersonId  int64     `db:"personId"`
	VehicleId int64     `db:"vehicleId"`
}

var person = schema.New(schema.FromModelRef(Person{}))
var vehicleOwnership = schema.New(schema.FromModelRef(VehicleOwnership{}))
var vehicle = schema.New(schema.FromModelRef(Vehicle{}))

func TestSelectBasicQuery(t *testing.T) {
	// Select All
	testSelectBuilder(t, "SELECT ALL",
		Select(Columns("*")).
			From(person),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`,
	)

	testSelectBuilder(t, "SELECT SPECIFIED FIELDS",
		Select(Columns("id", "fullName", "gender")).
			From(person),
		`SELECT "Person"."id", "Person"."fullName" FROM "Person"`,
	)

	testSelectBuilder(t, "SELECT WITH ALIAS TABLE",
		Select(Columns("*")).
			From(person, opt.As("p")),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`,
	)

	testSelectBuilder(t, "SELECT WITH LIMITED RESULT",
		Select(Columns(person, "*")).
			From(person).
			Limit(10),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10`,
	)

	testSelectBuilder(t, "SELECT WITH SKIPPED RESULT",
		Select(Columns("*")).
			From(person).
			Skip(1),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" OFFSET 1`,
	)

	testSelectBuilder(t, "SELECT WITH LIMITED AND SKIPPED RESULT",
		Select(Columns("*")).
			From(person).
			Limit(10).
			Skip(10),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10 OFFSET 10`,
	)

	testSelectBuilder(t, "SELECT WITH ORDER BY",
		Select(Columns("*")).
			From(person).
			OrderBy("createdAt"),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" ORDER BY "Person"."createdAt" ASC`,
	)

	testSelectBuilder(t, "SELECT WITH ORDER BY WITH ALIAS TABLE",
		Select(Columns("*")).
			From(person, opt.As("p")).
			OrderBy("createdAt", opt.SortDirection(op.Descending), opt.Schema(person)),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC`,
	)

	testSelectBuilder(t, "SELECT WITH ORDER BY, LIMIT AND SKIP",
		Select(Columns("*")).
			From(person, opt.As("p")).
			OrderBy("createdAt", opt.SortDirection(op.Descending)).
			Limit(10).
			Skip(0),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC LIMIT 10 OFFSET 0`,
	)

	testSelectBuilder(t, "SELECT WITH ORDER BY USING UNDECLARED COLUMN",
		Select(Columns("*")).
			From(person, opt.As("p")).
			OrderBy("age", opt.SortDirection(op.Descending)),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`,
	)

	testSelectBuilder(t, "SELECT WITH WHERE BY PK",
		Select(Columns("*")).
			From(person).
			Where(Equal(Column(person.PrimaryKey()))),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" WHERE "Person"."id" = ?`,
	)

	testSelectBuilder(t, "SELECT WITH WHERE AND",
		Select(Columns("*")).
			From(person, opt.As("p")).
			Where(Equal(Column(person.PrimaryKey())), Equal(Column("fullName"))),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" WHERE "p"."id" = ? AND "p"."fullName" = ?`,
	)

	testSelectBuilder(t, "SELECT WITH WHERE AND OR NESTED",
		Select(Columns("*")).
			From(person, opt.As("p")).
			Where(
				Or(
					Like(Column("fullName")),
					NotLike(Column("fullName")),
					ILike(Column("fullName")),
					NotILike(Column("fullName")),
					And(
						In(Column("id")),
						NotIn(Column("id")),
					),
				),
				And(
					Equal(Column("id")),
					NotEqual(Column("id")),
					LessThan(Column("id")),
					LessThanEqual(Column("id")),
					GreaterThan(Column("id")),
					GreaterThanEqual(Column("id")),
					Or(
						Between(Column("createdAt")),
						NotBetween(Column("updatedAt")),
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

//func TestIsExists(t *testing.T) {
//	testSelectBuilder(t, "IS EXISTS",
//		Select().
//			From(person).Where(Like(person, "fullName")),
//		`SELECT COUNT("Person"."id") > 0 AS "isExists" FROM "Person" WHERE "Person"."fullName" LIKE ?`,
//	)
//}

func TestSelectJoin(t *testing.T) {
	testSelectBuilder(t, "INNER JOIN 2 TABLE",
		Select(Columns("*")).
			Select(Columns("*", opt.Schema(vehicleOwnership))).
			From(person).
			Join(vehicleOwnership, Equal(Column("id"), On("personId")), opt.JoinMethod(op.InnerJoin)),
		`SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" INNER JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`,
	)

	testSelectBuilder(t, "LEFT JOIN 2 TABLE",
		Select(Columns("*")).
			Select(Columns("*", opt.Schema(vehicleOwnership))).
			From(person).
			Join(vehicleOwnership, Equal(Column("id"), On("personId")), opt.JoinMethod(op.LeftJoin)),
		`SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" LEFT JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`,
	)

	testSelectBuilder(t, "RIGHT JOIN 2 TABLE",
		Select(Columns("*")).
			Select(Columns("*", opt.Schema(vehicleOwnership))).
			From(person).
			Join(vehicleOwnership, Equal(Column("id"), On("personId")), opt.JoinMethod(op.RightJoin)),
		`SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" RIGHT JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`,
	)

	testSelectBuilder(t, "FULL JOIN 2 TABLE",
		Select(Columns("*")).
			Select(Columns("*", opt.Schema(vehicleOwnership))).
			From(person).
			Join(vehicleOwnership, Equal(Column("id"), On("personId")), opt.JoinMethod(op.FullJoin)),
		`SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" FULL OUTER JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`,
	)

	testSelectBuilder(t, "MANY TO MANY",
		Select(Columns("*")).
			Select(Columns("*", opt.Schema(vehicleOwnership))).
			Select(Columns("*", opt.Schema(vehicle))).
			From(person).
			Join(vehicleOwnership, Equal(Column("id"), On("personId"))).
			Join(vehicle, Equal(Column("vehicleId", opt.Schema(vehicleOwnership)), On("id"))),
		`SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId", "Vehicle"."createdAt" AS "Vehicle.createdAt", "Vehicle"."updatedAt" AS "Vehicle.updatedAt", "Vehicle"."id" AS "Vehicle.id", "Vehicle"."name" AS "Vehicle.name", "Vehicle"."category" AS "Vehicle.category" FROM "Person" INNER JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId" INNER JOIN "Vehicle" ON "VehicleOwnership"."vehicleId" = "Vehicle"."id"`,
	)

	testSelectBuilder(t, "MANY TO MANY WITH ALIAS",
		Select(Columns("*")).
			Select(Columns("*", opt.Schema(vehicleOwnership))).
			Select(Columns("*", opt.Schema(vehicle))).
			From(person, opt.As("p")).
			Join(vehicleOwnership, Equal(Column("id"), On("personId")), opt.As("vo")).
			Join(vehicle, Equal(Column("vehicleId", opt.Schema(vehicleOwnership)), On("id")), opt.As("v")),
		`SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName", "vo"."createdAt" AS "vo.createdAt", "vo"."updatedAt" AS "vo.updatedAt", "vo"."id" AS "vo.id", "vo"."personId" AS "vo.personId", "vo"."vehicleId" AS "vo.vehicleId", "v"."createdAt" AS "v.createdAt", "v"."updatedAt" AS "v.updatedAt", "v"."id" AS "v.id", "v"."name" AS "v.name", "v"."category" AS "v.category" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" INNER JOIN "Vehicle" AS "v" ON "vo"."vehicleId" = "v"."id"`,
	)

	testSelectBuilder(t, "INNER JOIN WITH FILTER",
		Select(Columns("*")).
			From(person, opt.As("p")).
			Join(vehicleOwnership,
				Equal(Column("id"), On("personId")),
				opt.JoinMethod(op.InnerJoin), opt.As("vo"),
			).
			Where(GreaterThanEqual(Column("createdAt", opt.Schema(vehicleOwnership)))),
		`SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" WHERE "vo"."createdAt" >= ?`,
	)

	testSelectBuilder(t, "INNER JOIN WITH AND CONDITION",
		Select(Columns("*")).
			Select(Columns(opt.Schema(vehicleOwnership), "vehicleId")).
			From(person, opt.As("p")).
			Join(vehicleOwnership,
				And(
					Equal(Column("id"), On("personId")),
					GreaterThan(Column("createdAt", opt.Schema(vehicleOwnership))),
				),
				opt.JoinMethod(op.InnerJoin), opt.As("vo"),
			),
		`SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName", "vo"."vehicleId" AS "vo.vehicleId" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" AND "vo"."createdAt" > ?`,
	)
}

func testSelectBuilder(t *testing.T, expectation string, b *SelectBuilder, expected string) {
	actual := b.Build()
	if actual != expected {
		t.Errorf("%s: FAILED\n  > got different generated query. Query = %s", expectation, actual)
	} else {
		t.Logf("%s: PASSED", expectation)
	}
}
