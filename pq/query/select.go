package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
	"github.com/nbs-go/nsql/schema"
	"strings"
)

func Select() *SelectBuilder {
	b := SelectBuilder{
		fields:   []query.SelectFieldWriter{},
		tables:   map[string]string{},
		orderBys: []query.OrderBy{},
	}
	return &b
}

type SelectBuilder struct {
	fields   []query.SelectFieldWriter
	from     query.FromWriter
	where    query.WhereWriter
	tables   map[string]string // To store registered table information
	limit    *int
	skip     *int
	orderBys []query.OrderBy
}

func (b *SelectBuilder) Columns(s *schema.Schema, c1 string, cn ...string) *SelectBuilder {
	// Init columns
	columns := []string{c1}

	// Combine fields if set more than 1
	if len(cn) > 0 {
		columns = append(columns, cn...)
	}

	// Create new schema writer
	w := newColumnWriter(s, columns)

	// add to selected fields
	b.fields = append(b.fields, w)

	return b
}

func (b *SelectBuilder) From(s *schema.Schema, args ...string) *SelectBuilder {
	w := tableWriter{tableName: s.TableName}

	// If from is being overwritten, then panic
	if b.from != nil {
		panic(fmt.Errorf("cannot overwrite existing from"))
	}

	// If alias is set, then add alias
	if len(args) > 0 && args[0] != "" {
		alias := args[0]
		b.tables[s.TableName] = alias
		w.SetAlias(alias)
	} else {
		b.tables[s.TableName] = ""
	}

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

func (b *SelectBuilder) OrderBy(s *schema.Schema, col string, args ...op.SortDirection) *SelectBuilder {
	// Skip if column is not available
	if !s.IsColumnExist(col) {
		return b
	}

	// Get direction
	direction := op.Ascending
	if len(args) > 0 {
		direction = args[0]
	}

	b.orderBys = append(b.orderBys, query.OrderBy{
		TableName: s.TableName,
		Column:    col,
		Direction: direction,
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
	// Generate select query
	fQueries := make([]string, len(b.fields))
	for i, f := range b.fields {
		// If selected fields is from table that is already registered
		if alias, ok := b.tables[f.GetTableName()]; ok {
			// If alias is set, then set alias
			if alias != "" {
				f.SetTableAlias(alias)
			}
			fQueries[i] = f.SelectQuery()
		}
	}
	selectQuery := strings.Join(fQueries, query.Separator)

	// Generate from query
	from := b.from.FromQuery()

	// Combine query
	q := fmt.Sprintf("SELECT %s FROM %s", selectQuery, from)

	// Generate where query
	if b.where != nil {
		// Set table aliases
		setAliasOnWhereWriter(b.where, b.tables)

		// Generate query
		where := b.where.WhereQuery()

		// If not empty, then add
		if where != "" {
			q += fmt.Sprintf(` WHERE %s`, where)
		}
	}

	// Add order by
	if count := len(b.orderBys); count > 0 {
		arr := make([]string, count)
		for i, v := range b.orderBys {
			tableName := v.TableName

			// Override table name if aliases is set
			if alias, ok := b.tables[tableName]; ok && alias != "" {
				tableName = alias
			}

			// Create order by query
			oq := fmt.Sprintf(`"%s"."%s"`, tableName, v.Column)

			// Add direction for descending
			if v.Direction == op.Descending {
				oq += " DESC"
			}

			// Append query
			arr[i] = oq
		}
		orderBy := strings.Join(arr, query.Separator)
		q += fmt.Sprintf(" ORDER BY %s", orderBy)
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

func setAliasOnWhereWriter(ww query.WhereWriter, aliases map[string]string) {
	// Switch type
	switch w := ww.(type) {
	case query.LogicalWhereWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			setAliasOnWhereWriter(cw, aliases)
		}
	case query.ComparisonWhereWriter:
		// Get alias
		if a, ok := aliases[w.GetTableName()]; ok && a != "" {
			w.SetTableAlias(a)
		}
	}
}
