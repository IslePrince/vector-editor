package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         int
	DataDir      string
	InkscapePath string
	MaxWorkers   int
}

func Load() *Config {
	return &Config{
		Port:         envInt("VEC_PORT", 8092),
		DataDir:      envStr("VEC_DATA_DIR", "./data"),
		InkscapePath: envStr("VEC_INKSCAPE_PATH", "inkscape"),
		MaxWorkers:   envInt("VEC_MAX_WORKERS", 4),
	}
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
