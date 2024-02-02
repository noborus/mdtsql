package mdtsql

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/yuin/goldmark/ast"
)

type Importer struct {
	tableName string
	caption   bool
	tables    []table
	node      ast.Node
	source    []byte
}

func Analyze(fileName string, caption bool) (*Importer, error) {
	if fileName == "" {
		return nil, fmt.Errorf("require markdown file")
	}
	if fileName == "-" {
		fileName = "stdin"
	}
	if idx := strings.Index(fileName, "::"); idx != -1 {
		fileName = fileName[:idx]
	}
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

	r := MDTReader{}
	if err := r.parse(f); err != nil {
		return nil, err
	}
	im := &Importer{}
	for i, node := range r.tables {
		table, err := tableNode(r.source, node)
		if err != nil {
			return nil, err
		}
		table.tableName = strconv.Itoa(i)
		im.tables = append(im.tables, table)
	}
	return im, nil
}

func (im *Importer) Dump(w io.Writer) {
	for _, table := range im.tables {
		fmt.Fprintf(w, "Table Name: [%s]\n", table.tableName)
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
