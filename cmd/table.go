package cmd

import (
	"os"

	"github.com/noborus/trdsql"
	"github.com/spf13/cobra"
)

// tableCmd represents the table command
var tableCmd = &cobra.Command{
	Use: "table",

	Short: "SQL(SELECT * FROM table) for markdown table and tabular data",
	Long: `Execute SQL(SELECT * FROM table) queries on markdown table and tabular data.
This command allows you to run SQL queries against tables formatted in Markdown within the specified files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return execTable(args[0])
	},
}

func execTable(table string) error {
	if Debug {
		trdsql.EnableDebug()
	}
	query := "SELECT * FROM " + table
	writer := newWriter(os.Stdout, os.Stderr)
	trd := trdsql.NewTRDSQL(
		trdsql.NewImporter(
			trdsql.InHeader(Header),
			trdsql.InPreRead(100),
		),
		trdsql.NewExporter(writer),
	)
	return trd.Exec(query)
}

func init() {
	rootCmd.AddCommand(tableCmd)
}
