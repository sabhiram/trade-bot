package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"text/template"
	"time"

	bittrex "github.com/toorop/go-bittrex"
)

////////////////////////////////////////////////////////////////////////////////

type (
	ExecFunc   func(t *Trade, args map[string]interface{}) error
	UpdateFunc func(t *Trade, args map[string]interface{}) error
)

////////////////////////////////////////////////////////////////////////////////

type Input struct {
	prompt string
	key    string
}

////////////////////////////////////////////////////////////////////////////////

// Trade represents the required data to represent the appropriate
// trading condition.  It contains a list of variables to fetch from
// the user which are stored in a map.  It also contains a template
// expression which when evaluated against the map should resolve in
// a truthy value.
type Trade struct {
	evaluate string
	inputs   []*Input
	execute  ExecFunc
	update   UpdateFunc

	Currency      string
	TargetBalance *Balance
	BTCBalance    *Balance
	USDTBalance   *Balance
}

func (t *Trade) Setup(currency string, target, btc, usdt *Balance) error {
	fmt.Printf("SETUP CALLED: %#v\n", currency)
	t.Currency = currency
	t.TargetBalance = target
	t.BTCBalance = btc
	t.USDTBalance = usdt
	return nil
}

func (t *Trade) Evaluate(args map[string]interface{}) (string, error) {
	fmt.Printf("EVAL CALLED: %#v\n", args)
	tpl, err := template.New("eval").Parse(t.evaluate)
	if err != nil {
		return "error", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, args); err != nil {
		return "error", err
	}

	return buf.String(), nil
}

func (t *Trade) Run(currency string, args map[string]interface{}, refreshDuration time.Duration) error {
	// Always update the trade before doing anything else. This will cause the
	// default values to be setup correctly.  Update should also be called
	// if the evaluate returns false for the next tick.
	if err := t.update(t, args); err != nil {
		return err
	}

	for {
		res, err := t.Evaluate(args)
		fatalOnError(err)

		switch res {
		case "true":
			return t.execute(t, args)
		case "false":
			if err := t.update(t, args); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid result %s", res)
		}

		<-time.After(refreshDuration)
	}

	return errors.New("unknown run error")
}

var tradeMap = map[string]*Trade{
	"limit-sell": &Trade{
		evaluate: `{{ ge .Current .Target }}`,
		inputs: []*Input{
			{prompt: "Sell Limit (in BTC): ", key: "SellLimit"},
			{prompt: "Sell Price (in BTC): ", key: "SellPrice"},
		},
		update: func(t *Trade, args map[string]interface{}) error {
			args["FOO"] = "bar"
			return nil
		},
		execute: func(t *Trade, args map[string]interface{}) error {
			log.Printf("Trade exec called %#v\n", args)
			return nil
		},
	},
	"high-low": &Trade{
		evaluate: `{{ ge 1 1 }}`,
		update: func(t *Trade, args map[string]interface{}) error {
			return nil
		},
		execute: func(t *Trade, args map[string]interface{}) error {
			return nil
		},
	},
	"stop-loss": &Trade{
		evaluate: `{{ ge 1 1 }}`,
		update: func(t *Trade, args map[string]interface{}) error {
			return nil
		},
		execute: func(t *Trade, args map[string]interface{}) error {
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

	if err := trade.Setup(currency, target, btc, usdt); err != nil {
		return err
	}

	// Iterate through the inputs as specified by the trade.
	// Build a map and pass it to the exec function.
	m := map[string]interface{}{}
	for _, inp := range trade.inputs {
		m[inp.key] = getUserInput(inp.prompt)
	}

	return trade.Run(currency, m, cli.refreshInterval)
}

////////////////////////////////////////////////////////////////////////////////
