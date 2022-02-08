package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	opt "github.com/nbs-go/nsql/query/option"
	"github.com/nbs-go/nsql/schema"
)

type DeleteBuilder struct {
	schema *schema.Schema
	where  query.WhereWriter
}

func (b *DeleteBuilder) Build(args ...interface{}) string {
	// Get variable format option
	opts := opt.EvaluateOptions(args)
	format, ok := opts.GetVariableFormat()
	if !ok {
		// If var format is not defined, then set default to query.NamedVar
		format = query.BindVar
	}

	// Set variable format in conditions
	if b.where == nil {
		// Set where to id
		b.where = Equal(Column(b.schema.PrimaryKey(), opt.Schema(b.schema)))
	}

	// Set format in conditions
	setUpdateFormat(b.where, b.schema, format)

	// Write where
	where := b.where.WhereQuery()

	return fmt.Sprintf(`DELETE FROM "%s" WHERE %s`, b.schema.TableName(), where)
}

func (b *DeleteBuilder) Where(w query.WhereWriter) *DeleteBuilder {
	b.where = w
	return b
}

func Delete(s *schema.Schema) *DeleteBuilder {
	return &DeleteBuilder{
		schema: s,
	}
}
