package cmd

import (
	"os"

	"github.com/noborus/mdtsql"
	"github.com/noborus/trdsql"
)

func queryExec(fileName string, query string, caption bool) error {
	w := outFormat(os.Stdout, os.Stderr)
	if Debug {
		trdsql.EnableDebug()
	}

	return mdtsql.MarkdownQuery(fileName, query, caption, w)
}

func analyzeDump(fileName string, caption bool) error {
	im, err := mdtsql.Analyze(fileName, caption)
	if err != nil {
		return err
	}
	im.Dump(os.Stdout)
	return nil
}
