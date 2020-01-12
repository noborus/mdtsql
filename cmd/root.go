package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/noborus/trdsql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mdtsql",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.SetOutput(os.Stderr)
		rootCmd.Println(err)
		os.Exit(1)
	}
}

// Header is an output header specification(CSV and RAW only).
var Header bool

// Delimiter is a delimiter specification (CSV ans RAW only).
var Delimiter string

// OutFormat is an output format specification.
var OutFormat string

// Query is SQL specification.
var Query string

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mdtsql.yaml)")

	rootCmd.PersistentFlags().StringVarP(&OutFormat, "OutFormat", "o", "at", "output format=at|csv|ltsv|json|tbln|raw|md|vf")
	rootCmd.PersistentFlags().StringVarP(&Delimiter, "Delimiter", "d", ",", "output delimiter (CSV only)")
	rootCmd.PersistentFlags().BoolVarP(&Header, "Header", "O", false, "output header (CSV only)")

	rootCmd.PersistentFlags().StringVarP(&Query, "Query", "q", "", "query")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
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

func outFormat() trdsql.Writer {
	var format trdsql.Format
	switch strings.ToUpper(OutFormat) {
	case "CSV":
		format = trdsql.CSV
	case "LTSV":
		format = trdsql.LTSV
	case "JSON":
		format = trdsql.JSON
	case "TBLN":
		format = trdsql.TBLN
	case "RAW":
		format = trdsql.RAW
	case "MD":
		format = trdsql.MD
	case "AT":
		format = trdsql.AT
	case "VF":
		format = trdsql.VF
	case "JSONL":
		format = trdsql.JSONL
	default:
		format = trdsql.AT
	}
	w := trdsql.NewWriter(
		trdsql.OutFormat(format),
		trdsql.OutDelimiter(Delimiter),
		trdsql.OutHeader(Header),
	)
	return w
}
