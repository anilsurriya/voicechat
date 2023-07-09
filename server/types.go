package vcserver

import (
	"strconv"

	"github.com/gorilla/websocket"
)

type roomID string

type VcServer struct {
	counter int
	Addr    string
	Prefix  string
	rooms   map[roomID]*room
}

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

type Message struct {
	Msg    []byte
	Sender *client
}

type room struct {
	forward chan Message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	close   chan bool
	id      int
}

func NewServer(addr string, prefix string) *VcServer {
	return &VcServer{
		counter: 0,
		Addr:    addr,
		Prefix:  prefix,
		rooms:   make(map[roomID]*room),
	}
}

func (vs *VcServer) NewRoom() *room {
	nroom := &room{
		id:      vs.counter,
		forward: make(chan Message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		close:   make(chan bool),
	}
	vs.rooms[roomID(strconv.Itoa(vs.counter))] = nroom
	vs.counter++

	go nroom.Run()

	return nroom
}
