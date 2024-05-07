package query

import (
	"errors"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/option"
	"strings"
)

func Lower(col nsql.ColumnWriter, args ...interface{}) *LowerColumnWriter {
	if col == nil {
		panic(errors.New("nsql: column cannot be nil"))
	}

	if col.GetColumn() == AllColumns {
		panic(errors.New("nsql: all column (*) is not supported"))
	}

	// Evaluate options
	opts := option.EvaluateOptions(args)

	// Get alias
	as, _ := opts.GetString(option.AsKey)

	return &LowerColumnWriter{
		as:           as,
		ColumnWriter: col,
	}
}

// LowerColumnWriter implements query.SelectWriter that wrap column with LOWER() function
type LowerColumnWriter struct {
	as string
	nsql.ColumnWriter
}

func (c *LowerColumnWriter) SelectQuery() string {
	var b strings.Builder
	b.WriteString("LOWER(")
	b.WriteString(c.ColumnWriter.ColumnQuery())
	b.WriteString(")")

	if c.as != "" {
		b.WriteString(` AS "`)
		b.WriteString(c.as)
		b.WriteString(`"`)
	}

	return b.String()
}

func (c *LowerColumnWriter) ColumnQuery() string {
	if c.as != "" {
		return `"` + c.as + `"`
	}

	var b strings.Builder
	b.WriteString("LOWER(")
	b.WriteString(c.ColumnWriter.ColumnQuery())
	b.WriteString(")")

	return b.String()
}

func (c *LowerColumnWriter) IsAllColumns() bool {
	return false
}
