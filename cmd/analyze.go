package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze markdown",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("require query")
		}
		fileName := args[0]
		im, err := importer(fileName, Caption)
		if err != nil {
			return err
		}
		err = im.Analyze()
		if err != nil {
			return err
		}
		im.Dump(os.Stdout)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
