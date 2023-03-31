package cmd

import (
	"fmt"
	"os"

	"github.com/interfacerproject/zenflows-bank/config"
	"github.com/spf13/cobra"
)

var inputFile string
var outputFile string

func init() {
	cobra.OnInitialize(config.Init)
	rootCmd.PersistentFlags().StringVar(&inputFile, "input", "list.csv", "input file, supported format are csv and xlsx. Defaults to list.csv")
	rootCmd.PersistentFlags().StringVar(&outputFile, "output", "result.csv", "output file, supported format are csv and xlsx. Defaults to result.csv")
}

var rootCmd = &cobra.Command{
	Use:   "bank",
	Short: "Zenflows bank manage token conversion",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
