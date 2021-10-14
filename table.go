package mdtsql

import (
	"fmt"
	"io"
	"os"

	"github.com/noborus/trdsql"
	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/extension/ast"
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

func toText(buf []byte) string {
	if len(buf) > 0 {
		return string(buf)
	}
	return ""
}

func (im *Importer) tableNode(node ast.Node) table {
	t := table{}
	for n := node.FirstChild(); n != nil; n = n.NextSibling() {
		switch n.Kind() {
		case gast.KindTableHeader:
			i := 0
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				col := toText(c.Text(im.source))
				if col == "" {
					col = fmt.Sprintf("c%d", i+1)
				}
				t.names = append(t.names, col)
				i++
			}
		case gast.KindTableRow:
			row := []string{}
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				rawText := []byte{}
				for i := 0; i < c.Lines().Len(); i++ {
					line := c.Lines().At(i)
					rawText = append(rawText, line.Value(im.source)...)
				}
				row = append(row, string(rawText))
			}
			data := make([]interface{}, len(row))
			for i, col := range row {
				data[i] = col
			}
			t.body = append(t.body, data)
		default:
			fmt.Fprintf(os.Stderr, "unknown node:")
			fmt.Fprintf(os.Stderr, "%v:%v\n", n.Kind(), n.Type())
		}
	}
	t.types = make([]string, len(t.names))
	for i := 0; i < len(t.names); i++ {
		t.types[i] = trdsql.DefaultDBType
	}
	return t
}
