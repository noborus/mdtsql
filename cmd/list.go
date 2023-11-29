package cmd

import (
	"log"
	"os"

	"github.com/noborus/mdtsql"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName := ""
		if len(args) >= 1 {
			fileName = args[0]
		}
		if err := analyzeDump(fileName, Caption); err != nil {
			log.Fatal(err)
		}
	},
}

func analyzeDump(fileName string, caption bool) error {
	im, err := mdtsql.Analyze(fileName, caption)
	if err != nil {
		return err
	}
	im.Dump(os.Stdout)
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
