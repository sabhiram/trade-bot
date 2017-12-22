package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	bittrex "github.com/toorop/go-bittrex"
)

////////////////////////////////////////////////////////////////////////////////

type ExecFunc func(args map[string]interface{}) error

type Input struct {
	prompt string
	key    string
}

// Trade represents the required data to represent the appropriate
// trading condition.  It contains a list of variables to fetch from
// the user which are stored in a map.  It also contains a template
// expression which when evaluated against the map should resolve in
// a truthy value.
type Trade struct {
	evaluate string
	inputs   []*Input
	execute  ExecFunc
}

var tradeMap = map[string]*Trade{
	"limit-sell": &Trade{
		evaluate: `{{ ge .Current .Target }}`,
		inputs: []*Input{
			{prompt: "Sell Limit", key: "SellLimit"},
			{prompt: "Sell Price", key: "SellPrice"},
		},
		execute: func(args map[string]interface{}) error {
			log.Printf("EXECUTE LIMIT ORDER\nGOT ARGS: %#v\n", args)
			return nil
		},
	},
	"high-low": &Trade{
		evaluate: `{{ }}`,
		execute: func(args map[string]interface{}) error {
			return nil
		},
	},
	"stop-loss": &Trade{
		evaluate: `{{ }}`,
		execute: func(args map[string]interface{}) error {
			return nil
		},
	},
}

////////////////////////////////////////////////////////////////////////////////

var (
	stdinOnce    sync.Once
	stdinScanner *bufio.Scanner
)

func getUserInput(msg string) string {
	stdinOnce.Do(func() {
		stdinScanner = bufio.NewScanner(os.Stdin)
	})

	if stdinScanner != nil {
		fmt.Printf(msg)
		stdinScanner.Scan()
		return stdinScanner.Text()
	}

	return ""
}

////////////////////////////////////////////////////////////////////////////////

func printSummary(s *bittrex.MarketSummary) {
	fmt.Printf(`
Market Summary:
===============
High:       ` + s.High.String() + `
Low:        ` + s.Low.String() + `
Ask:        ` + s.Ask.String() + `
Bid:        ` + s.Bid.String() + `
Last:       ` + s.Last.String() + `
Volume:     ` + s.Volume.String() + `
`)
}

////////////////////////////////////////////////////////////////////////////////

func runCmd(cmd, currency string, market *bittrex.MarketSummary, target, btc, usdt *Balance) error {
	fmt.Printf(`
Available %s balance %f.
Available USDT balance %f.
Available BTC balance %f.
`, currency, target.Available, usdt.Available, btc.Available)
	printSummary(market)

	trade, ok := tradeMap[cmd]
	if !ok {
		return fmt.Errorf("invalid command (%s) specified", cmd)
	}

	// Iterate through the inputs as specified by the trade.
	// Build a map and pass it to the exec function.
	m := map[string]interface{}{}
	for _, inp := range trade.inputs {
		m[inp.key] = getUserInput(inp.prompt)
	}

	// TODO: Evaluate until the chain of conditions are true.
	// Once this is the case we 'execute' the action for the
	// trade.

	// TODO: For now - blindly 'execute' the trade.
	return trade.execute(m)
}

////////////////////////////////////////////////////////////////////////////////
