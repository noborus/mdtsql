package mdtsql

import (
	"fmt"
	"os"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/noborus/trdsql"
)

type table struct {
	header []string
	body   [][]string
}

func (t table) string() string {
	ret := fmt.Sprintf("| %s |\n", strings.Join(t.header, " | "))
	for _, b := range t.body {
		ret += fmt.Sprintf("| %s |\n", strings.Join(b, " | "))
	}
	return ret
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
				t.header = r
			}
		case *ast.TableBody:
			for _, row := range table.GetChildren() {
				t.body = append(t.body, tableRow(row))
			}
		default:
			ast.Print(os.Stdout, node)
		}
	}
	return t
}

type Importer struct {
	tableName  string
	caption    bool
	tableNames []string
	node       ast.Node
}

func (im *Importer) Import(db *trdsql.DB, query string) (string, error) {
	return query, im.nodeAST(db, im.node)
}

func NewImporter(tableName string, md []byte, caption bool) Importer {
	parser := parser.New()
	o := parser.Parse(md)
	im := Importer{
		tableName: tableName,
		caption:   caption,
		node:      o,
	}
	return im
}

func already(tableNames []string, tableName string) bool {
	for _, v := range tableNames {
		if tableName == v {
			return true
		}
	}
	return false
}

func (im *Importer) tableImport(db *trdsql.DB, t table) error {
	l := len(t.header)
	types := make([]string, l)
	for i := 0; i < l; i++ {
		types[i] = "text"
	}
	tableName := db.QuotedName(im.tableName)
	for i := 2; already(im.tableNames, tableName); i++ {
		tableName = db.QuotedName(fmt.Sprintf("%s_%d", im.tableName, i))
	}
	im.tableNames = append(im.tableNames, tableName)
	err := db.CreateTable(tableName, t.header, types, true)
	if err != nil {
		return err
	}
	columnName := strings.Join(t.header, ", ")
	for _, row := range t.body {
		column := "'" + strings.Join(row, "', '") + "'"
		sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnName, column)
		_, err = db.Tx.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func (im *Importer) nodeAST(db *trdsql.DB, node ast.Node) error {
	switch node := node.(type) {
	case *ast.Heading:
		if im.caption {
			im.tableName = text(ast.GetFirstChild(node))
		}
	case *ast.Text:
		if im.caption {
			im.tableName = text(node)
		}
	case *ast.Table:
		err := im.tableImport(db, tableNode(node))
		if err != nil {
			return err
		}
	default:
		for _, node := range node.GetChildren() {
			err := im.nodeAST(db, node)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
