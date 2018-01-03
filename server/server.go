/* Note: This package auto-generates files based on the below directives:
This one validates the install directory and ensures that the esc binary
exists on the system:
//go:generate go run gen/verify.go

This one takes the "static/" directory and compresses it into
`static/static.go` which is then used as a file handler in the server.
//go:generate esc -o static/static.go -pkg static -prefix "static/" static

To update stuff, run `go generate ./server` from the root of this project.
*/
package server

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/sabhiram/trade-bot/app"
	"github.com/sabhiram/trade-bot/hub"
	"github.com/sabhiram/trade-bot/server/static"
)

////////////////////////////////////////////////////////////////////////////////

const (
	cUseLocalFS = false
)

////////////////////////////////////////////////////////////////////////////////

var (
	wsUpgrader = websocket.Upgrader{}
)

////////////////////////////////////////////////////////////////////////////////

// Server handles all websocket, HTTP API and file requests.
type Server struct {
	*http.Server

	app   *app.App // app engine
	wshub *hub.Hub // websocket hub
}

// New returns an instance of Server.
func New(addr string, a *app.App) (*Server, error) {
	wsh, err := hub.New()
	if err != nil {
		return nil, err
	}

	s := &Server{
		Server: &http.Server{
			Addr: addr,
		},

		app:   a,
		wshub: wsh,
	}

	return s, s.setupRoutes()
}

func (s *Server) Start() {
	fmt.Printf("Kicking off webserver at: %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("error :: webserver died :: %s\n", err.Error())
	}
}

////////////////////////////////////////////////////////////////////////////////

func (s *Server) todoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("TODO handler hit!\n")
		w.Write([]byte("TODO"))
	}
}

func (s *Server) setupRoutes() error {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(static.FS(cUseLocalFS)))
	mux.Handle("/ws", s.wsHandler())
	mux.Handle("/api", s.todoHandler())

	s.Handler = mux
	return nil
}

////////////////////////////////////////////////////////////////////////////////
