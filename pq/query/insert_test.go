package query

import (
	"github.com/nbs-go/nsql/schema"
	"github.com/nbs-go/nsql/test_utils"
	"testing"
)

func TestInsert(t *testing.T) {
	// Init schema
	s := schema.New(schema.FromModelRef(new(Person)))

	// Test #1
	test_utils.CompareString(t, "INSERT ALL COLUMNS",
		Insert(s, "*").Build(),
		`INSERT INTO "Person"("createdAt", "updatedAt", "fullName") VALUES (:createdAt, :updatedAt, :fullName)`)

	// Test #2
	test_utils.CompareString(t, "INSERT SPECIFIED COLUMNS",
		Insert(s, "id", "fullName", "age").Build(),
		`INSERT INTO "Person"("id", "fullName") VALUES (:id, :fullName)`,
	)

	// Init schema without auto increment
	s = schema.New(schema.FromModelRef(new(Person)), schema.AutoIncrement(false))

	// Test #3
	test_utils.CompareString(t, "INSERT ALL (NO AUTO INCREMENT)",
		Insert(s, "*").Build(),
		`INSERT INTO "Person"("createdAt", "updatedAt", "id", "fullName") VALUES (:createdAt, :updatedAt, :id, :fullName)`,
	)
}
