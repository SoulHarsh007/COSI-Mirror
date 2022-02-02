package mirrormap

import (
	"github.com/COSI_Lab/Mirror/logging"
)

type hub struct {
	// Map of clients TODO: Make a HashSet?
	clients map[*client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// registers a client from the hub
	register chan *client

	// unregister a client from the hub
	unregister chan *client
}

func (hub *hub) count() int {
	return len(hub.clients)
}

func (hub *hub) run() {
	for {
		select {
		case client := <-hub.register:
			// registers a client
			hub.clients[client] = true
			logging.Log(logging.Info, "Registered client", client.conn.RemoteAddr())
		case client := <-hub.unregister:
			// unregister a client
			delete(hub.clients, client)
			close(client.send)
			logging.Log(logging.Info, "Unregistered client", client.conn.RemoteAddr())
		case message := <-hub.broadcast:
			// broadcasts the message to all clients
			for client := range hub.clients {
				select {
				case client.send <- message:
				default:
					// if sending to a client blocks we drop the client
					close(client.send)
					delete(hub.clients, client)
					logging.Log(logging.Warn, "Dropped client", client.conn.RemoteAddr())
				}
			}
		}
	}
}
