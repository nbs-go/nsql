package query

import (
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/option"
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
	return Select(Column("*")).From(s.schema).Where(Equal(Column(s.schema.PrimaryKey()))).Build()
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

func (s *SchemaBuilder) Count(where nsql.WhereWriter) string {
	return Select(Count(s.schema.PrimaryKey(), option.As("count"))).From(s.schema).Where(where).Build()
}

func (s *SchemaBuilder) IsExists(where nsql.WhereWriter) string {
	return Select(GreaterThan(Count(s.schema.PrimaryKey()), IntVar(0), option.As("isExists"))).
		From(s.schema).Where(where).Build()
}
