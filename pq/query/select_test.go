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
	actual := query.Select(query.Column("*")).
		From(person).Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectWithAlias(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))

	actual := query.Select(query.Column("*")).From(pSchema).Build()

	expected := `SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectSpecifiedFields(t *testing.T) {
	actual := query.Select(query.Columns("id", "fullName", "gender", option.Schema(person))).
		From(person).Build()

	expected := `SELECT "Person"."id", "Person"."fullName" FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectSpecifiedFieldsWithoutDeclareSchema(t *testing.T) {
	actual := query.Select(query.Columns("id", "fullName", "gender")).
		From(person).Build()

	expected := `SELECT "Person"."id", "Person"."fullName" FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectWithLimitedResult(t *testing.T) {
	actual := query.Select(query.Column("*", option.Schema(person))).
		From(person).
		Limit(10).
		Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectWithSkippedResult(t *testing.T) {
	actual := query.Select(query.Column("*")).
		From(person).
		Skip(1).
		Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" OFFSET 1`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectWithLimitAndSkippedResult(t *testing.T) {
	actual := query.Select(query.Column("*")).
		From(person).
		Limit(10).
		Skip(10).
		Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10 OFFSET 10`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectWithOrderBy(t *testing.T) {
	actual := query.Select(query.Column("*")).
		From(person).
		OrderBy("createdAt").
		Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" ORDER BY "Person"."createdAt" ASC`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectWithOrderByAndAliasTable(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))

	actual := query.Select(query.Column("*")).
		From(pSchema).
		OrderBy("createdAt", option.SortDirection(op.Descending), option.Schema(pSchema)).
		Build()

	expected := `SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectWithOrderByLimitSkippedAndAliasTable(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))

	actual := query.Select(query.Column("*")).
		From(pSchema).
		OrderBy("createdAt", option.SortDirection(op.Descending)).
		Limit(10).
		Skip(0).
		Build()

	expected := `SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" ORDER BY "p"."createdAt" DESC LIMIT 10 OFFSET 0`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectWithOrderByUsingUndeclaredColumn(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))

	actual := query.Select(query.Column("*")).
		From(pSchema).
		OrderBy("age", option.SortDirection(op.Descending)).
		Build()

	expected := `SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectWithWhereByPrimaryKey(t *testing.T) {
	actual := query.Select(query.Column("*")).
		From(person).
		Where(query.Equal(query.Column(person.PrimaryKey()))).
		Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" WHERE "Person"."id" = ?`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectWithWhereAnd(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))

	actual := query.Select(query.Column("*")).
		From(pSchema).
		Where(query.Equal(query.Column(pSchema.PrimaryKey())), query.Equal(query.Column("fullName"))).
		Build()

	expected := `SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" WHERE "p"."id" = ? AND "p"."fullName" = ?`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectWithWhereAndOrNested(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))

	actual := query.Select(query.Column("*")).
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
		).
		Build()

	expected := `SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p" WHERE ("p"."fullName" LIKE ? OR "p"."fullName" NOT LIKE ? OR "p"."fullName" ILIKE ? OR "p"."fullName" NOT ILIKE ? OR ("p"."id" IN (?) AND "p"."id" NOT IN (?, ?))) AND ("p"."id" = ? AND "p"."id" != ? AND "p"."id" < ? AND "p"."id" <= ? AND "p"."id" > ? AND "p"."id" >= ? AND ("p"."createdAt" BETWEEN ? AND ? OR "p"."updatedAt" NOT BETWEEN ? AND ?))`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectCountAll(t *testing.T) {
	actual := query.Select(query.Count("*")).
		From(person).
		Build()

	expected := `SELECT COUNT(*) FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectCountById(t *testing.T) {
	actual := query.Select(query.Count("id", option.Schema(person))).
		From(person).
		Build()

	expected := `SELECT COUNT("Person"."id") FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectCountByIdWithAliasField(t *testing.T) {
	actual := query.Select(query.Count("id", option.Schema(person), option.As("count"))).
		From(person).
		Build()

	expected := `SELECT COUNT("Person"."id") AS "count" FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectCountByIdWithAliasTable(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))

	actual := query.Select(query.Count("id", option.Schema(pSchema))).
		From(pSchema).
		Build()

	expected := `SELECT COUNT("p"."id") FROM "Person" AS "p"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestPanicCount(t *testing.T) {
	defer test_utils.RecoverPanic(t, "NO COLUMN DECLARE ON COUNT", `column "age" is not declared in schema "Person"`)()
	query.Select(query.Count("age", option.Schema(person))).From(person)
}

func TestEmptyWhere(t *testing.T) {
	actual := query.Select(query.Column("*")).
		From(person).Where(nil).
		Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
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
	actual := query.Select(query.Column("*"), query.Column("*", option.Schema(vehicleOwnership))).
		From(person).
		Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId")), option.JoinMethod(op.InnerJoin)).
		Build()

	expected := `SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" INNER JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectLeftJoinTwoTable(t *testing.T) {
	actual := query.Select(query.Column("*"), query.Column("*", option.Schema(vehicleOwnership))).
		From(person).
		Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId")), option.JoinMethod(op.LeftJoin)).
		Build()

	expected := `SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" LEFT JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectRightJoinTwoTable(t *testing.T) {
	actual := query.Select(query.Column("*"), query.Column("*", option.Schema(vehicleOwnership))).
		From(person).
		Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId")), option.JoinMethod(op.RightJoin)).
		Build()
	expected := `SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" RIGHT JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectFullJoinTwoTable(t *testing.T) {
	actual := query.Select(query.Column("*"), query.Column("*", option.Schema(vehicleOwnership))).
		From(person).
		Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId")), option.JoinMethod(op.FullJoin)).
		Build()

	expected := `SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId" FROM "Person" FULL OUTER JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectManyToManyJoinTable(t *testing.T) {
	actual := query.Select(query.Column("*"), query.Column("*", option.Schema(vehicleOwnership)), query.Column("*", option.Schema(vehicle))).
		From(person).
		Join(vehicleOwnership, query.Equal(query.Column("id"), query.On("personId"))).
		Join(vehicle, query.Equal(query.Column("vehicleId", option.Schema(vehicleOwnership)), query.On("id"))).
		Build()

	expected := `SELECT "Person"."createdAt" AS "Person.createdAt", "Person"."updatedAt" AS "Person.updatedAt", "Person"."id" AS "Person.id", "Person"."fullName" AS "Person.fullName", "VehicleOwnership"."createdAt" AS "VehicleOwnership.createdAt", "VehicleOwnership"."updatedAt" AS "VehicleOwnership.updatedAt", "VehicleOwnership"."id" AS "VehicleOwnership.id", "VehicleOwnership"."personId" AS "VehicleOwnership.personId", "VehicleOwnership"."vehicleId" AS "VehicleOwnership.vehicleId", "Vehicle"."createdAt" AS "Vehicle.createdAt", "Vehicle"."updatedAt" AS "Vehicle.updatedAt", "Vehicle"."id" AS "Vehicle.id", "Vehicle"."name" AS "Vehicle.name", "Vehicle"."category" AS "Vehicle.category" FROM "Person" INNER JOIN "VehicleOwnership" ON "Person"."id" = "VehicleOwnership"."personId" INNER JOIN "Vehicle" ON "VehicleOwnership"."vehicleId" = "Vehicle"."id"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectInnerJoinTableWithFilter(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	voSchema := schema.New(schema.FromModelRef(VehicleOwnership{}), schema.As("vo"))

	actual := query.Select(query.Column("*")).
		From(pSchema).
		Join(voSchema,
			query.Equal(query.Column("id"), query.On("personId")),
			option.JoinMethod(op.InnerJoin),
		).
		Where(query.GreaterThanEqual(query.Column("createdAt", option.Schema(voSchema)))).
		Build()

	expected := `SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" WHERE "vo"."createdAt" >= ?`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectJoinTableWithAndCondition(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	voSchema := schema.New(schema.FromModelRef(VehicleOwnership{}), schema.As("vo"))

	actual := query.Select(query.Column("*"), query.Column("vehicleId", option.Schema(voSchema))).
		From(pSchema).
		Join(voSchema,
			query.And(
				query.Equal(query.Column("id"), query.On("personId")),
				query.GreaterThan(query.Column("createdAt", option.Schema(voSchema))),
			),
			option.JoinMethod(op.InnerJoin),
		).Build()

	expected := `SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName", "vo"."vehicleId" AS "vo.vehicleId" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" AND "vo"."createdAt" > ?`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectEmptyWhere(t *testing.T) {
	actual := query.Select(query.Column("*")).
		From(person).
		Where(query.And()).
		Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectSingleWhere(t *testing.T) {
	actual := query.Select(query.Column("*")).
		From(person).
		Where(query.And(
			query.Equal(query.Column("id")),
		)).
		Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" WHERE "Person"."id" = ?`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestSelectEmptyWhereNested(t *testing.T) {
	actual := query.Select(query.Column("*")).
		From(person).
		Where(query.And(
			query.Equal(query.Column("id")),
			query.Or(),
		)).
		Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" WHERE "Person"."id" = ?`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
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

	actual := b.Select(query.Column("*")).Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestFromConstructorReplaceSelectCountAll(t *testing.T) {
	b := query.From(person)

	actual := b.Select(query.Count("*")).Build()

	expected := `SELECT COUNT(*) FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestFromConstructorSelectAllWithAlias(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	b := query.From(pSchema)

	actual := b.Select(query.Column("*")).Build()

	expected := `SELECT "p"."createdAt", "p"."updatedAt", "p"."id", "p"."fullName" FROM "Person" AS "p"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestResetOrderBy(t *testing.T) {
	b := query.From(person).Select(query.Column("*"))

	// generate first query select all
	actual := b.OrderBy("createdAt").Build()
	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" ORDER BY "Person"."createdAt" ASC`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}

	// reset order by
	b.ResetOrderBy()

	// set new order by
	actual = b.OrderBy("updatedAt", option.SortDirection(op.Descending)).Build()
	expected = `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" ORDER BY "Person"."updatedAt" DESC`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestUndeclaredOrderByColumn(t *testing.T) {
	actual := query.Select(query.Column("*")).
		From(person).
		OrderBy("gender", option.Schema(person)).
		Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestUndeclaredWhereColumn(t *testing.T) {
	actual := query.Select(query.Column("*")).
		From(person).
		Where(query.Equal(query.Column("gender"), option.Schema(person))).
		Build()

	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestReusableWhereCondition(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	b := query.Select(query.Column("*")).From(pSchema).Where(query.GreaterThan(query.Column("createdAt")))

	b.Build()

	actual := b.Select(query.Count("*")).Build()

	expected := `SELECT COUNT(*) FROM "Person" AS "p" WHERE "p"."createdAt" > ?`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestResetLimitSkip(t *testing.T) {
	b := query.From(person).Select(query.Column("*")).Skip(0).Limit(10)

	// first time build query
	actual := b.Build()
	expected := `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person" LIMIT 10 OFFSET 0`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}

	// reset skip, limit and build query
	actual = b.ResetSkip().ResetLimit().Build()

	expected = `SELECT "Person"."createdAt", "Person"."updatedAt", "Person"."id", "Person"."fullName" FROM "Person"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestJoinBindVarCondition(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	voSchema := schema.New(schema.FromModelRef(VehicleOwnership{}), schema.As("vo"))

	actual := query.Select(query.Column("*")).
		From(pSchema).
		Join(voSchema,
			query.And(
				query.Equal(query.Column("id"), query.On("personId")),
				query.Equal(query.Column("vehicleId", option.Schema(voSchema)), query.BindVar()),
			),
			option.JoinMethod(op.InnerJoin),
		).
		Where(query.GreaterThanEqual(query.Column("createdAt", option.Schema(voSchema)))).
		Build()

	expected := `SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" AND "vo"."vehicleId" = ? WHERE "vo"."createdAt" >= ?`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
}

func TestJoinInvalidBindVarCondition(t *testing.T) {
	pSchema := schema.New(schema.FromModelRef(Person{}), schema.As("p"))
	voSchema := schema.New(schema.FromModelRef(VehicleOwnership{}), schema.As("vo"))

	actual := query.Select(query.Column("*")).
		From(pSchema).
		Join(voSchema,
			query.And(
				query.Equal(query.Column("id"), query.On("personId")),
				query.Equal(query.Column("vehicleId"), query.BindVar()), // invalid condition vehicleId without comparator
			),
			option.JoinMethod(op.InnerJoin),
		).
		Where(query.GreaterThanEqual(query.Column("createdAt", option.Schema(voSchema)))).
		Build()

	expected := `SELECT "p"."createdAt" AS "p.createdAt", "p"."updatedAt" AS "p.updatedAt", "p"."id" AS "p.id", "p"."fullName" AS "p.fullName" FROM "Person" AS "p" INNER JOIN "VehicleOwnership" AS "vo" ON "p"."id" = "vo"."personId" WHERE "vo"."createdAt" >= ?`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different generated query. Actual = %s", expected, actual)
	}
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
