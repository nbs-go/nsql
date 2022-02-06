package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
	opt "github.com/nbs-go/nsql/query/option"
	"github.com/nbs-go/nsql/schema"
	"strings"
)

type selectTable struct {
	schema *schema.Schema
	as     string
}

func Select(args ...interface{}) *SelectBuilder {
	b := SelectBuilder{
		fields:   []query.SelectWriter{},
		orderBys: []query.OrderByWriter{},
		tables:   map[string]selectTable{},
	}

	b.Select(args...)

	return &b
}

type SelectBuilder struct {
	fields   []query.SelectWriter
	from     query.FromWriter
	where    query.WhereWriter
	orderBys []query.OrderByWriter
	limit    *int
	skip     *int
	tables   map[string]selectTable
}

func (b *SelectBuilder) Select(args ...interface{}) *SelectBuilder {
	// Evaluate options
	opts := opt.EvaluateOptions(args)

	// Check if select option
	if countCol, ok := opts.GetString(opt.CountKey); ok && countCol != "" {
		b.selectCount(countCol, opts)
		return b
	}

	// Get table columnSchemaWriter
	var tableName string
	s := opts.GetSchema()

	if s == nil {
		// If schema not set, then will set selected fields using table that is defined in "FROM"
		tableName = fromTableFlag
	} else {
		tableName = s.TableName
	}

	// Get columnSchemaWriter
	inputColumns := opts.GetStringArray(opt.ColumnsKey)

	// If input has columnSchemaWriter
	if count := len(inputColumns); count > 0 {
		// If all columnSchemaWriter is set, then add single columnWriter
		if inputColumns[0] == AllColumns {
			b.fields = append(b.fields, &columnWriter{
				name:      inputColumns[0],
				tableName: tableName,
			})
		} else {
			// Set columnSchemaWriter as schema columnSchemaWriter
			b.fields = append(b.fields, &columnSchemaWriter{
				schema:    s,
				columns:   inputColumns,
				tableName: tableName,
			})
		}
	}

	return b
}

func (b *SelectBuilder) selectCount(column string, options *opt.Options) {
	// Get options
	s := options.GetSchema()
	as, _ := options.GetString(opt.AsKey)

	// Init writer
	w := selectCountWriter{as: as}

	// If all column, then create
	if column == "*" {
		w.ColumnWriter = &columnWriter{
			name:      column,
			tableName: forceWriteFlag,
		}
		w.allColumn = true
	} else {
		// If column is invalid or not in schema, then skip
		if !s.IsColumnExist(column) {
			return
		}

		// Set writer
		w.ColumnWriter = &columnWriter{
			name:      column,
			tableName: s.TableName,
		}
	}

	// Set column to select fields
	b.fields = append(b.fields, &w)
}

func (b *SelectBuilder) From(s *schema.Schema,
	args ...interface{}) *SelectBuilder {
	// Evaluate options
	opts := opt.EvaluateOptions(args)

	// Get alias
	as, _ := opts.GetString(opt.AsKey)

	// Create writer
	w := tableWriter{
		tableName: s.TableName,
		as:        as,
	}

	// Add table and set FROM
	b.addTable(s, as)
	b.from = &w

	return b
}

func (b *SelectBuilder) Where(w1 query.WhereWriter, wn ...query.WhereWriter) *SelectBuilder {
	// If wn is not set, then set where filter
	if len(wn) == 0 {
		b.where = w1
		return b
	}

	// Else, set where with AND logical operators
	where := append([]query.WhereWriter{w1}, wn...)
	b.where = &whereLogicalWriter{
		op:         op.And,
		conditions: where,
	}
	return b
}

func (b *SelectBuilder) OrderBy(col string, args ...interface{}) *SelectBuilder {
	// Evaluate options
	opts := opt.EvaluateOptions(args)

	// Get schema
	s := opts.GetSchema()

	var tableName string
	if s != nil {
		// Skip if columnWriter is not available
		if !s.IsColumnExist(col) {
			return b
		}
		tableName = s.TableName
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
	//// Generate where query
	//if b.where != nil {
	//	// Set table aliases
	//	setAliasOnWhereWriter(b.where, b.tables)
	//
	//	// Generate query
	//	where := b.where.WhereQuery()
	//
	//	// If not empty, then add
	//	if where != "" {
	//		q += fmt.Sprintf(` WHERE %s`, where)
	//	}
	//}
	//
	//// Add order by

	//}
	//
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

	return t.as
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
	var writers []query.SelectWriter
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
			exp, eOk := f.(query.Expander)
			if eOk {
				// Expand with given schema and replace writer
				f = exp.Expand(opt.Schema(table.schema))
			}
		}

		// Set alias
		f.SetTableAs(table.as)

		// Push to select writer list
		writers = append(writers, f)
	}

	// Generate select query
	fields := make([]string, len(writers))
	for i, w := range writers {
		fields[i] = w.SelectQuery()
	}

	return strings.Join(fields, query.Separator)
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
	var writers []query.OrderByWriter
	for _, f := range b.orderBys {
		tableName := f.GetTableName()

		// Get existing table, if not then filter out writer
		table, ok := b.tables[tableName]
		if !ok {
			continue
		}

		// Set alias
		f.SetTableAs(table.as)

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
	q := strings.Join(arr, query.Separator)
	return fmt.Sprintf(" ORDER BY %s", q)
}

func (b *SelectBuilder) writeWhereQuery() string {
	if b.where == nil {
		return ""
	}

	// Replace fromTableFlag with FROM Table Name
	from := b.getFromSchema()
	resolveFromOfWhereWriters(b.where, from)

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
	return from.schema
}

func (b *SelectBuilder) addTable(s *schema.Schema, as string) {
	b.tables[s.TableName] = selectTable{
		schema: s,
		as:     as,
	}
}
