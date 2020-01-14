package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/noborus/mdtsql"
	"github.com/noborus/trdsql"
)

func exec(fileName string, query string, caption bool) error {
	var f io.Reader
	if fileName != "stdin" {
		var err error
		f, err = os.Open(fileName)
		if err != nil {
			return err
		}
	} else {
		f = os.Stdin
	}
	md, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	tableName := filepath.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))])

	if Debug {
		trdsql.EnableDebug()
	}

	importer := mdtsql.NewImporter(tableName, md, caption)

	trd := trdsql.NewTRDSQL(
		&importer,
		trdsql.NewExporter(
			outFormat(),
		),
	)
	return trd.Exec(query)
}
