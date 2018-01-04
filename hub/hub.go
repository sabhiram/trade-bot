package hub

////////////////////////////////////////////////////////////////////////////////

import (
	"github.com/sabhiram/trade-bot/server/socket"
)

////////////////////////////////////////////////////////////////////////////////

type Hub struct {
	sockets map[*socket.Socket]struct{}

	broadcastCh  chan []byte
	registerCh   chan *socket.Socket
	unregisterCh chan *socket.Socket
}

func New() (*Hub, error) {
	return &Hub{
		sockets: map[*socket.Socket]struct{}{},

		broadcastCh:  make(chan []byte),
		registerCh:   make(chan *socket.Socket),
		unregisterCh: make(chan *socket.Socket),
	}, nil
}

////////////////////////////////////////////////////////////////////////////////

func (h *Hub) RegisterSocket(s *socket.Socket) {
	h.registerCh <- s
}

func (h *Hub) UnregisterSocket(s *socket.Socket) {
	h.unregisterCh <- s
}

func (h *Hub) Broadcast(msg []byte) {
	h.broadcastCh <- msg
}

func (h *Hub) Run() {
	for {
		select {
		case socket := <-h.registerCh:
			h.sockets[socket] = struct{}{}
		case socket := <-h.unregisterCh:
			if _, ok := h.sockets[socket]; ok {
				delete(h.sockets, socket)
				socket.Close()
			}
		case msg := <-h.broadcastCh:
			for socket := range h.sockets {
				socket.Send(msg)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
