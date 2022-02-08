package query

import (
	"github.com/nbs-go/nsql/query"
	opt "github.com/nbs-go/nsql/query/option"
	"github.com/nbs-go/nsql/schema"
)

func Schema(s *schema.Schema) *SchemaBuilder {
	return &SchemaBuilder{
		schema: s,
	}
}

type SchemaBuilder struct {
	schema *schema.Schema
}

func (s *SchemaBuilder) Schema() *schema.Schema {
	return s.schema
}

func (s *SchemaBuilder) FindByPK() string {
	return Select().From(s.schema).Where(Equal(Column(s.schema.PrimaryKey()))).Build()
}

func (s *SchemaBuilder) Insert() string {
	return Insert(s.schema, AllColumns).Build()
}

func (s *SchemaBuilder) Update() string {
	where := Equal(Column(s.schema.PrimaryKey()))
	return Update(s.schema, AllColumns).Where(where).Build()
}

func (s *SchemaBuilder) Delete() string {
	where := Equal(Column(s.schema.PrimaryKey()))
	return Delete(s.schema).Where(where).Build()
}

func (s *SchemaBuilder) Count(where query.WhereWriter) string {
	return Select(opt.Count(s.schema.PrimaryKey(), opt.As("count"))).From(s.schema).Where(where).Build()
}
