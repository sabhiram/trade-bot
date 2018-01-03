// Package app abstracts all timers and queries that the engine processes.
package app

////////////////////////////////////////////////////////////////////////////////

import (
	bittrex "github.com/toorop/go-bittrex"

	"github.com/sabhiram/trade-bot/app/db"
	"github.com/sabhiram/trade-bot/types"
)

////////////////////////////////////////////////////////////////////////////////

// App encapsulates the checkers for all pending conditions for the given
// session.  The app instance will be used to issue new requests to the upstream
// APIs and push state to various clients via open/subscribed websockets.
type App struct {
	config *types.Config    // app config
	db     *db.DB           // local "database" of tracked session(s)
	client *bittrex.Bittrex // bittrex client
}

func New(config *types.Config) (*App, error) {
	d, err := db.New(config.DbPath)
	if err != nil {
		return nil, err
	}

	app := &App{
		config: config,
		db:     d,
		client: bittrex.New(config.ApiKey, config.Secret),
	}

	return app, nil
}

////////////////////////////////////////////////////////////////////////////////
