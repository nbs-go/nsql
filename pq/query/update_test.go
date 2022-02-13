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

type Transaction struct {
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
	Id        int64     `db:"id"`
	Status    string    `db:"status"`
	Version   int64     `db:"version"`
}

func TestUpdate(t *testing.T) {
	// Init schema
	s := schema.New(schema.FromModelRef(new(Transaction)))

	// Test #1
	test_utils.CompareString(t, "UPDATE ALL COLUMNS",
		query.Update(s, "*").Build(),
		`UPDATE "Transaction" SET "createdAt" = :createdAt, "updatedAt" = :updatedAt, "status" = :status, "version" = :version WHERE "id" = :id`)

	// Test #2
	test_utils.CompareString(t, "UPDATE SPECIFIED COLUMNS",
		query.Update(s, "status", "price").Build(),
		`UPDATE "Transaction" SET "status" = :status WHERE "id" = :id`,
	)

	// Test #3
	test_utils.CompareString(t, "UPDATE WITH MULTIPLE CONDITION",
		query.Update(s, "status").
			Where(
				query.And(
					query.Equal(query.Column("id")),
					query.Equal(query.Column("version")),
				),
			).
			Build(),
		`UPDATE "Transaction" SET "status" = :status WHERE "id" = :id AND "version" = :version`,
	)

	// Test #4
	test_utils.CompareString(t, "UPDATE WITH BIND VAR",
		query.Update(s, "status").
			Where(
				query.And(
					query.Equal(query.Column("id")),
					query.Equal(query.Column("version")),
				),
			).
			Build(option.VariableFormat(op.BindVar)),
		`UPDATE "Transaction" SET "status" = ? WHERE "id" = ? AND "version" = ?`,
	)
}
