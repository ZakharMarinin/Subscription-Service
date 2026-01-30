package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	HttpServer HttpServer `yaml:"http_server"`
	Storage    Storage    `yaml:"storage"`
}

type Storage struct {
	Addr string `yaml:"addr" env-default:":8081"`
}

type HttpServer struct {
	Addr        string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func MustLoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, assuming variables are set in environment")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	var cfg Config
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	cfg.Storage.Addr = os.Getenv("POSTGRES_URL")

	return &cfg
}
