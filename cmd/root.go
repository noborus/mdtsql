package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/noborus/mdtsql"
	"github.com/noborus/trdsql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "mdtsql",
	Args:  cobra.MinimumNArgs(1),
	Short: "Execute SQL for markdown table",
	Long: `Execute SQL for table in markdown.
The result can be output to CSV, JSON, LTSV, YAML, Markdown, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		if Ver {
			fmt.Printf("mdtsql version %s rev:%s\n", Version, Revision)
			return
		}
		cmd.Help()
	},
}

var (
	// Version represents the version.
	Version string
	// Revision set "git rev-parse --short HEAD".
	Revision string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string, revision string) {
	Version = version
	Revision = revision
	if err := rootCmd.Execute(); err != nil {
		rootCmd.SetOutput(os.Stderr)
		rootCmd.Println(err)
		os.Exit(1)
	}
}

// Header is an output header specification(CSV and RAW only).
var Header bool

// Ver is version information.
var Ver bool

// Debug is debug print.
var Debug bool

// Delimiter is a delimiter specification (CSV ans RAW only).
var Delimiter string

// OutFormat is an output format specification.
var OutFormat string

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mdtsql.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Ver, "version", "v", false, "display version information")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "debug print")
	rootCmd.PersistentFlags().BoolVarP(&mdtsql.Caption, "caption", "c", false, "caption table name")

	rootCmd.PersistentFlags().StringVarP(&OutFormat, "OutFormat", "o", "md", "output format=at|csv|ltsv|json|jsonl|tbln|raw|md|vf|yaml")
	rootCmd.PersistentFlags().StringVarP(&Delimiter, "Delimiter", "d", ",", "output delimiter (CSV only)")
	rootCmd.PersistentFlags().BoolVarP(&Header, "Header", "O", false, "output header (CSV only)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mdtsql" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mdtsql")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func newWriter(outStream io.Writer, errStream io.Writer) trdsql.Writer {
	format := trdsql.OutputFormat(strings.ToUpper(OutFormat))
	w := trdsql.NewWriter(
		trdsql.OutFormat(format),
		trdsql.OutDelimiter(Delimiter),
		trdsql.OutHeader(Header),
		trdsql.OutStream(outStream),
		trdsql.ErrStream(errStream),
	)
	return w
}
