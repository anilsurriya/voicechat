package vcserver

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				if msg.Sender == client {
					continue
				}
				client.send <- msg.Msg
			}
		case <-r.close:
			r.closeRoom()
		}
	}
}

func (r *room) closeRoom() {
	close(r.forward)
	close(r.join)
	close(r.leave)
	close(r.close)
	for client := range r.clients {
		delete(r.clients, client)
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
