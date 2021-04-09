package recv

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	sync.Mutex
	conns map[int]*websocket.Conn
}

var rooms = make(map[int]*Room)
var channels = make(map[int]chan string)
var counters = make(map[int]int)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunServer() {
	http.HandleFunc("/ws", handleConnections)

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	params, ok := r.URL.Query()["code"]

	if !ok || len(params[0]) == 0 {
		log.Println("Url Param 'code' is missing")
		return
	}
	roomCodeStr := params[0]
	roomCode, err := strconv.Atoi(roomCodeStr)
	if err != nil {
		ws.WriteMessage(0, []byte("Provide a 5 digit room code"))
		return
	}
	if len(roomCodeStr) != 5 {
		ws.WriteMessage(0, []byte("Provide a 5 digit room code"))
		return
	}

	log.Println("Url Param 'roomCode' is: " + fmt.Sprint(roomCode))

	counters[roomCode] += 1
	userid := counters[roomCode]

	if rooms[roomCode] == nil {
		rooms[roomCode] = &Room{conns: make(map[int]*websocket.Conn)}
	}
	rooms[roomCode].conns[userid] = ws

	if channels[roomCode] == nil {
		channels[roomCode] = make(chan string)
	}
	go handleRoom(roomCode)

	for {
		_, bytes, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(rooms, roomCode)
			break
		}
		msg := string(bytes)
		channels[roomCode] <- msg
	}
}

func handleRoom(roomCode int) {
	room := rooms[roomCode]

	for {
		msg := <-channels[roomCode]
		for userid, ws := range room.conns {
			fmt.Println("Room", roomCode, ": Sending", msg, "to", userid)
			err := ws.WriteMessage(1, []byte(msg))
			if err != nil {
				log.Printf("error: %v", err)
				ws.Close()
				delete(room.conns, userid)
			}
		}
	}
}
