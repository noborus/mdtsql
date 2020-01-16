package mdtsql

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/noborus/trdsql"
	"github.com/olekukonko/tablewriter"
)

type Importer struct {
	tableName  string
	caption    bool
	tableNames []string
	tables     []table
	node       ast.Node
}

func NewImporter(tableName string, md []byte, caption bool) Importer {
	parser := parser.New()
	im := Importer{
		tableName: tableName,
		caption:   caption,
		node:      parser.Parse(md),
	}
	return im
}

func (im *Importer) Dump(w io.Writer) {
	for i, table := range im.tables {
		fmt.Fprintf(w, "Table Name: [%s]\n", im.tableNames[i])
		typeTable := tablewriter.NewWriter(w)
		typeTable.SetAutoFormatHeaders(false)
		typeTable.SetHeader([]string{"column name", "type"})
		for _, name := range table.names {
			typeTable.Append([]string{name, "text"})
		}
		typeTable.Render()
		fmt.Fprintf(w, "\n")
	}
}

func (im *Importer) Import(db *trdsql.DB, query string) (string, error) {
	err := im.parseNode(im.node)
	if err != nil {
		return "", err
	}
	for i, table := range im.tables {
		err := im.tableImport(db, im.tableNames[i], table)
		if err != nil {
			return "", err
		}
	}
	return query, nil
}

func (im *Importer) Analyze() error {
	err := im.parseNode(im.node)
	if err != nil {
		return err
	}
	return nil
}

func (im *Importer) parseNode(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Heading:
		if im.caption {
			im.tableName = text(ast.GetLastChild(node))
		}
	case *ast.Text:
		if im.caption {
			im.tableName = text(node)
		}
	case *ast.Table:
		tableName := im.tableName
		for i := 2; already(im.tableNames, tableName); i++ {
			tableName = fmt.Sprintf("%s_%d", im.tableName, i)
		}
		im.tableNames = append(im.tableNames, tableName)
		im.tables = append(im.tables, tableNode(node))
	default:
		for _, node := range node.GetChildren() {
			err := im.parseNode(node)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (im *Importer) tableImport(db *trdsql.DB, tableName string, t table) error {
	err := db.CreateTable(db.QuotedName(tableName), t.names, t.types, true)
	if err != nil {
		return err
	}
	return db.Import(db.QuotedName(tableName), t.names, t)
}

func already(tableNames []string, tableName string) bool {
	for _, v := range tableNames {
		if tableName == v {
			return true
		}
	}
	return false
}
