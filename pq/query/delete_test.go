package query

import (
	"github.com/nbs-go/nsql/query"
	opt "github.com/nbs-go/nsql/query/option"
	"github.com/nbs-go/nsql/schema"
	"github.com/nbs-go/nsql/test_utils"
	"testing"
)

func TestDelete(t *testing.T) {
	// Init schema
	s := schema.New(schema.FromModelRef(new(Transaction)))

	// Test #1
	test_utils.CompareString(t, "DELETE BY ID",
		Delete(s).Build(),
		`DELETE FROM "Transaction" WHERE "id" = ?`)

	// Test #2
	test_utils.CompareString(t, "DELETE BY MULTIPLE CONDITION",
		Delete(s).Where(
			And(
				Equal(Column("id")),
				Equal(Column("version")),
			),
		).Build(),
		`DELETE FROM "Transaction" WHERE "id" = ? AND "version" = ?`)

	// Test #3
	test_utils.CompareString(t, "DELETE BY MULTIPLE CONDITION (NAMED VAR)",
		Delete(s).Where(
			And(
				Equal(Column("id")),
				Equal(Column("version")),
			),
		).Build(opt.VariableFormat(query.NamedVar)),
		`DELETE FROM "Transaction" WHERE "id" = :id AND "version" = :version`)
}
