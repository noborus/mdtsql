package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/noborus/mdtsql"
	"github.com/noborus/trdsql"
)

func importer(fileName string, caption bool) (*mdtsql.Importer, error) {
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

	if Debug {
		trdsql.EnableDebug()
	}
	im := mdtsql.NewImporter(tableName, md, caption)
	return &im, nil
}

func exec(fileName string, query string, caption bool) error {
	importer, err := importer(fileName, caption)
	if err != nil {
		return err
	}
	trd := trdsql.NewTRDSQL(
		importer,
		trdsql.NewExporter(
			outFormat(),
		),
	)
	return trd.Exec(query)
}
