package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

var scheme = "http"
var wsscheme = "ws"
var port = 8000
var domain = fmt.Sprintf("localhost:%d", port)

func main() {
	fmt.Println("Connecting to recv.live...")

	flag.Parse()
	arg := flag.Arg(0)

	var roomCode int
	if arg == "" {
		roomCode = createRoom()
	} else {
		var err error
		roomCode, err = strconv.Atoi(arg)
		if err != nil {
			fmt.Println("Provide a 5 digit room code")
			return
		}
		if len(arg) != 5 {
			fmt.Println("Provide a 5 digit room code")
			return
		}
	}

	conn := connect(roomCode)

	fmt.Println("Your Room code is:", roomCode)
	go readMessages(conn)

	input := bufio.NewReader(os.Stdin)
	for {
		text, _ := input.ReadString('\n')
		err := conn.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Fatal("error while writing to websocket:", err)
		}
	}
}

func createRoom() int {
	url := fmt.Sprintf("%s://%s/connect", scheme, domain)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("error while getting room code:", err)
	}
	roomCodeBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error while reading room code:", err)
	}
	roomCode, err := strconv.Atoi(string(roomCodeBytes))
	if err != nil {
		log.Fatal("error while decoding room code:", err)
	}

	return roomCode
}

func connect(roomCode int) *websocket.Conn {
	url := fmt.Sprintf("%s://%s/ws?code=%d", wsscheme, domain, roomCode)
	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		log.Fatal("error while connecting to websocket:", err)
	}

	return conn
}

func readMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("error while reading from websocket:", err)
		}
		fmt.Printf("%s", msg)
	}
}
