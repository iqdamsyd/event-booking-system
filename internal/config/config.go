package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	JWTSecret []byte
	RedisUrl  string
}

var config *Config

func init() {
	godotenv.Load()
	config = &Config{
		DBUrl:     os.Getenv("DATABASE_URL"),
		JWTSecret: []byte(os.Getenv("JWT_SECRET")),
		RedisUrl:  os.Getenv("REDIS_URL"),
	}
}

func GetConfig() *Config {
	return config
}
