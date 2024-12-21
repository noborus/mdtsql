package mdtsql

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Analyze parses the markdown file and returns the table information.
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

	var reader io.Reader = os.Stdin
	if fileName != "stdin" {
		f, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f
	}

	r := MDTReader{}
	r.caption = Caption
	if err := r.parse(reader); err != nil {
		return nil, err
	}
	if len(r.tables) == 0 {
		return nil, fmt.Errorf("no markdown table found")
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

// Dump outputs the table information.
func Dump(w io.Writer, tables []table) {
	for _, table := range tables {
		fmt.Fprintf(w, "Table Name: [%s]\n", table.tableName)
		typeTable := tablewriter.NewWriter(w)
		typeTable.SetAutoFormatHeaders(false)
		typeTable.SetHeader([]string{"column name", "type"})
		for n, name := range table.names {
			typeTable.Append([]string{name, table.types[n]})
		}
		typeTable.Render()
		fmt.Fprintf(w, "\n")
	}
}
