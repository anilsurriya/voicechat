package vcserver

import (
	"log"
	"net/http"
	"strings"
)

func (vs *VcServer) closeRoom(id string) {
	room, found := vs.rooms[roomID(id)]
	if !found {
		return
	}
	room.closeRoom()
	delete(vs.rooms, roomID(id))
}

func (vs *VcServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.URL.Path)
	reqPath := req.URL.Path

	id := ""
	if strings.HasPrefix(reqPath, vs.Prefix) {
		id = reqPath[len(vs.Prefix):]
	}

	if id == "" {
		return
	}

	room, found := vs.rooms[roomID(id)]

	if !found {
		log.Println(id, ": Not found")
		return
	}

	room.ServeHTTP(w, req)
}
