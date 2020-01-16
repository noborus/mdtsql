package mdtsql

import (
	"fmt"
	"io"
	"os"

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

func toText(nodes []ast.Node) string {
	var ret string
	for _, node := range nodes {
		switch node := node.(type) {
		case *ast.Text, *ast.Code:
			l := (node).AsLeaf()
			if l == nil {
				continue
			}
			ret += string(l.Literal)
		case *ast.Link:
			ret += toText(node.Children)
		default:
			fmt.Fprintf(os.Stderr, "unknown node:%#v\n", node)
		}
	}
	return ret
}

func tableCell(node ast.Node) string {
	switch node := node.(type) {
	case *ast.TableCell:
		return toText(node.Children)
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
		fmt.Fprintf(os.Stderr, "unknown node:%#v\n", node)
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
		case *ast.TableBody, *ast.TableFooter:
			for _, row := range table.GetChildren() {
				r := tableRow(row)
				data := make([]interface{}, len(r))
				for i, col := range r {
					data[i] = col
				}
				t.body = append(t.body, data)
			}
		default:
			ast.Print(os.Stderr, node)
		}
	}
	t.types = make([]string, len(t.names))
	for i := 0; i < len(t.names); i++ {
		t.types[i] = trdsql.DefaultDBType
	}
	return t
}
