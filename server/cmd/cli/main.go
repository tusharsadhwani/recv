package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	. "github.com/tusharsadhwani/recv/constants" //lint:ignore ST1001 importing constants
)

func main() {
	cfg := GetConfig()
	fmt.Printf("Connecting to %s...\n", cfg.domain)

	var roomCode int
	if cfg.roomCode == "" {
		roomCode = createRoom()
	} else {
		var err error
		if len(cfg.roomCode) != RoomCodeLength {
			fmt.Printf("Provide a %d digit room code\n", RoomCodeLength)
			return
		}
		roomCode, err = strconv.Atoi(cfg.roomCode)
		if err != nil {
			fmt.Printf("Provide a %d digit room code\n", RoomCodeLength)
			return
		}
	}

	conn := connect(roomCode)

	fmt.Println("Your Room code is:", roomCode)
	go readMessages(conn)
	writeMessages(conn)
}

func createRoom() int {
	cfg := GetConfig()
	url := fmt.Sprintf("%s://%s/connect", cfg.scheme, cfg.domain)
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
	cfg := GetConfig()
	url := fmt.Sprintf("%s://%s/ws?code=%d", cfg.wsscheme, cfg.domain, roomCode)
	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		log.Fatal("error while connecting to websocket:", err)
	}

	return conn
}

func readMessages(conn *websocket.Conn) {
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		msg := string(msgBytes)
		if strings.HasPrefix(msg, "https://") {
			fmt.Println("Received link: " + msg)
		} else {
			fmt.Println(msg)
		}
	}
}

func writeMessages(conn *websocket.Conn) {
	input := bufio.NewReader(os.Stdin)
	for {
		text, _ := input.ReadBytes('\n')
		text = bytes.TrimRight(text, "\n")
		err := conn.WriteMessage(websocket.TextMessage, text)
		if err != nil {
			log.Fatal(err)
		}
	}
}
