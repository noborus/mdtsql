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
				t.header = tableRow(row)
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

func (im *Im) tableImport(db *trdsql.DB, t table) error {
	im.num += 1
	l := len(t.header)
	types := make([]string, l)
	for i := 0; i < l; i++ {
		types[i] = "text"
	}
	fmt.Println(t.header)
	fmt.Println(types)
	tableName := fmt.Sprintf("%s_%d", im.tableName, im.num)
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

func (im *Im) nodeAST(db *trdsql.DB, node ast.Node) error {
	switch node := node.(type) {
	case *ast.Document:
		for _, node := range node.GetChildren() {
			err := im.nodeAST(db, node)
			if err != nil {
				return err
			}
		}
	case *ast.Heading:
		//fmt.Printf("%s\n", text(ast.GetFirstChild(node)))
	case *ast.Table:
		err := im.tableImport(db, tableNode(node))
		if err != nil {
			return err
		}
	default:
		//fmt.Printf("%T\n", node)
	}

	return nil
}

type Im struct {
	tableName string
	node      ast.Node
	num       int
}

func (im Im) Import(db *trdsql.DB, query string) (string, error) {
	return query, im.nodeAST(db, im.node)
}

func NewIm(tableName string, md []byte) Im {
	parser := parser.New()
	o := parser.Parse(md)
	im := Im{
		tableName: tableName,
		node:      o,
	}
	return im
}
