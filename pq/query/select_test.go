package query_test

import (
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"github.com/nbs-go/nsql/test_utils"
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

type Location struct {
	Id       int64  `db:"id"`
	Name     string `db:"name"`
	IsActive bool   `db:"isActive"`
}

type Route struct {
	Id            int64 `db:"id"`
	OriginId      int64 `db:"originId"`
	DestinationId int64 `db:"destinationId"`
}

var person = schema.New(schema.FromModelRef(Person{}))
var vehicleOwnership = schema.New(schema.FromModelRef(VehicleOwnership{}))
var vehicle = schema.New(schema.FromModelRef(Vehicle{}))

func TestSelectAll(t *testing.T) {
	// Select All
	testSelectBuilder(t, "SELECT ALL",
		query.Select(query.Column("*")).
			From(person),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`,
	)
}

func TestSelectWithAlias(t *testing.T) {
	personAs := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	testSelectBuilder(t, "SELECT WITH ALIAS TABLE",
		query.Select(query.Column("*")).
			From(personAs),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`,
	)
}

func TestSelectSpecifiedFields(t *testing.T) {
	testSelectBuilder(t, "SELECT SPECIFIED FIELDS",
		query.Select(query.Columns("id", "fullName", "gender", option.Schema(person))).
			From(person),
		`SELECT "Person"."id", "Person"."fullName" FROM "Person"`,
	)
}

func TestSelectSpecifiedFieldsWithoutDeclareSchema(t *testing.T) {
	testSelectBuilder(t, "SELECT SPECIFIED FIELDS (WITHOUT DECLARE SCHEMA)",
		query.Select(query.Columns("id", "fullName", "gender")).
			From(person),
		`SELECT "Person"."id", "Person"."fullName" FROM "Person"`,
	)
}

func TestSelectWithLimitedResult(t *testing.T) {
	testSelectBuilder(t, "SELECT WITH LIMITED RESULT",
		query.Select(query.Column("*", option.Schema(person))).
			From(person).
			Limit(10),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10`,
	)
}

func TestSelectWithSkippedResult(t *testing.T) {
	testSelectBuilder(t, "SELECT WITH SKIPPED RESULT",
		query.Select(query.Column("*")).
			From(person).
			Skip(1),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" OFFSET 1`,
	)
}

func TestSelectWithLimitAndSkippedResult(t *testing.T) {
	testSelectBuilder(t, "SELECT WITH LIMITED AND SKIPPED RESULT",
		query.Select(query.Column("*")).
			From(person).
			Limit(10).
			Skip(10),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10 OFFSET 10`,
	)
}

func TestSelectWithOrderBy(t *testing.T) {
	testSelectBuilder(t, "SELECT WITH ORDER BY",
		query.Select(query.Column("*")).
			From(person).
			OrderBy("createdAt"),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" ORDER BY "Person"."createdAt" ASC`,
	)
}

func TestSelectWithOrderByAndAliasTable(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	testSelectBuilder(t, "SELECT WITH ORDER BY WITH ALIAS TABLE",
		query.Select(query.Column("*")).
			From(pSchema).
			OrderBy("createdAt", option.SortDirection(op.Descending), option.Schema(pSchema)),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC`,
	)
}

func TestSelectWithOrderByLimitSkippedAndAliasTable(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	testSelectBuilder(t, "SELECT WITH ORDER BY, LIMIT AND SKIP",
		query.Select(query.Column("*")).
			From(pSchema).
			OrderBy("createdAt", option.SortDirection(op.Descending)).
			Limit(10).
			Skip(0),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC LIMIT 10 OFFSET 0`,
	)
}

func TestSelectWithOrderByUsingUndeclaredColumn(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	testSelectBuilder(t, "SELECT WITH ORDER BY USING UNDECLARED COLUMN",
		query.Select(query.Column("*")).
			From(pSchema).
			OrderBy("age", option.SortDirection(op.Descending)),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`,
	)
}

func TestSelectWithWhereByPrimaryKey(t *testing.T) {
	testSelectBuilder(t, "SELECT WITH WHERE BY PK",
		query.Select(query.Column("*")).
			From(person).
			Where(query.Equal(query.Column(person.PrimaryKey()))),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" WHERE "Person"."id" = ?`,
	)
}

func TestSelectWithWhereAnd(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	testSelectBuilder(t, "SELECT WITH WHERE AND",
		query.Select(query.Column("*")).
			From(pSchema).
			Where(query.Equal(query.Column(pSchema.PrimaryKey())), query.Equal(query.Column("fullName"))),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" WHERE "p"."id" = ? AND "p"."fullName" = ?`,
	)
}

func TestSelectWithWhereAndOrNested(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	testSelectBuilder(t, "SELECT WITH WHERE AND OR NESTED",
		query.Select(query.Column("*")).
			From(pSchema).
			Where(
				query.Or(
					query.Like(query.Column("fullName")),
					query.NotLike(query.Column("fullName")),
					query.ILike(query.Column("fullName")),
					query.NotILike(query.Column("fullName")),
					query.And(
						query.In(query.Column("id"), 1),
						query.NotIn(query.Column("id"), 2),
					),
				),
				query.And(
					query.Equal(query.Column("id")),
					query.NotEqual(query.Column("id")),
					query.LessThan(query.Column("id")),
					query.LessThanEqual(query.Column("id")),
					query.GreaterThan(query.Column("id")),
					query.GreaterThanEqual(query.Column("id")),
					query.Or(
						query.Between(query.Column("createdAt")),
						query.NotBetween(query.Column("updatedAt")),
					),
				),
			),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" WHERE ("p"."fullName" LIKE ? OR "p"."fullName" NOT LIKE ? OR "p"."fullName" ILIKE ? OR "p"."fullName" NOT ILIKE ? OR ("p"."id" IN (?) AND "p"."id" NOT IN (?, ?))) AND ("p"."id" = ? AND "p"."id" != ? AND "p"."id" < ? AND "p"."id" <= ? AND "p"."id" > ? AND "p"."id" >= ? AND ("p"."createdAt" BETWEEN ? AND ? OR "p"."updatedAt" NOT BETWEEN ? AND ?))`,
	)
}

func TestSelectCountAll(t *testing.T) {
	testSelectBuilder(t, "COUNT ALL",
		query.Select(query.Count("*")).
			From(person),
		`SELECT COUNT(*) FROM "Person"`,
	)
}

func TestSelectCountById(t *testing.T) {
	testSelectBuilder(t, "COUNT BY ID",
		query.Select(query.Count("id", option.Schema(person))).
			From(person),
		`SELECT COUNT("Person"."id") FROM "Person"`,
	)
}

func TestSelectCountByIdWithAliasField(t *testing.T) {
	testSelectBuilder(t, "COUNT BY ID WITH ALIAS FIELD",
		query.Select(query.Count("id", option.Schema(person), option.As("count"))).
			From(person),
		`SELECT COUNT("Person"."id") AS "count" FROM "Person"`,
	)
}

func TestSelectCountByIdWithAliasTable(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	testSelectBuilder(t, "COUNT BY ID WITH ALIAS TABLE",
		query.Select(query.Count("id", option.Schema(pSchema))).
			From(pSchema),
		`SELECT COUNT("p"."id") FROM "Person" AS "p"`,
	)
}

func TestPanicCount(t *testing.T) {
	defer test_utils.RecoverPanic(t, "NO COLUMN DECLARE ON COUNT", `column "age" is not declared in schema "Person"`)()
	query.Select(query.Count("age", option.Schema(person))).From(person)
}

func TestEmptyWhere(t *testing.T) {
	testSelectBuilder(t, "NIL WHERE",
		query.Select(query.Column("*")).
			From(person).Where(nil),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`,
	)
}

func TestIsExists(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))

	actual := query.Select(query.GreaterThan(query.Count("id"), query.IntVar(0), option.As("isExists"))).
		From(pSchema).Build()
	expected := `SELECT COUNT("p"."id") > 0 AS "isExists" FROM "Person" AS "p"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectJoin_ManyToManyWithTableAlias(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	voSchema := schema.New(schema.FromModelRef(VehicleOwnership{}), schema.As("vo"))
	vSchema := schema.New(schema.FromModelRef(Vehicle{}), schema.As("v"))

	// Assert
	actual := query.Select(
		query.Column("*"),
		query.Column("*", option.Schema(voSchema)),
		query.Column("*", option.Schema(vSchema)),
	).
		From(pSchema).
		Join(voSchema, query.Equal(query.Column("id"), query.On("personId"))).
		Join(vSchema, query.Equal(query.Column("vehicleId", option.Schema(voSchema)), query.On("id"))).
		Build()
	expected := `SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName", "vo"."createdAt" AS "vo.createdAt", "vo"."updatedAt" AS "vo.updatedAt", "vo"."id" AS "vo.id", "vo"."personId" AS "vo.personId", "vo"."vehicleId" AS "vo.vehicleId", "v"."createdAt" AS "v.createdAt", "v"."updatedAt" AS "v.updatedAt", "v"."id" AS "v.id", "v"."name" AS "v.name", "v"."category" AS "v.category" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" INNER JOIN "Vehicle" AS "v" ON "vo"."vehicleId" = "v"."id"`
	if actual != expected {
		t.Errorf("%s - FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestSelectInnerJoinTwoTable(t *testing.T) {
	testSelectBuilder(t, "INNER JOIN 2 TABLE",
		query.Select(query.Column("*"), query.Column("*", option.Schema(vehicleOwnership))).
			From(person).
			Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId")), option.JoinMethod(op.InnerJoin)),
		`SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" INNER JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`,
	)
}

func TestSelectLeftJoinTwoTable(t *testing.T) {
	testSelectBuilder(t, "LEFT JOIN 2 TABLE",
		query.Select(query.Column("*"), query.Column("*", option.Schema(vehicleOwnership))).
			From(person).
			Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId")), option.JoinMethod(op.LeftJoin)),
		`SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" LEFT JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`,
	)
}

func TestSelectRightJoinTwoTable(t *testing.T) {
	testSelectBuilder(t, "RIGHT JOIN 2 TABLE",
		query.Select(query.Column("*"), query.Column("*", option.Schema(vehicleOwnership))).
			From(person).
			Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId")), option.JoinMethod(op.RightJoin)),
		`SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" RIGHT JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`,
	)
}

func TestSelectFullJoinTwoTable(t *testing.T) {
	testSelectBuilder(t, "FULL JOIN 2 TABLE",
		query.Select(query.Column("*"), query.Column("*", option.Schema(vehicleOwnership))).
			From(person).
			Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId")), option.JoinMethod(op.FullJoin)),
		`SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" FULL OUTER JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`,
	)
}

func TestSelectManyToManyJoinTable(t *testing.T) {
	testSelectBuilder(t, "MANY TO MANY",
		query.Select(query.Column("*"), query.Column("*", option.Schema(vehicleOwnership)), query.Column("*", option.Schema(vehicle))).
			From(person).
			Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId"))).
			Join(vehicle, query.Equal(query.Column("vehicleId", option.Schema(vehicleOwnership)), query.On("id"))),
		`SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId", "Vehicle"."createdAt" AS "Vehicle.createdAt", "Vehicle"."updatedAt" AS "Vehicle.updatedAt", "Vehicle"."id" AS "Vehicle.id", "Vehicle"."name" AS "Vehicle.name", "Vehicle"."category" AS "Vehicle.category" FROM "Person" INNER JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId" INNER JOIN "Vehicle" ON "VehicleOwnership"."vehicleId" = "Vehicle"."id"`,
	)
}

func TestSelectInnerJoinTableWithFilter(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	voSchema := schema.New(schema.FromModelRef(VehicleOwnership{}), schema.As("vo"))
	testSelectBuilder(t, "INNER JOIN WITH FILTER",
		query.Select(query.Column("*")).
			From(pSchema).
			Join(voSchema,
				query.Equal(query.Column("id"), query.On("personId")),
				option.JoinMethod(op.InnerJoin),
			).
			Where(query.GreaterThanEqual(query.Column("createdAt", option.Schema(voSchema)))),
		`SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" WHERE "vo"."createdAt" >= ?`,
	)
}

func TestSelectJoinTableWithAndCondition(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	voSchema := schema.New(schema.FromModelRef(VehicleOwnership{}), schema.As("vo"))
	testSelectBuilder(t, "INNER JOIN WITH AND CONDITION",
		query.Select(query.Column("*"), query.Column("vehicleId", option.Schema(voSchema))).
			From(pSchema).
			Join(voSchema,
				query.And(
					query.Equal(query.Column("id"), query.On("personId")),
					query.GreaterThan(query.Column("createdAt", option.Schema(voSchema))),
				),
				option.JoinMethod(op.InnerJoin),
			),
		`SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName", "vo"."vehicleId" AS "vo.vehicleId" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" AND "vo"."createdAt" > ?`,
	)
}

func TestSelectEmptyWhere(t *testing.T) {
	testSelectBuilder(t, "Empty Where",
		query.Select(query.Column("*")).
			From(person).
			Where(query.And()),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`)
}

func TestSelectSingleWhere(t *testing.T) {
	testSelectBuilder(t, "Single Where",
		query.Select(query.Column("*")).
			From(person).
			Where(query.And(
				query.Equal(query.Column("id")),
			)),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" WHERE "Person"."id" = ?`)
}

func TestSelectEmptyWhereNested(t *testing.T) {
	testSelectBuilder(t, "Empty Where Nested",
		query.Select(query.Column("*")).
			From(person).
			Where(query.And(
				query.Equal(query.Column("id")),
				query.Or(),
			)),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" WHERE "Person"."id" = ?`)
}

func testSelectBuilder(t *testing.T, expectation string, b *query.SelectBuilder, expected string) {
	actual := b.Build()
	if actual != expected {
		t.Errorf("%s: FAILED\n  > got different generated query. Query = %s", expectation, actual)
	} else {
		t.Logf("%s: PASSED", expectation)
	}
}

func BenchmarkJoinManyToMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query.Select(
			query.Column("*"),
			query.Column("*", option.Schema(vehicleOwnership)),
			query.Column("*", option.Schema(vehicle)),
		).
			From(person, option.As("p")).
			Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId")), option.As("vo")).
			Join(vehicle, query.Equal(query.Column("vehicleId", option.Schema(vehicleOwnership)), query.On("id")), option.As("v")).
			Build()
	}
}

func TestFromConstructorSelectAll(t *testing.T) {
	b := query.From(person)
	test_utils.CompareString(t, "SELECT ALL",
		b.Select(query.Column("*")).Build(),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`,
	)
}

func TestFromConstructorReplaceSelectCountAll(t *testing.T) {
	b := query.From(person)
	test_utils.CompareString(t, "REPLACE SELECT COUNT ALL",
		b.Select(query.Count("*")).Build(),
		`SELECT COUNT(*) FROM "Person"`,
	)
}

func TestFromConstructorSelectAllWithAlias(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	b := query.From(pSchema)
	test_utils.CompareString(t, "SELECT ALL WITH ALIAS",
		b.Select(query.Column("*")).Build(),
		`SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`,
	)
}

func TestResetOrderBy(t *testing.T) {
	b := query.From(person).Select(query.Column("*"))

	test_utils.CompareString(t, "SELECT ALL",
		b.OrderBy("createdAt").Build(),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" ORDER BY "Person"."createdAt" ASC`,
	)

	b.ResetOrderBy()

	test_utils.CompareString(t, "REPLACE SELECT COUNT ALL",
		b.OrderBy("updatedAt", option.SortDirection(op.Descending)).Build(),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" ORDER BY "Person"."updatedAt" DESC`,
	)
}

func TestUndeclaredOrderByColumn(t *testing.T) {
	b := query.Select(query.Column("*")).From(person).OrderBy("gender", option.Schema(person))

	test_utils.CompareString(t, "NO ORDER BY",
		b.Build(),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`,
	)
}

func TestUndeclaredWhereColumn(t *testing.T) {
	b := query.Select(query.Column("*")).From(person).Where(query.Equal(query.Column("gender"), option.Schema(person)))

	test_utils.CompareString(t, "NO ORDER BY",
		b.Build(),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`,
	)
}

func TestReusableWhereCondition(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	b := query.Select(query.Column("*")).From(pSchema).Where(query.GreaterThan(query.Column("createdAt")))
	b.Build()

	q := b.Select(query.Count("*")).Build()

	test_utils.CompareString(t, "REUSABLE WHERE CONDITION",
		q,
		`SELECT COUNT(*) FROM "Person" AS "p" WHERE "p"."createdAt" > ?`,
	)
}

func TestResetLimitSkip(t *testing.T) {
	b := query.From(person).Select(query.Column("*")).Skip(0).Limit(10)

	test_utils.CompareString(t, "SELECT ALL WITH LIMIT AND SKIP",
		b.Build(),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10 OFFSET 0`,
	)

	b.ResetSkip().ResetLimit()

	test_utils.CompareString(t, "SELECT ALL",
		b.Build(),
		`SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`,
	)
}

func TestJoinBindVarCondition(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	voSchema := schema.New(schema.FromModelRef(VehicleOwnership{}), schema.As("vo"))
	testSelectBuilder(t, "INNER JOIN WITH BIND VAR CONDITION",
		query.Select(query.Column("*")).
			From(pSchema).
			Join(voSchema,
				query.And(
					query.Equal(query.Column("id"), query.On("personId")),
					query.Equal(query.Column("vehicleId", option.Schema(voSchema)), query.BindVar()),
				),
				option.JoinMethod(op.InnerJoin),
			).
			Where(query.GreaterThanEqual(query.Column("createdAt", option.Schema(voSchema)))),
		`SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" AND "vo"."vehicleId" = ? WHERE "vo"."createdAt" >= ?`,
	)
}

func TestJoinInvalidBindVarCondition(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	voSchema := schema.New(schema.FromModelRef(VehicleOwnership{}), schema.As("vo"))
	testSelectBuilder(t, "INNER JOIN WITH NO COLUMN BIND VAR CONDITION",
		query.Select(query.Column("*")).
			From(pSchema).
			Join(voSchema,
				query.And(
					query.Equal(query.Column("id"), query.On("personId")),
					query.Equal(query.Column("vehicleId"), query.BindVar()),
				),
				option.JoinMethod(op.InnerJoin),
			).
			Where(query.GreaterThanEqual(query.Column("createdAt", option.Schema(voSchema)))),
		`SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" WHERE "vo"."createdAt" >= ?`,
	)
}

func TestSelectAsField(t *testing.T) {
	// Process actual value
	actual := query.Select(query.Column("id", option.As("personId"))).
		From(person).
		Build()

	// Compare
	expected := `SELECT "Person"."id" AS "personId" FROM "Person"`
	if actual != expected {
		t.Errorf("%s: FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestIsNull(t *testing.T) {
	// Process actual value
	actual := query.Select(query.Column("id")).
		From(person).
		Where(query.IsNull(query.Column("fullName"))).
		Build()

	// Compare
	expected := `SELECT "Person"."id" FROM "Person" WHERE "Person"."fullName" IS NULL`
	if actual != expected {
		t.Errorf("%s: FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestIsNotNull(t *testing.T) {
	// Process actual value
	actual := query.Select(query.Column("id")).
		From(person).
		Where(query.IsNotNull(query.Column("fullName"))).
		Build()

	// Compare
	expected := `SELECT "Person"."id" FROM "Person" WHERE "Person"."fullName" IS NOT NULL`
	if actual != expected {
		t.Errorf("%s: FAILED\n  > got different generated query. Query = %s", expected, actual)
	}
}

func TestPrintOptionAsDeprecationWarning(t *testing.T) {
	// option.As in From query
	actual := query.Select(
		query.Column("*"),
	).
		From(person).
		Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId")), option.As("vo")).
		Build()
	expected := `SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName" FROM "Person" INNER JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > Unexpected actual value.\n  > Actual = %s", expected, actual)
	}
}

func TestSameTableDifferentSchema(t *testing.T) {
	origin := schema.New(schema.FromModelRef(Location{}), schema.As("o"))
	dest := schema.New(schema.FromModelRef(Location{}), schema.As("d"))
	route := schema.New(schema.FromModelRef(Route{}), schema.As("r"))

	// Assert
	actual := query.Select(
		query.Column("*"),
		query.Column("*", option.Schema(origin)),
		query.Column("*", option.Schema(dest)),
	).
		From(route).
		Join(origin, query.Equal(query.Column("originId"), query.On("id"))).
		Join(dest, query.Equal(query.Column("destinationId"), query.On("id"))).
		Where(
			query.Equal(query.Column("originId")),
			query.Equal(query.Column("isActive", option.Schema(origin))),
			query.Equal(query.Column("isActive", option.Schema(dest))),
		).
		Build()
	expected := `SELECT "r"."id" AS "r.id", "r"."originId" AS "r.originId", "r"."destinationId" AS "r.destinationId", "o"."id" AS "o.id", "o"."name" AS "o.name", "o"."isActive" AS "o.isActive", "d"."id" AS "d.id", "d"."name" AS "d.name", "d"."isActive" AS "d.isActive" FROM "Route" AS "r" INNER JOIN "Location" AS "o" ON "r"."originId" = "o"."id" INNER JOIN "Location" AS "d" ON "r"."destinationId" = "d"."id" WHERE "r"."originId" = ? AND "o"."isActive" = ? AND "d"."isActive" = ?`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}
