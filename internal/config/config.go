package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpPort string `env:"HTTP_PORT"`

	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`

	JWTSecret string `env:"JWT_SECRET"`
}

func New(filePath string) (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(filePath, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cleanenv.ReadConfig %w", err)
	}
	return cfg, err
}
