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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return exec(args)
	},
}

func exec(args []string) error {
	if Debug {
		trdsql.EnableDebug()
	}
	query := strings.Join(args, " ")

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
