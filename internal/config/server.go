package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Host              string
	Port              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func (srvCfg *ServerConfig) Load() {
	godotenv.Load()

	srvCfg.Host = os.Getenv("SERVER_HOST")
	srvCfg.Port = getEnv("SERVER_PORT", "8080")
	srvCfg.ReadTimeout = getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second)
	srvCfg.ReadHeaderTimeout = getEnvDuration("SERVER_READ_HEADER_TIMEOUT", 5*time.Second)
	srvCfg.WriteTimeout = getEnvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second)
	srvCfg.IdleTimeout = getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second)
	srvCfg.ShutdownTimeout = getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 5*time.Second)
}

func (srvCfg *ServerConfig) Addr() string {
	return srvCfg.Host + ":" + srvCfg.Port
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
