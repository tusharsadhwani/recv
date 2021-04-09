package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/tusharsadhwani/recv"
)

func main() {
	flag.Parse()
	arg := flag.Arg(0)

	if arg == "" {
		go recv.RunServer()
		//TODO: send text to channel here
		//TODO: generate a room code and send messages to that client only
	}

	roomCode, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Println("Provide a 5 digit room code")
		os.Exit(1)
	}
	if len(arg) != 5 {
		fmt.Println("Provide a 5 digit room code")
		os.Exit(1)
	}

	Receive(roomCode)
}

// Receive receives messages sent on the server
func Receive(roomCode int) {
	//TODO: implement something
}
