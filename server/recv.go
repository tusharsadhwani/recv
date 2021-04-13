package recv

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/tusharsadhwani/recv/utils"
)

// TODO: Delete unused rooms every few minutes,
// by adding a "last joined" timestamp, and deleting a room if it's empty
// and no-one has joined in the last X minutes

// TODO: figure out where to add mutex locks for the rooms, or remove mutex entirely

type Room struct {
	sync.Mutex
	conns    map[int]*websocket.Conn
	messages [][]byte
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
		return true
	},
}

func HandleConnect(w http.ResponseWriter, r *http.Request) {
	for {
		roomCode := utils.GenerateRoomCode()
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
		ws.WriteMessage(websocket.TextMessage, []byte("URL parameter 'code' is missing\n"))
		return
	}
	roomCodeStr := params[0]
	roomCode, err := strconv.Atoi(roomCodeStr)
	if err != nil {
		ws.WriteMessage(
			websocket.TextMessage,
			[]byte(fmt.Sprintf("Provide a %d digit room code\n", utils.RoomCodeLength)),
		)
		return
	}
	if len(roomCodeStr) != utils.RoomCodeLength {
		ws.WriteMessage(
			websocket.TextMessage,
			[]byte(fmt.Sprintf("Provide a %d digit room code\n", utils.RoomCodeLength)),
		)
		return
	}
	if rooms[roomCode] == nil {
		ws.WriteMessage(websocket.TextMessage, []byte("This room doesn't exist\n"))
		close(*channels[roomCode])
		delete(rooms, roomCode)
		delete(channels, roomCode)
		return
	}
	if channels[roomCode] == nil {
		ws.WriteMessage(websocket.TextMessage, []byte("This room doesn't exist\n"))
		if rooms[roomCode] != nil {
			for _, conn := range rooms[roomCode].conns {
				conn.Close()
			}
		}
		return
	}

	counters[roomCode] += 1
	userid := counters[roomCode]
	room := rooms[roomCode]
	room.conns[userid] = ws

	// Send older messages to client
	for _, msg := range room.messages {
		ws.WriteMessage(websocket.TextMessage, msg)
	}

	for {
		_, bytes, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error while reading message: %v", err)
			ws.Close()
			delete(room.conns, userid)
			if len(room.conns) == 0 {
				delete(rooms, roomCode)
				close(*channels[roomCode])
				delete(channels, roomCode)
			}
			break
		}

		// TODO: Add a limit to the old message cache
		room.messages = append(room.messages, bytes)

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
