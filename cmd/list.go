package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/interfacerproject/zenflows-bank/config"
	"github.com/interfacerproject/zenflows-bank/storage"
	"github.com/interfacerproject/zenflows-bank/zenflows"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

func RequestPerson(id string, za zenflows.Agent, rc chan []string) {
	person, err := za.GetPerson(id)

	if err != nil {
		log.Println(err.Error())
		rc <- []string{id}
	} else {
		rc <- []string{id, person.EthereumAddress}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List accounts and their token amount",
	Run: func(cmd *cobra.Command, args []string) {
		storage := &storage.TTStorage{}
		err := storage.Init(config.Config.TTHost, config.Config.TTUser, config.Config.TTPass)
		if err != nil {
			log.Fatal(err.Error())
		}

		za := zenflows.Agent{
			Sk:          config.Config.ZenflowsSk,
			ZenflowsUrl: config.Config.ZenflowsUrl,
		}

		balances, err := storage.Balances()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
		}
		ethChan := make(chan []string)
		for k := range balances {
			go RequestPerson(k, za, ethChan)
		}
		for i := 0; i < len(balances); i++ {
			val := <-ethChan
			if len(val) > 1 {
				balances[val[0]].EthereumAddress = val[1]
			}
		}
		wr := csv.NewWriter(os.Stdout)
		wr.Write([]string{
			"ID",
			"EthereumAddress",
			"Idea",
			"Strengths",
		})
		for k, v := range balances {
			wr.Write([]string{
				k,
				v.EthereumAddress,
				strconv.FormatInt(v.Idea, 10),
				strconv.FormatInt(v.Strengths, 10),
			})
		}
		wr.Flush()
	},
}
