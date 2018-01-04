package socket

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

////////////////////////////////////////////////////////////////////////////////

type Socket struct {
	conn   *websocket.Conn
	sendCh chan []byte
}

func New(c *websocket.Conn) *Socket {
	return &Socket{
		conn:   c,
		sendCh: make(chan []byte, 1024),
	}
}

////////////////////////////////////////////////////////////////////////////////

func (s *Socket) Read() {
	defer s.conn.Close()
	s.conn.SetReadLimit(1024)
	s.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	s.conn.SetPongHandler(func(string) error {
		s.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		fmt.Printf("Waiting for read...\n")
		mt, msg, err := s.conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); !ok {
				fmt.Printf("wsHandler :: unable to read ws :: %s\n", err.Error())
			} else {
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

func (s *Socket) Write() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		s.conn.Close()
	}()

	for {
		fmt.Printf("Waiting for write on sendCh or ticker...\n")
		select {
		case msg, ok := <-s.sendCh:
			s.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			// Handle case when the hub / app closes a socket.
			if !ok {
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Write the message to the next websocket writer available.
			w, err := s.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			// Drain the buffered send channel.
			n := len(s.sendCh)
			for i := 0; i < n; i++ {
				w.Write(<-s.sendCh)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			s.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

func (s *Socket) Send(msg []byte) {
	s.sendCh <- msg
}

func (s *Socket) Close() {
	close(s.sendCh)
}

////////////////////////////////////////////////////////////////////////////////
