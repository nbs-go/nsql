package query_test

import (
	"github.com/nbs-go/nsql/mysql/query"
	"github.com/nbs-go/nsql/schema"
	"github.com/nbs-go/nsql/test_utils"
	"testing"
)

func TestInsert(t *testing.T) {
	// Init schema
	s := schema.New(schema.FromModelRef(new(Person)))

	// Test #1
	test_utils.CompareString(t, "INSERT ALL COLUMNS",
		query.Insert(s, "*").Build(),
		"INSERT INTO `Person`(`createdAt`, `updatedAt`, `fullName`) VALUES (:createdAt, :updatedAt, :fullName)")

	// Test #2
	test_utils.CompareString(t, "INSERT SPECIFIED COLUMNS",
		query.Insert(s, `id`, `fullName`, "age").Build(),
		"INSERT INTO `Person`(`id`, `fullName`) VALUES (:id, :fullName)",
	)

	// Init schema without auto increment
	s = schema.New(schema.FromModelRef(new(Person)), schema.AutoIncrement(false))

	// Test #3
	test_utils.CompareString(t, "INSERT ALL (NO AUTO INCREMENT)",
		query.Insert(s, "*").Build(),
		"INSERT INTO `Person`(`createdAt`, `updatedAt`, `id`, `fullName`) VALUES (:createdAt, :updatedAt, :id, :fullName)",
	)
}
