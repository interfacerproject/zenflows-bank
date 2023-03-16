package cmd

import (
	"fmt"
	"os"

	"github.com/interfacerproject/zenflows-bank/config"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(config.Init)
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
