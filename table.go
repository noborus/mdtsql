package mdtsql

import (
	"io"
)

type table struct {
	tableName string
	names     []string
	types     []string
	body      [][]interface{}
}

// Names returns the column names.
func (t table) Names() ([]string, error) {
	return t.names, nil
}

// Types returns the column types.
func (t table) Types() ([]string, error) {
	return t.types, nil
}

// PreReadRow returns the body.
func (t table) PreReadRow() [][]any {
	return t.body
}

// ReadRow only returns EOF.
func (t table) ReadRow(row []any) ([]any, error) {
	return nil, io.EOF
}
