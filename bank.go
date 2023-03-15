package main

import (
	"fmt"
	"log"
	"os"
	"flag"
	"encoding/csv"
	"strconv"
)

type Bank struct {
	config Config
	storage *TTStorage
	dryRun bool
	zenflowsAgent ZenflowsAgent
}

type Config struct {
	ZenflowsUrl string
	TTHost      string
	TTUser      string
	TTPass      string
}

type Token struct {
	EthereumAddress string
	Idea     int64
	Strength int64
}

func loadEnvConfig() Config {
	return Config{
		ZenflowsUrl: fmt.Sprintf("%s/api", os.Getenv("ZENFLOWS_URL")),
		TTHost:      os.Getenv("TT_HOST"),
		TTUser:      os.Getenv("TT_USER"),
		TTPass:      os.Getenv("TT_PASS"),
	}
}

func RequestPerson(id string, bank *Bank, rc chan []string) *http.Response {
	person, err := bank.zenflowsAgent.GetPerson(id)
	if err != nil {
		panic(err)
	}
	rc <- []string{id, person.EthereumAddress}
}

var cmds map[string]func(*Bank, []string) = map[string]func(*Bank, []string) {
	"list": func(bank *Bank, _ []string) {
		balances, err := bank.storage.Balances()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
		}
		wr := csv.NewWriter(os.Stdout)
		for k,v := range balances {
			wr.Write([]string{k, strconv.FormatInt(v.Idea, 10), strconv.FormatInt(v.Strength, 10)})
		}
		wr.Flush()
	},
	"airdrop": func(bank *Bank, values []string) {
		var err error
		var minIdea, minStrength, fabcoinAmount int64
		if len(values) == 0 {
			fmt.Printf("Usage: %s airdrop min_idea min_strength fabcoin_amount\n", os.Args[0])
			os.Exit(1)
		}
		minIdea, err = strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
		}
		minStrength, err = strconv.ParseInt(values[1], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
		}
		fabcoinAmount, err = strconv.ParseInt(values[2], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
		}

		balances, err := bank.storage.Balances()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
		}
		ethChan := make(chan []string)
		for i := 0; i < len(balances); i++ {
			val := <-ethChan
			balances[val[0]].EthereumAddress = val[1]
		}

		wr := csv.NewWriter(os.Stdout)
		for k,v := range balances {
			var fabcoin int64 = 0
			if v.Idea > minIdea && v.Strength > minStrength {
				fabcoin = fabcoinAmount
			}
			txid := ""
			wr.Write([]string{
				v.EthereumAddress,
				strconv.FormatInt(v.Idea, 10),
				strconv.FormatInt(v.Strength, 10),
				strconv.FormatInt(fabcoin, 10),
				txid,
			})
		}
		wr.Flush()
	},
}

func main() {
	config := loadEnvConfig()
	storage := &TTStorage{}
	err := storage.Init(config.TTHost, config.TTUser, config.TTPass)
	if err != nil {
		log.Fatal(err.Error())
	}
	za := ZenflowsAgent{
		Sk:          os.Getenv("ZENFLOWS_SK"),
		ZenflowsUrl: config.ZenflowsUrl,
	}
	bank := &Bank{config: config, storage: storage, zenflowsAgent: za}

	dryRun := flag.Bool("dry-run", false, "Run in dry-run mode")
	flag.Parse()

	bank.dryRun = *dryRun

	values := flag.Args()
	if len(values) == 0 {
		fmt.Printf("Usage: %s {list|amount} ...\n", os.Args[0])
        flag.PrintDefaults()
        os.Exit(1)
	}

	if v, ok := cmds[values[0]]; ok {
		v(bank, values[1:])
    } else {
		fmt.Printf("Unknown command %s", values[0])
        os.Exit(1)
    }

}
