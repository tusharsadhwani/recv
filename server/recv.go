package recv

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	sync.Mutex
	conns map[int]*websocket.Conn
}

type RoomID = int

type Message struct {
	userid int
	text   string
}

var rooms = make(map[RoomID]*Room)
var channels = make(map[RoomID]*chan Message)
var counters = make(map[RoomID]int)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //TODO: CORS
	},
}

func setupCORS(w *http.ResponseWriter) {
	// (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func HandleConnect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for {
		roomCode := 10000 + rand.Intn(90000)
		if rooms[roomCode] == nil {
			rooms[roomCode] = &Room{conns: make(map[int]*websocket.Conn)}
			channel := make(chan Message)
			channels[roomCode] = &channel
			go handleRoom(roomCode)
			w.Write([]byte(fmt.Sprint(roomCode)))
			return
		}
	}
}

func HandleWebsockets(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	params, ok := r.URL.Query()["code"]

	if !ok || len(params[0]) == 0 {
		ws.WriteMessage(websocket.TextMessage, []byte("Url Param 'code' is missing"))
		return
	}
	roomCodeStr := params[0]
	roomCode, err := strconv.Atoi(roomCodeStr)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Provide a 5 digit room code"))
		return
	}
	if len(roomCodeStr) != 5 {
		ws.WriteMessage(websocket.TextMessage, []byte("Provide a 5 digit room code"))
		return
	}
	if rooms[roomCode] == nil {
		ws.WriteMessage(websocket.TextMessage, []byte("This room doesn't exist"))
		close(*channels[roomCode])
		delete(rooms, roomCode)
		delete(channels, roomCode)
		return
	}
	if channels[roomCode] == nil {
		ws.WriteMessage(websocket.TextMessage, []byte("This room doesn't exist"))
		if rooms[roomCode] != nil {
			for _, conn := range rooms[roomCode].conns {
				conn.Close()
			}
		}
		return
	}

	counters[roomCode] += 1
	userid := counters[roomCode]
	rooms[roomCode].conns[userid] = ws

	for {
		_, bytes, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error while reading message: %v", err)
			ws.Close()
			delete(rooms[roomCode].conns, userid)
			if len(rooms[roomCode].conns) == 0 {
				delete(rooms, roomCode)
				close(*channels[roomCode])
				delete(channels, roomCode)
			}
			break
		}

		msg := Message{
			userid: userid,
			text:   string(bytes),
		}
		*channels[roomCode] <- msg
	}
}

func handleRoom(roomCode int) {
	room := rooms[roomCode]
	channel := *channels[roomCode]

	for msg := range channel {
		for userid, ws := range room.conns {
			if msg.userid == userid {
				continue
			}

			// fmt.Printf("Room %d: Sending %q to %d\n", roomCode, msg, userid)
			err := ws.WriteMessage(websocket.TextMessage, []byte(msg.text))
			if err != nil {
				log.Printf("error while writing message: %v\n", err)
				delete(room.conns, userid)
				if len(room.conns) == 0 {
					ws.Close()
					close(channel)
					delete(rooms, roomCode)
					delete(channels, roomCode)
					break
				}
			}
		}
	}
}