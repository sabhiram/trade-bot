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
      trade-bot -pair <pair> [-refresh 5s] <command>

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
		pair string   // currency trading pair
		args []string // other command line args

		refreshInterval time.Duration // conditions check refresh interval
		apiKey          string        // bittrex api key
		secret          string        // bittrex secret
	}{}
)

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
		log.Printf(usage)
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
		Timeout: time.Second * 10,
	}
	client := bittrex.NewWithCustomHttpClient(cli.apiKey, cli.secret, h)

	for {
		summary, err := client.GetMarketSummary(cli.pair)
		fatalOnError(err)

		if len(summary) > 0 {
			lastValue := summary[0].Last.String()
			log.Printf("Last value: %s\n", lastValue)
		}

		<-time.After(cli.refreshInterval)
	}
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetPrefix("")
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	cli.apiKey = getenvFatal("BITTREX_API_KEY")
	cli.secret = getenvFatal("BITTREX_SECRET")

	flag.StringVar(&cli.pair, "pair", "", "trading pair to look-up")
	flag.StringVar(&cli.pair, "p", "", "trading pair to look-up (short)")

	var refIntStr string
	flag.StringVar(&refIntStr, "refresh", "5s", "refresh interval duration")
	flag.StringVar(&refIntStr, "r", "5s", "refresh interval duration (short)")

	flag.Parse()

	if len(cli.pair) == 0 {
		usageErr(fmt.Errorf("required argument -pair/-p missing"))
	}

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
