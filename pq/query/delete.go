package query

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/schema"
)

type DeleteBuilder struct {
	schema *schema.Schema
	where  nsql.WhereWriter
}

func (b *DeleteBuilder) Build(args ...interface{}) string {
	// Get variable format option
	opts := option.EvaluateOptions(args)
	format, ok := opts.GetVariableFormat()
	if !ok {
		// If var format is not defined, then set default to query.NamedVar
		format = nsql.BindVar
	}

	// Set variable format in conditions
	if b.where == nil {
		// Set where to id
		b.where = Equal(Column(b.schema.PrimaryKey(), option.Schema(b.schema)))
	}

	// Set format in conditions
	setUpdateFormat(b.where, b.schema, format)

	// Write where
	where := b.where.WhereQuery()

	return fmt.Sprintf(`DELETE FROM "%s" WHERE %s`, b.schema.TableName(), where)
}

func (b *DeleteBuilder) Where(w nsql.WhereWriter) *DeleteBuilder {
	b.where = w
	return b
}

func Delete(s *schema.Schema) *DeleteBuilder {
	return &DeleteBuilder{
		schema: s,
	}
}
