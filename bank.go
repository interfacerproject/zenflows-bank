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
}

type Config struct {
	ZenflowsUrl string
	TTHost      string
	TTUser      string
	TTPass      string
}

type Token struct {
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

		wr := csv.NewWriter(os.Stdout)
		for k,v := range balances {
			var fabcoin int64 = 0
			if v.Idea > minIdea && v.Strength > minStrength {
				fabcoin = fabcoinAmount
			}
			txid := ""
			wr.Write([]string{
				k,
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
	bank := &Bank{config: config, storage: storage}

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
