package server

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

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

////////////////////////////////////////////////////////////////////////////////
