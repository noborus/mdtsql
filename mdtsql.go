package mdtsql

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

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

func MarkdownQuery(fileName string, query string, caption bool, w trdsql.Writer) error {
	if fileName == "" {
		fileName = "stdin"
	}

	importer, err := importer(fileName, caption)
	if err != nil {
		return err
	}
	trd := trdsql.NewTRDSQL(
		importer,
		trdsql.NewExporter(
			w,
		),
	)
	return trd.Exec(query)
}

func Analyze(fileName string, caption bool) (*Importer, error) {
	if fileName == "" {
		return nil, fmt.Errorf("require markdown file")
	}
	if fileName == "-" {
		fileName = "stdin"
	}
	im, err := importer(fileName, caption)
	if err != nil {
		return nil, err
	}
	err = im.Analyze()
	if err != nil {
		return nil, err
	}
	return im, nil
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

func (im *Importer) ImportContext(ctx context.Context, db *trdsql.DB, query string) (string, error) {
	err := im.parseNode(im.node)
	if err != nil {
		return "", err
	}
	for i, table := range im.tables {
		err := im.tableImport(ctx, db, im.tableNames[i], table)
		if err != nil {
			return "", err
		}
	}
	return query, nil
}

func (im *Importer) Import(db *trdsql.DB, query string) (string, error) {
	ctx := context.Background()
	return im.ImportContext(ctx, db, query)
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
			im.tableName = toText(node.Children)
		}
	case *ast.Text:
		if im.caption {
			im.tableName = string(node.AsLeaf().Literal)
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

func (im *Importer) tableImport(ctx context.Context, db *trdsql.DB, tableName string, t table) error {
	err := db.CreateTableContext(ctx, db.QuotedName(tableName), t.names, t.types, true)
	if err != nil {
		return err
	}
	return db.ImportContext(ctx, db.QuotedName(tableName), t.names, t)
}

func already(tableNames []string, tableName string) bool {
	for _, v := range tableNames {
		if tableName == v {
			return true
		}
	}
	return false
}

func importer(fileName string, caption bool) (*Importer, error) {
	var f io.Reader
	if fileName != "stdin" {
		var err error
		f, err = os.Open(fileName)
		if err != nil {
			return nil, err
		}
	} else {
		f = os.Stdin
	}
	md, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	tableName := filepath.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))])

	im := NewImporter(tableName, md, caption)
	return &im, nil
}
