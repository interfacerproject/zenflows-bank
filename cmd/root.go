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
	rootCmd.PersistentFlags().StringVar(&inputFile, "input", "list.csv", "input file, supported format is csv. Defaults to list.csv")
	rootCmd.PersistentFlags().StringVar(&outputFile, "output", "list.csv", "output file, supported format are csv and xlsx. Defaults to list.csv")
}

var rootCmd = &cobra.Command{
	Use:   "bank",
	Short: "Zenflows bank manage token conversion",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Config)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
