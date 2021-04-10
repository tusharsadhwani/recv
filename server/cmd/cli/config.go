package main

import (
	"flag"
	"os"
)

type Config struct {
	roomCode string
	domain   string
	scheme   string
	wsscheme string
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		InitConfig()
	}
	return config
}

func InitConfig() {
	if os.Getenv("APP_ENV") == "dev" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8000"
		}
		config = &Config{
			roomCode: os.Getenv("ROOM"),
			domain:   "localhost:" + port,
			scheme:   "http",
			wsscheme: "ws",
		}
	} else {
		flag.Parse()
		roomCode := flag.Arg(0)
		config = &Config{
			roomCode: roomCode,
			domain:   "recv.live",
			scheme:   "https",
			wsscheme: "wss",
		}
	}
}
