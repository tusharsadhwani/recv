package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/tusharsadhwani/recv"
)

func main() {
	RunServer()
}

func RunServer() {
	port := getPort()

	http.HandleFunc("/connect", recv.HandleConnect)
	http.HandleFunc("/echo", recv.WebsocketEcho)
	http.HandleFunc("/ws", recv.HandleWebsockets)
	http.Handle("/", http.FileServer(http.Dir("../web")))

	fmt.Printf("http server started on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getPort() int {
	defaultPort := 8000

	if os.Getenv("APP_ENV") == "dev" {
		portStr, exists := os.LookupEnv("PORT")
		if !exists {
			return defaultPort
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatalln("error while reading PORT:", err)
		}
		return port
	}

	var port int
	flag.IntVar(&port, "p", defaultPort, "server port")
	flag.Parse()
	return port
}
