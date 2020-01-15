package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "SQL query command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("require query")
		}
		query := args[0]
		fileName := "stdin"
		if len(args) >= 2 {
			fileName = args[1]
		}
		return exec(fileName, query, Caption)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
