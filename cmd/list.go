package cmd

import (
	"log"
	"os"

	"github.com/noborus/mdtsql"
	"github.com/noborus/trdsql"
	"github.com/spf13/cobra"
)

// listCmd represents the list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List and analyze SQL dumps",
	Long: `List and analyze SQL dumps from a specified file.
 This command parses the SQL dump file and displays information about the tables contained within.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName := ""
		if len(args) >= 1 {
			fileName = args[0]
		}
		if err := analyzeDump(fileName); err != nil {
			log.Fatal(err)
		}
	},
}

func analyzeDump(fileName string) error {
	tables, err := mdtsql.Analyze(fileName)
	if err == nil {
		mdtsql.Dump(os.Stdout, tables)
		return nil
	}

	// fallback to trdsql.Analyze.
	if err := trdsql.Analyze(fileName, trdsql.NewAnalyzeOpts(), trdsql.NewReadOpts()); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
