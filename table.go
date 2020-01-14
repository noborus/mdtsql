package mdtsql

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/noborus/trdsql"
)

type table struct {
	names []string
	types []string
	body  [][]interface{}
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

func text(node ast.Node) string {
	if node == nil {
		return ""
	}
	l := (node).AsLeaf()
	if l == nil {
		return ""
	}
	return string(l.Literal)
}

func tableCell(node ast.Node) string {
	switch node := node.(type) {
	case *ast.TableCell:
		return text(ast.GetFirstChild(node))
	default:
		return ""
	}
}

func tableRow(node ast.Node) []string {
	row := []string{}
	switch node := node.(type) {
	case *ast.TableRow:
		for _, n := range node.GetChildren() {
			row = append(row, tableCell(n))
		}
		return row
	default:
		return []string{}
	}
}

func tableNode(node ast.Node) table {
	t := table{}
	for _, table := range node.GetChildren() {
		switch table := table.(type) {
		case *ast.TableHeader:
			for _, row := range table.GetChildren() {
				r := tableRow(row)
				for i, col := range r {
					if col == "" {
						r[i] = fmt.Sprintf("c%d", i+1)
					}
				}
				t.names = r
			}
		case *ast.TableBody:
			for _, row := range table.GetChildren() {
				r := tableRow(row)
				data := make([]interface{}, len(r))
				for i, col := range r {
					data[i] = col
				}
				t.body = append(t.body, data)
			}
		default:
			// ast.Print(os.Stdout, node)
		}
	}
	t.types = make([]string, len(t.names))
	for i := 0; i < len(t.names); i++ {
		t.types[i] = trdsql.DefaultDBType
	}
	return t
}
