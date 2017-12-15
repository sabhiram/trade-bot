package main

////////////////////////////////////////////////////////////////////////////////

import (
	"log"
	"net/http"
	"os"
	"time"

	bittrex "github.com/toorop/go-bittrex"
)

////////////////////////////////////////////////////////////////////////////////

var (
	apiKey string
	secret string
)

////////////////////////////////////////////////////////////////////////////////

func fatalOnError(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %s\n", err.Error())
	}
}

func main() {
	h := &http.Client{
		Timeout: time.Second * 10,
	}
	client := bittrex.NewWithCustomHttpClient(apiKey, secret, h)

	markets, err := client.GetMarkets()
	fatalOnError(err)

	log.Printf("Markets: %#v\n", markets)
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetPrefix("")
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	apiKey = os.Getenv("BITTREX_API_KEY")
	if len(apiKey) == 0 {
		log.Fatalf("BITTREX_API_KEY environment variable missing\n")
	}

	secret = os.Getenv("BITTREX_SECRET")
	if len(secret) == 0 {
		log.Fatalf("BITTREX_SECRET environment variable missing\n")
	}
}

////////////////////////////////////////////////////////////////////////////////
