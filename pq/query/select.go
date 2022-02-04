package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/schema"
	"strings"
)

func Select() *SelectBuilder {
	b := SelectBuilder{
		fields:     []query.SelectFieldWriter{},
		aliases:    map[string]string{},
		aliasesRev: map[string]string{},
		orderBys:   []query.OrderBy{},
	}
	return &b
}

type SelectBuilder struct {
	fields     []query.SelectFieldWriter
	from       query.FromWriter
	aliases    map[string]string
	aliasesRev map[string]string
	limit      *int
	skip       *int
	orderBys   []query.OrderBy
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

	// If alias is set, then add alias
	if len(args) > 0 && args[0] != "" {
		alias := args[0]
		b.aliases[s.TableName] = alias
		b.aliasesRev[alias] = s.TableName
		w.SetAlias(alias)
	}

	b.from = &w

	return b
}

func (b *SelectBuilder) OrderBy(s *schema.Schema, col string, args ...query.SortDirection) *SelectBuilder {
	// Skip if column is not available
	if !s.IsColumnExist(col) {
		return b
	}

	// Get direction
	direction := query.Ascending
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
		if a, ok := b.aliases[f.GetTableName()]; ok {
			f.SetTableAlias(a)
		}
		fQueries[i] = f.SelectQuery()
	}
	selectQuery := strings.Join(fQueries, query.Separator)

	// Generate from query
	from := b.from.FromQuery()

	// Combine query
	q := fmt.Sprintf("SELECT %s FROM %s", selectQuery, from)

	// Add order by
	if count := len(b.orderBys); count > 0 {
		arr := make([]string, count)
		for i, v := range b.orderBys {
			tableName := v.TableName
			// Override table name if aliases is set
			if a, ok := b.aliases[tableName]; ok {
				tableName = a
			}

			// Create order by query
			oq := fmt.Sprintf(`"%s"."%s"`, tableName, v.Column)

			// Add direction for descending
			if v.Direction == query.Descending {
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
