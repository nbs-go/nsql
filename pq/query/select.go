package query

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/schema"
	"strings"
)

func Select(column1 nsql.SelectWriter, columnN ...nsql.SelectWriter) *SelectBuilder {
	b := SelectBuilder{
		fields:   []nsql.SelectWriter{},
		orderBys: []nsql.OrderByWriter{},
		tables:   map[string]nsql.Table{},
	}

	// Merge columns
	columns := append([]nsql.SelectWriter{column1}, columnN...)

	// Set to columns to selected fields
	b.fields = columns

	return &b
}

type SelectBuilder struct {
	fields   []nsql.SelectWriter
	from     nsql.FromWriter
	where    nsql.WhereWriter
	orderBys []nsql.OrderByWriter
	limit    *int
	skip     *int
	tables   map[string]nsql.Table
}

func (b *SelectBuilder) From(s *schema.Schema, args ...interface{}) *SelectBuilder {
	// Evaluate options
	opts := option.EvaluateOptions(args)

	// Get alias
	as, _ := opts.GetString(option.AsKey)

	// Create writer
	w := newTableWriter(s.TableName(), as)

	// Add table and set FROM
	b.addTable(s, as)
	b.from = w

	return b
}

func (b *SelectBuilder) Join(s *schema.Schema, onCondition nsql.WhereWriter, args ...interface{}) *SelectBuilder {
	// Evaluate options
	opts := option.EvaluateOptions(args)
	joinMethod := opts.GetJoinMethod()
	as, _ := opts.GetString(option.AsKey)

	// Resolve fromTableFlag reference
	resolveFromTableFlag(onCondition, b.getFromSchema())

	// Resolve joinTableFlag reference
	resolveJoinTableFlag(onCondition, s)

	// Set table aliases
	joinTable := nsql.Table{
		Schema: s,
		As:     as,
	}
	setJoinTableAs(onCondition, &joinTable, b.tables)

	// Create join writer
	var w = joinWriter{
		method:      joinMethod,
		table:       &joinTable,
		onCondition: onCondition,
	}

	// Set join
	b.from.Join(&w)

	// Add table
	b.addTable(s, as)

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

	var tableName string
	if s != nil {
		// Skip if columnWriter is not available
		if !s.IsColumnExist(col) {
			return b
		}
		tableName = s.TableName()
	} else {
		tableName = fromTableFlag
	}

	// Get direction
	direction := opts.GetSortDirection()

	b.orderBys = append(b.orderBys, &orderByWriter{
		ColumnWriter: &columnWriter{
			name:      col,
			tableName: tableName,
		},
		direction: direction,
	})

	return b
}

func (b *SelectBuilder) Limit(n int) *SelectBuilder {
	b.limit = &n
	return b
}

func (b *SelectBuilder) Skip(n int) *SelectBuilder {
	b.skip = &n
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

func (b *SelectBuilder) getTableAlias(tableName string) string {
	t, ok := b.tables[tableName]
	if !ok {
		return ""
	}

	return t.As
}

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
		table, ok := b.tables[tableName]
		if !ok {
			continue
		}

		// If columnWriter is set for all
		if f.IsAllColumns() {
			// Check if columnWriter writer is expandable
			exp, eOk := f.(nsql.Expander)
			if eOk {
				// Expand with given schema and replace writer
				f = exp.Expand(option.Schema(table.Schema))
			}
		}

		// Set alias
		f.SetTableAs(table.As)

		// Push to select writer list
		writers = append(writers, f)
	}

	// Set format if query has joins
	if len(b.tables) > 1 {
		for _, w := range writers {
			w.SetFormat(nsql.SelectJoinColumn)
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
		tableName := f.GetTableName()

		// Get existing table, if not then filter out writer
		table, ok := b.tables[tableName]
		if !ok {
			continue
		}

		// Set alias
		f.SetTableAs(table.As)

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
	w := filterWhereWriters(b.where, b.tables)

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
	fromTable := b.from.GetTableName()
	from := b.tables[fromTable]
	return from.Schema
}

func (b *SelectBuilder) addTable(s *schema.Schema, as string) {
	b.tables[s.TableName()] = nsql.Table{
		Schema: s,
		As:     as,
	}
}
