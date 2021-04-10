package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tusharsadhwani/recv"
)

func main() {
	RunServer()
}

func RunServer() {
	http.HandleFunc("/connect", recv.HandleConnect)
	http.HandleFunc("/ws", recv.HandleWebsockets)
	http.Handle("/", http.FileServer(http.Dir("../web")))

	fmt.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
