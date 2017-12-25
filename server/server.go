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
	"github.com/sabhiram/trade-bot/server/hub"
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

	wshub *hub.Hub // websocket hub
}

// New returns an instance of Server.
func New(addr string) (*Server, error) {
	wsh, err := hub.New()
	if err != nil {
		return nil, err
	}

	s := &Server{
		Server: &http.Server{
			Addr: addr,
		},

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

func (s *Server) wsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("wsHandler :: error :: %s\n", err.Error())
			return
		}
		defer c.Close()

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				if _, ok := err.(*websocket.CloseError); !ok {
					fmt.Printf("wsHandler :: unable to read ws :: %s\n", err.Error())
				} else {
					// todo: unsubscribe from hub here...
					fmt.Printf("wsHandler :: connection closed\n")
				}
				break
			}

			switch mt {
			case websocket.TextMessage:
				fmt.Printf("wsHandler :: got message :: %s\n", string(msg))
			default:
				fmt.Printf("wsHandler :: unknown message type :: %d\n", mt)
			}
		}
	}
}

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
