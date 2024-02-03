/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
