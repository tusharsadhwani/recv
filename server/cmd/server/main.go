package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	filenameBytes, err := ioutil.ReadAll(req.Body)
	if err != nil || len(filenameBytes) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide a file name in request body."))
		return
	}

	filename := string(filenameBytes)
	url, err := GeneratePresignedURL(filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Write([]byte(url))
}

// GeneratePresignedURL generates a pre-signed url for file upload on S3
func GeneratePresignedURL(filename string) (string, error) {
	cfg := GetConfig()

	s3 := simples3.New(cfg.S3Region, cfg.S3AccessKey, cfg.S3SecretKey)

	timestamp := time.Now()
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	randomHex := strings.ReplaceAll(uuid.String(), "-", "")
	objectKey := fmt.Sprintf("%x-%s-%s", timestamp.Unix(), randomHex, filename)

	url := s3.GeneratePresignedURL(simples3.PresignedInput{
		Bucket:        cfg.S3Bucket,
		ObjectKey:     objectKey,
		Method:        "PUT",
		Timestamp:     timestamp,
		ExpirySeconds: 1000,
	})
	return url, nil
}
