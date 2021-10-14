package mdtsql

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/noborus/trdsql"
	"github.com/olekukonko/tablewriter"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	gast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

type Importer struct {
	tableName  string
	caption    bool
	tableNames []string
	tables     []table
	node       ast.Node
	source     []byte
}

func NewImporter(tableName string, md []byte, caption bool) Importer {
	gmd := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	parser := gmd.Parser()
	im := Importer{
		tableName: tableName,
		caption:   caption,
		node:      parser.Parse(text.NewReader(md)),
		source:    md,
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
	if err = im.Analyze(); err != nil {
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
	if err := im.parseNode(im.node); err != nil {
		return "", err
	}
	for i, table := range im.tables {
		if err := im.tableImport(ctx, db, im.tableNames[i], table); err != nil {
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
	if err := im.parseNode(im.node); err != nil {
		return err
	}
	return nil
}

func (im *Importer) parseNode(node ast.Node) error {
	switch node.Type() {
	case ast.TypeDocument:
		for n := node.FirstChild(); n != nil; n = n.NextSibling() {
			if err := im.parseNode(im.node); err != nil {
				return err
			}
		}
	case ast.TypeBlock:
		if node.Kind() == gast.KindTable {
			im.tables = append(im.tables, im.tableNode(node))
			tableName := im.tableName
			for i := 2; already(im.tableNames, tableName); i++ {
				tableName = fmt.Sprintf("%s_%d", im.tableName, i)
			}
			im.tableNames = append(im.tableNames, tableName)
			return nil
		}

		switch node.Kind() {
		case ast.KindHeading, ast.KindParagraph:
			if im.caption {
				im.tableName = string(node.Text(im.source))
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown node:")
		fmt.Fprintf(os.Stderr, "%v:%v\n", node.Kind(), node.Type())
	}
	return nil
}

func (im *Importer) tableImport(ctx context.Context, db *trdsql.DB, tableName string, t table) error {
	if err := db.CreateTableContext(ctx, db.QuotedName(tableName), t.names, t.types, true); err != nil {
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
