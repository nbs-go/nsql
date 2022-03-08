package query_test

import (
	"github.com/nbs-go/nsql/mysql/query"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/schema"
	"github.com/nbs-go/nsql/test_utils"
	"testing"
)

func TestDelete(t *testing.T) {
	// Init schema
	s := schema.New(schema.FromModelRef(new(Transaction)))

	// Test #1
	test_utils.CompareString(t, "DELETE BY ID",
		query.Delete(s).Build(),
		"DELETE FROM `Transaction` WHERE `id` = ?")

	// Test #2
	test_utils.CompareString(t, "DELETE BY MULTIPLE CONDITION",
		query.Delete(s).Where(
			query.And(
				query.Equal(query.Column("id")),
				query.Equal(query.Column("version")),
			),
		).Build(),
		"DELETE FROM `Transaction` WHERE `id` = ? AND `version` = ?")

	// Test #3
	test_utils.CompareString(t, "DELETE BY MULTIPLE CONDITION (NAMED VAR)",
		query.Delete(s).Where(
			query.And(
				query.Equal(query.Column("id")),
				query.Equal(query.Column("version")),
			),
		).Build(option.VariableFormat(op.NamedVar)),
		"DELETE FROM `Transaction` WHERE `id` = :id AND `version` = :version")
}
