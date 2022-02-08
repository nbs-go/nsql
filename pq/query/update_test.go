package query

import (
	"github.com/nbs-go/nsql/query"
	opt "github.com/nbs-go/nsql/query/option"
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
		Update(s, "*").Build(),
		`UPDATE "Transaction" SET "createdAt" = :createdAt, "updatedAt" = :updatedAt, "status" = :status, "version" = :version WHERE "id" = :id`)

	// Test #2
	test_utils.CompareString(t, "UPDATE SPECIFIED COLUMNS",
		Update(s, "status", "price").Build(),
		`UPDATE "Transaction" SET "status" = :status WHERE "id" = :id`,
	)

	// Test #3
	test_utils.CompareString(t, "UPDATE WITH MULTIPLE CONDITION",
		Update(s, "status").
			Where(
				And(
					Equal(Column("id")),
					Equal(Column("version")),
				),
			).
			Build(),
		`UPDATE "Transaction" SET "status" = :status WHERE "id" = :id AND "version" = :version`,
	)

	// Test #4
	test_utils.CompareString(t, "UPDATE WITH BIND VAR",
		Update(s, "status").
			Where(
				And(
					Equal(Column("id")),
					Equal(Column("version")),
				),
			).
			Build(opt.VariableFormat(query.BindVar)),
		`UPDATE "Transaction" SET "status" = ? WHERE "id" = ? AND "version" = ?`,
	)
}
