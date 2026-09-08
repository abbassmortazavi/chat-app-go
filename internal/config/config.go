package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// go run ./cmd/api -config ./config/dev.env
type HttpServer struct {
	HttpAddress string `env:"HTTP_ADDRESS" envDefault:":8084"`
}
type Config struct {
	ENV        string `env:"ENV" envDefault:"development"`
	DBPath     string `env:"DB_PATH" envDefault:"sqlite/development"`
	DBName     string `env:"DB_NAME" envDefault:"chat-app"`
	HttpServer HttpServer
	JwtKey     string `env:"JWT_KEY" envDefault:"supersecret"`
}

func LoadConfig() *Config {
	cfg := &Config{}
	var envPath string
	flag.StringVar(&envPath, "config", "", "set the path to the environment file")
	flag.Parse()
	if envPath == "" {
		envPath = os.Getenv("CONFIG_PATH")
	}
	if envPath == "" {
		envPath = "config/dev.env"
	}
	err := cleanenv.ReadConfig(envPath, cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
