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
	_DOMAIN := "recv.live" // Your website
	_SECURE := true        // Whether the website uses https

	isDev := os.Getenv("APP_ENV") == "dev"

	if isDev {
		roomCode := os.Getenv("ROOM")
		port := os.Getenv("PORT")
		if port == "" {
			port = "8000"
		}
		config = &Config{
			roomCode: roomCode,
			domain:   "localhost:" + port,
			scheme:   "http",
			wsscheme: "ws",
		}
	} else {
		flag.Parse()
		roomCode := flag.Arg(0)
		var scheme string
		var wsscheme string
		if _SECURE {
			scheme = "https"
			wsscheme = "wss"
		} else {
			scheme = "http"
			wsscheme = "ws"
		}

		config = &Config{
			roomCode: roomCode,
			domain:   _DOMAIN,
			scheme:   scheme,
			wsscheme: wsscheme,
		}
	}
}
