package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BucketName string
	ObjectKey  string
	Region     string
	APIKey     string
}

func NewConfigFromEnv() (*Config, error) {
	_ = godotenv.Load()

	bucket := os.Getenv("BUCKET_NAME")
	key := os.Getenv("FILE_KEY")

	if bucket == "" || key == "" {
		return nil, fmt.Errorf("Not set env")
	}

	return &Config{
		BucketName: bucket,
		ObjectKey:  key,
		Region:     "ap-northeast-1",
	}, nil
}
