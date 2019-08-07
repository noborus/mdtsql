package psutilsql

import (
	"github.com/noborus/trdsql"
)

func SliceQuery(slice interface{}, tableName string, query string, out trdsql.Format) error {
	// trdsql.EnableDebug()
	importer := trdsql.NewSliceImporter(tableName, slice)
	writer := trdsql.NewWriter(trdsql.OutFormat(out))
	trd := trdsql.NewTRDSQL(importer, trdsql.NewExporter(writer))
	err := trd.Exec(query)
	return err
}

func readerQuery(reader Reader, query string, out trdsql.Format) error {
	importer, err := NewMultiImporter(reader)
	if err != nil {
		return err
	}
	writer := trdsql.NewWriter(trdsql.OutFormat(out))
	trd := trdsql.NewTRDSQL(importer, trdsql.NewExporter(writer))
	err = trd.Exec(query)
	return err
}
