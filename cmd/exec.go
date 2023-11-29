package cmd

import (
	"os"

	"github.com/noborus/mdtsql"
	"github.com/noborus/trdsql"
)

func queryExec(fileName string, query string, caption bool) error {
	w := newWriter(os.Stdout, os.Stderr)
	if Debug {
		trdsql.EnableDebug()
	}

	return mdtsql.MarkdownQuery(fileName, query, caption, w)
}
