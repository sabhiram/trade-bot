package main

////////////////////////////////////////////////////////////////////////////////

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	bittrex "github.com/toorop/go-bittrex"
)

////////////////////////////////////////////////////////////////////////////////

const usage = `trade-bot usage:

  $ BITTREX_API_KEY=<key> BITTREX_SECRET=<secret> \
      trade-bot [-refresh 5s] <command>

  All 'command's will trigger once a target is hit.  Each command will
  query the user for all required parameters and confirm before being
  deployed.  Once deployed the script will query the last trade value
  for the 'pair' every 'refresh' seconds (default 5s).  Once any/all of
  the conditions are met, the appropriate action is taken.

  Where 'command' can include:

    limit-sell          -   regular on-limit sale
    stop-loss           -   simple stop loss
    high-low            -   stop-loss + limit-sell (first one wins)

  For authorizing transactions and conducting market queries, the API
  key and secret need to be provided as environment variables. The two
  required ones are:

    BITTREX_API_KEY     -   api-key from bittrex
    BITTREX_SECRET      -   secret for above key

  The api-key and secret can be accessed once you have logged in to
  bittrex in the "Settings" section under "Manage API Keys".

  Note: You must have 2-factor authentication enabled to make new keys

  Issues:
  =======

  If you find this software useful, help out by filing issues or
  suggestions here: https://github.com/sabhiram/trade-bot/issues.

  Disclaimer:
  ===========

  This is open source software.  Use at your own risk.

`

const version = `0.0.1`

////////////////////////////////////////////////////////////////////////////////

var (
	cli = struct {
		args []string // other command line args

		refreshInterval time.Duration // conditions check refresh interval
		apiKey          string        // bittrex api key
		secret          string        // bittrex secret
	}{}
)

////////////////////////////////////////////////////////////////////////////////

type Balance struct {
	Currency  string
	Available float64
	Total     float64
}

////////////////////////////////////////////////////////////////////////////////

func fatalOnError(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %s\n", err.Error())
	}
}

func usageErr(err error) {
	if err != nil {
		log.Fatalf("Usage Error: %s\n%s", err.Error(), usage)
	} else {
		fmt.Printf(usage)
	}
}

// getenvFatal returns the os env var matching `key` if found.  Otherwise
// throws a fatal error.
func getenvFatal(key string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		usageErr(fmt.Errorf("%s environment variable missing", key))
	}
	return v
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	cmd := strings.ToLower(cli.args[0])
	switch cmd {
	case "version":
		fmt.Printf("%s\n", version)
		return
	case "usage", "help", "-h":
		usageErr(nil)
		return
	case "limit-sell", "stop-loss", "high-low":
		break
	default:
		usageErr(fmt.Errorf("%s is an invalid command", cmd))
		return
	}

	h := &http.Client{
		Timeout: time.Second * 20,
	}
	client := bittrex.NewWithCustomHttpClient(cli.apiKey, cli.secret, h)

	fmt.Printf("Fetching bittrex balances for account...\n")
	balances, err := client.GetBalances()
	fatalOnError(err)

	bs := []*Balance{}
	for _, b := range balances {
		ava, _ := b.Available.Float64()
		bal, _ := b.Balance.Float64()
		if bal > 0.0 {
			bs = append(bs, &Balance{
				Currency:  strings.ToUpper(b.Currency),
				Available: ava,
				Total:     bal,
			})
		}
	}

	fmt.Printf("Found the following balances:\n")
	for i, bal := range bs {
		fmt.Printf("% 3d. % 6s : %f available\n", i+1, bal.Currency, bal.Available)
	}

	input := getUserInput(`Which coin do you want to setup (ex: "PIVX"): `)
	input = strings.ToUpper(input)

	var (
		target *Balance // chosen currency balance
		btc    *Balance // BTC balance
		usdt   *Balance // USDT balance
	)
	for _, bal := range bs {
		switch bal.Currency {
		case input:
			target = bal
		case "BTC":
			btc = bal
		case "USDT":
			usdt = bal
		}
	}

	if target == nil {
		fmt.Printf("Currency (%s) not available.\n", input)
		os.Exit(0)
	}

	if target.Available > 0.0 {
		sourceCurrency := "BTC" // TODO: Fix this
		market := fmt.Sprintf("%s-%s", sourceCurrency, input)

		// TODO: This moves into the runCmd
		fmt.Printf("Querying market %s\n", market)
		summary, err := client.GetMarketSummary(market)
		fatalOnError(err)

		if len(summary) > 0 {
			err := runCmd(cmd, input, &summary[0], target, usdt, btc)
			fatalOnError(err)
		} else {
			fmt.Printf("Market summary does not exist for %s\n", market)
		}
	} else {
		fmt.Printf("Currency (%s) has no available balance\n", input)
	}
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetPrefix("")
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	cli.apiKey = getenvFatal("BITTREX_API_KEY")
	cli.secret = getenvFatal("BITTREX_SECRET")

	var refIntStr string
	flag.StringVar(&refIntStr, "refresh", "5s", "refresh interval duration")
	flag.StringVar(&refIntStr, "r", "5s", "refresh interval duration (short)")

	flag.Parse()

	var err error
	cli.refreshInterval, err = time.ParseDuration(refIntStr)
	if err != nil {
		usageErr(err)
	}

	cli.args = flag.Args()
	if len(cli.args) == 0 {
		cli.args = []string{"usage"}
	}
}

////////////////////////////////////////////////////////////////////////////////
