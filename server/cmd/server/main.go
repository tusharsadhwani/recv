package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rhnvrm/simples3"
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
	http.HandleFunc("/upload", GetPresignedURL)
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

func GetPresignedURL(w http.ResponseWriter, req *http.Request) {
	url, err := GeneratePresignedURL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Write([]byte(url))
}

// GeneratePresignedURL generates a pre-signed url for file upload on S3
func GeneratePresignedURL() (string, error) {
	cfg := GetConfig()

	s3 := simples3.New(cfg.S3Region, cfg.S3AccessKey, cfg.S3SecretKey)

	timestamp := time.Now()
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	randomHex := strings.ReplaceAll(uuid.String(), "-", "")
	objectKey := fmt.Sprintf("%x-%s", timestamp.Unix(), randomHex)

	url := s3.GeneratePresignedURL(simples3.PresignedInput{
		Bucket:        cfg.S3Bucket,
		ObjectKey:     objectKey,
		Method:        "PUT",
		Timestamp:     timestamp,
		ExpirySeconds: 1000,
	})
	return url, nil
}
