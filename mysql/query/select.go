package query

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/schema"
	"log"
	"strings"
)

func newSelectBuilder() *SelectBuilder {
	return &SelectBuilder{
		fields:    []nsql.SelectWriter{},
		orderBys:  []nsql.OrderByWriter{},
		schemaRef: map[schema.Reference]*schema.Schema{},
	}
}

func Select(column1 nsql.SelectWriter, columnN ...nsql.SelectWriter) *SelectBuilder {
	b := newSelectBuilder()
	b.Select(column1, columnN...)
	return b
}

func From(s *schema.Schema, args ...interface{}) *SelectBuilder {
	b := newSelectBuilder()
	b.From(s, args...)
	return b
}

type SelectBuilder struct {
	fields    []nsql.SelectWriter
	from      nsql.FromWriter
	where     nsql.WhereWriter
	orderBys  []nsql.OrderByWriter
	limit     *int64
	skip      *int64
	schemaRef map[schema.Reference]*schema.Schema
}

func (b *SelectBuilder) Select(column1 nsql.SelectWriter, columnN ...nsql.SelectWriter) *SelectBuilder {
	// Set existing
	b.fields = append([]nsql.SelectWriter{column1}, columnN...)
	return b
}

func (b *SelectBuilder) From(s *schema.Schema, args ...interface{}) *SelectBuilder {
	// Resolve as option, and print warning
	opts := option.EvaluateOptions(args)
	as, _ := opts.GetString(option.AsKey)
	if as != "" {
		log.Printf("nsql: warning: From() option setter option.As() is deprecated. Use schema.New() option setter schema.As() instead. See Breaking Changes Note => https://github.com/nbs-go/nsql#breaking-changes. (Schema = %s)\n", s.TableName())
	}
	// Create writer
	w := newTableWriter(s.TableName(), s.As())
	// Add table and set FROM
	b.addTable(s)
	b.from = w
	return b
}

func (b *SelectBuilder) Join(s *schema.Schema, onCondition nsql.WhereWriter, args ...interface{}) *SelectBuilder {
	// Evaluate options
	opts := option.EvaluateOptions(args)
	joinMethod := opts.GetJoinMethod()

	// Resolve as option, and print warning
	as, _ := opts.GetString(option.AsKey)
	if as != "" {
		log.Printf("nsql: warning: Join() option setter option.As() is deprecated. Use schema.New() option setter schema.As() instead. See Breaking Changes Note => https://github.com/nbs-go/nsql#breaking-changes. (Schema = %s)\n", s.TableName())
	}

	// Add table
	b.addTable(s)

	// Resolve fromTableFlag reference
	resolveFromTableFlag(onCondition, b.getFromSchema())

	// Resolve joinTableFlag reference
	resolveJoinTableFlag(onCondition, s)

	// Set table aliases
	joinTable := s
	setJoinTableAs(onCondition, joinTable, b.schemaRef)

	// Create join writer
	var w = joinWriter{
		method:      joinMethod,
		table:       joinTable,
		onCondition: onCondition,
	}

	// Set join
	b.from.Join(&w)

	return b
}

func (b *SelectBuilder) Where(w1 nsql.WhereWriter, wn ...nsql.WhereWriter) *SelectBuilder {
	// If wn is not set, then set where filter
	if len(wn) == 0 {
		b.where = w1
		return b
	}

	// Else, set where with AND logical operators
	where := append([]nsql.WhereWriter{w1}, wn...)
	b.where = &whereLogicWriter{
		op:         op.And,
		conditions: where,
	}
	return b
}

func (b *SelectBuilder) OrderBy(col string, args ...interface{}) *SelectBuilder {
	// Evaluate options
	opts := option.EvaluateOptions(args)

	// Get schema
	s := opts.GetSchema()

	var tableName, tableAs string
	if s != nil {
		// Skip if columnWriter is not available
		if !s.IsColumnExist(col) {
			return b
		}
		tableName = s.TableName()
		tableAs = s.As()
	} else {
		tableName = fromTableFlag
	}

	// Get direction
	direction := opts.GetSortDirection()

	b.orderBys = append(b.orderBys, &orderByWriter{
		ColumnWriter: &columnWriter{
			name:      col,
			tableName: tableName,
			tableAs:   tableAs,
		},
		direction: direction,
	})

	return b
}

func (b *SelectBuilder) ResetOrderBy() *SelectBuilder {
	b.orderBys = []nsql.OrderByWriter{}
	return b
}

func (b *SelectBuilder) Limit(n int64) *SelectBuilder {
	b.limit = &n
	return b
}

