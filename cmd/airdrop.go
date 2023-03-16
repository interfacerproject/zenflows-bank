package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

var minIdea int64
var minStrengths int64
var fabcoinsAmount int64

func init() {
	rootCmd.AddCommand(airdropCmd)
	airdropCmd.Flags().Int64Var(&minIdea, "minIdea", 10, "minimum amount of idea. Defaults to 10")
	airdropCmd.Flags().Int64Var(&minStrengths, "minStrengths", 10, "minimum amount of strengths. Defaults to 10")
	airdropCmd.Flags().Int64Var(&fabcoinsAmount, "fabcoinsAmount", 100, "fabcoinsAmount to transfer")
}

var airdropCmd = &cobra.Command{
	Use:   "airdrop",
	Short: "Make user earn fabcoins based on idea and strengths points",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(minIdea)
		fmt.Println(minStrengths)
		fmt.Println(fabcoinsAmount)
	},
}

