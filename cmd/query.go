package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/noborus/mdtsql"
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
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		md, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		trdsql.EnableDebug()
		d := mdtsql.NewIm(filepath.Base(path[:len(path)-len(filepath.Ext(path))]), md)
		trd := trdsql.NewTRDSQL(&d,
			trdsql.NewExporter(
				outFormat(),
			),
		)
		err = trd.Exec(Query)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
