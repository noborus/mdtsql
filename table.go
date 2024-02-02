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

func (t table) Names() ([]string, error) {
	return t.names, nil
}

func (t table) Types() ([]string, error) {
	return t.types, nil
}

func (t table) PreReadRow() [][]interface{} {
	return t.body
}

// ReadRow only returns EOF.
func (t table) ReadRow(row []interface{}) ([]interface{}, error) {
	return nil, io.EOF
}
