package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpPort string `env:"HTTP_PORT"`

	DBHOST     string `env:"DB_HOST"`
	DBPORT     int    `env:"DB_PORT"`
	DBUSER     string `env:"DB_USER"`
	DBPASSWORD string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_Name"`
}

func New(filePath string) (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(filePath, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cleanenv.ReadConfig %w", err)
	}
	return cfg, err
}
