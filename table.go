package nsql

import "github.com/nbs-go/nsql/schema"

type Table struct {
	Schema *schema.Schema
	As     string
}
