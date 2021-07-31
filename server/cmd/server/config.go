package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	S3Bucket    string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		InitConfig()
	}

	return config
}

func InitConfig() {
	dotEnvPath, err := filepath.Abs("./.env")
	if err != nil {
		log.Fatalln("Error reading .env file")
	}

	err = godotenv.Load(dotEnvPath)
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	S3Bucket := os.Getenv("S3_BUCKET")
	S3Region := os.Getenv("S3_REGION")
	S3AccessKey := os.Getenv("S3_ACCESS_KEY")
	S3SecretKey := os.Getenv("S3_SECRET_KEY")
	config = &Config{
		S3Bucket:    S3Bucket,
		S3Region:    S3Region,
		S3AccessKey: S3AccessKey,
		S3SecretKey: S3SecretKey,
	}
}
