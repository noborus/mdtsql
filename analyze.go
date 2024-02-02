package mdtsql

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func Analyze(fileName string) ([]table, error) {
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
	r.caption = Caption
	if err := r.parse(f); err != nil {
		return nil, err
	}
	tables := make([]table, 0, len(r.tables))
	for i, node := range r.tables {
		table, err := tableNode(r.source, node)
		if err != nil {
			return nil, err
		}
		if r.caption {
			table.tableName = r.tableNames[i]
		} else {
			table.tableName = strconv.Itoa(i)
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func Dump(w io.Writer, tables []table) {
	for _, table := range tables {
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
