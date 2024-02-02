package cmd

import (
	"os"
	"strings"

	"github.com/noborus/trdsql"
	"github.com/spf13/cobra"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Execute SQL queries on markdown table and tabular data",
	Long: `Execute SQL queries on markdown table and tabular data.
This command allows you to run SQL queries against tables formatted in Markdown within the specified files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return exec(args)
	},
}

func exec(args []string) error {
	if Debug {
		trdsql.EnableDebug()
	}
	query := strings.Join(args, " ")
	trdsql.EnableMultipleQueries()
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
	rootCmd.AddCommand(queryCmd)
}