func (b *SelectBuilder) ResetLimit() *SelectBuilder {
	b.limit = nil
	return b
}

func (b *SelectBuilder) Skip(n int64) *SelectBuilder {
	b.skip = &n
	return b
}

func (b *SelectBuilder) ResetSkip() *SelectBuilder {
	b.skip = nil
	return b
}

func (b *SelectBuilder) Build() string {
	selectQuery := b.writeSelectQuery()

	// Generate from query
	from := b.from.FromQuery()

	// Combine query
	q := fmt.Sprintf("SELECT %s FROM %s", selectQuery, from)

	// Generate where query
	whereQuery := b.writeWhereQuery()
	if whereQuery != "" {
		q += whereQuery
	}

	// Generate order by query
	orderBy := b.writeOrderByQuery()
	if orderBy != "" {
		q += orderBy
	}

	// Add limit
	if b.limit != nil {
		q += fmt.Sprintf(" LIMIT %d", *b.limit)
	}

	// Add skip
	if b.skip != nil {
		q += fmt.Sprintf(" OFFSET %d", *b.skip)
	}

	return q
}

// Private methods

func (b *SelectBuilder) writeSelectQuery() string {
	// Replace fromTableFlag with FROM Table Name
	from := b.getFromSchema()
	for _, f := range b.fields {
		if f.GetTableName() == fromTableFlag {
			f.SetSchema(from)
		}
	}

	// Prepare query.SelectWriter
	var writers []nsql.SelectWriter
	for _, f := range b.fields {
		// Get table name
		tableName := f.GetTableName()

		// If force flag is set, then add to writers list
		if tableName == forceWriteFlag {
			writers = append(writers, f)
		}

		// Get existing table, if not then filter out writer
		sRef := f.GetSchemaRef()
		table, ok := b.schemaRef[sRef]
		if !ok {
			continue
		}

		// If columnWriter is set for all
		if f.IsAllColumns() {
			// Check if columnWriter writer is expandable
			exp, eOk := f.(nsql.Expander)
			if eOk {
				// Expand with given schema and replace writer
				f = exp.Expand(option.Schema(table))
			}
		}

		// Set alias
		f.SetTableAs(table.As())

		// Push to select writer list
		writers = append(writers, f)
	}

	// Set format if query has joins
	if len(b.schemaRef) > 1 {
		for _, w := range writers {
			w.SetFormat(op.SelectJoinColumn)
		}
	}

	// Generate select query
	fields := make([]string, len(writers))
	for i, w := range writers {
		fields[i] = w.SelectQuery()
	}

	return strings.Join(fields, nsql.Separator)
}

func (b *SelectBuilder) writeOrderByQuery() string {
	// If empty, then return empty query
	if len(b.orderBys) == 0 {
		return ""
	}

	// Replace fromTableFlag with FROM Table Name
	from := b.getFromSchema()
	for _, f := range b.orderBys {
		if f.GetTableName() == fromTableFlag {
			f.SetSchema(from)
		}
	}

	// Prepare order by writers
	var writers []nsql.OrderByWriter
	for _, f := range b.orderBys {
		// Get existing table, if not then filter out writer
		sRef := f.GetSchemaRef()
		table, ok := b.schemaRef[sRef]
		if !ok {
			continue
		}

		// Set alias
		f.SetTableAs(table.As())

		writers = append(writers, f)
	}

	// If no writers, then return empty
	if len(writers) == 0 {
		return ""
	}

	// Generate query
	arr := make([]string, len(writers))
	for i, w := range writers {
		arr[i] = w.OrderByQuery()
	}
	q := strings.Join(arr, nsql.Separator)
	return fmt.Sprintf(" ORDER BY %s", q)
}

func (b *SelectBuilder) writeWhereQuery() string {
	if b.where == nil {
		return ""
	}

	// Replace fromTableFlag with FROM Table Name
	from := b.getFromSchema()
	resolveFromTableFlag(b.where, from)

	// Prepare query.SelectWriter
	w := filterWhereWriters(b.where, b.schemaRef)

	if w == nil {
		return ""
	}

	q := w.WhereQuery()

	if q == "" {
		return q
	}

	return " WHERE " + q
}

// getFromSchema retrieve schema that is defined in FROM
func (b *SelectBuilder) getFromSchema() *schema.Schema {
	sRef := b.from.GetSchemaRef()
	return b.schemaRef[sRef]
}

func (b *SelectBuilder) addTable(s *schema.Schema) {
	b.schemaRef[s.Ref()] = s
}
