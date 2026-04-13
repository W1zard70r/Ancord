package config

import "os"

type Config struct {
	DBURL  string
	JWTKey string
	Port   string
}

func Load() *Config {
	return &Config{
		DBURL:  os.Getenv("DB_URL"),
		JWTKey: os.Getenv("JWT_KEY"),
		Port:   os.Getenv("PORT"),
	}
}
