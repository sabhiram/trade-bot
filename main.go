package main

////////////////////////////////////////////////////////////////////////////////

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sabhiram/trade-bot/server"
)

////////////////////////////////////////////////////////////////////////////////

const (
	sitePath = "./site/build/default"
)

////////////////////////////////////////////////////////////////////////////////

const usage = `trade-bot usage:

  $ BITTREX_API_KEY=<key> BITTREX_SECRET=<secret> trade-bot

  For authorizing transactions and conducting market queries, the API
  key and secret need to be provided as environment variables. The two
  required ones are:

    BITTREX_API_KEY     -   api-key from bittrex
    BITTREX_SECRET      -   secret for above key

  The api-key and secret can be accessed once you have logged in to
  bittrex in the "Settings" section under "Manage API Keys".

  Note: You must have 2-factor authentication enabled to make new keys.

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
	s, err := server.New(":8100")
	fatalOnError(err)

	s.Start()
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
