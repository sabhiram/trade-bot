// Package app abstracts all timers and queries that the engine processes.
package app

////////////////////////////////////////////////////////////////////////////////

import (
	"strings"

	bittrex "github.com/toorop/go-bittrex"

	"github.com/sabhiram/trade-bot/app/db"
	"github.com/sabhiram/trade-bot/hub"
	"github.com/sabhiram/trade-bot/server/socket"
	"github.com/sabhiram/trade-bot/types"
)

////////////////////////////////////////////////////////////////////////////////

// App encapsulates the checkers for all pending conditions for the given
// session.  The app instance will be used to issue new requests to the upstream
// APIs and push state to various clients via open/subscribed websockets.
type App struct {
	config *types.Config // app config
	db     *db.DB        // local "database" of tracked session(s)
	hub    *hub.Hub
	client *bittrex.Bittrex // bittrex client
}

func New(config *types.Config, h *hub.Hub) (*App, error) {
	d, err := db.New(config.DbPath)
	if err != nil {
		return nil, err
	}

	app := &App{
		config: config,
		db:     d,
		hub:    h,
		client: bittrex.New(config.ApiKey, config.Secret),
	}

	return app, app.UpdateBalances(false)
}

////////////////////////////////////////////////////////////////////////////////

// UpdateBalances updates balances using the bitterx client and pushes the new
// balances to the `db`.
func (a *App) UpdateBalances(broadcast bool) error {
	// Fetch current balances from bittrex.
	balances, err := a.client.GetBalances()
	if err != nil {
		return err
	}

	// Pull relevant balances into our own format.
	bs := []*types.Balance{}
	for _, b := range balances {
		ava, _ := b.Available.Float64()
		bal, _ := b.Balance.Float64()
		if bal > 0.0 {
			bs = append(bs, &types.Balance{
				Currency:  strings.ToUpper(b.Currency),
				Available: ava,
				Total:     bal,
			})
		}
	}

	// Update the db.
	err = a.db.UpdateBalances(bs)
	if err != nil {
		return err
	}

	// Update clients (if we need to broadcast).
	if broadcast {
		return a.BroadcastBalances()
	}
	return nil
}

// BroadcastBalances pushes the latest balance state to all connected clients.
func (a *App) BroadcastBalances() error {
	bal, err := a.db.GetBalances()
	if err != nil {
		return err
	}

	bs, err := types.NewSocketMessage("Balance", bal).Marshal()
	if err != nil {
		return err
	}

	a.hub.Broadcast(bs)
	return nil
}

// SendBalances pushes the latest balance state to the specified socket.
func (a *App) SendBalances(sock *socket.Socket) error {
	bal, err := a.db.GetBalances()
	if err != nil {
		return err
	}

	bs, err := types.NewSocketMessage("Balance", bal).Marshal()
	if err != nil {
		return err
	}

	sock.Send(bs)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
