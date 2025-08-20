package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	JWTSecret []byte
}

var config *Config

func init() {
	godotenv.Load()
	config = &Config{
		DBUrl:     os.Getenv("DATABASE_URL"),
		JWTSecret: []byte(os.Getenv("JWT_SECRET")),
	}
}

func GetConfig() *Config {
	return config
}
