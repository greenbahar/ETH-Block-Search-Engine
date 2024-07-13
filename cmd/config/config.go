package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerConf    ServerConf
	EthClientConf EthClientConf
}

type ServerConf struct {
	ServerIP     string        `envconfig:"SERVER_IP"`
	ServerPort   string        `envconfig:"SERVER_PORT" default:"8000"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"5"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"5"`
}

type EthClientConf struct {
	EthereumHttpURL               string `envconfig:"HTTP_ETH_URL"`
	EthereumWSSURL                string `envconfig:"WSS_ETH_URL"`
	NumberOfRecentBlocks          int    `envconfig:"NUMBER_OF_RECENT_BLOCKS" default:"50"`
	NumberOfBlockProcessorWorkers int    `envconfig:"NUMBER_OF_BLOCK_PROCESSOR_WORKERS" default:"7"`
}

func LoadConfig(logger *log.Logger) *Config {
	return &Config{
		ServerConf: ServerConf{
			ServerIP:     getEnv("SERVER_IP", ""),
			ServerPort:   getEnv("SERVER_PORT", "8000"),
			ReadTimeout:  time.Duration(getEnvAsInt("READ_TIMEOUT", 5)) * time.Second,
			WriteTimeout: time.Duration(getEnvAsInt("WRITE_TIMEOUT", 5)) * time.Second,
		},
		EthClientConf: EthClientConf{
			EthereumHttpURL:               getEnv("HTTP_ETH_URL", ""),
			EthereumWSSURL:                getEnv("WSS_ETH_URL", ""),
			NumberOfRecentBlocks:          getEnvAsInt("NUMBER_OF_RECENT_BLOCKS", 50),
			NumberOfBlockProcessorWorkers: getEnvAsInt("NUMBER_OF_BLOCK_PROCESSOR_WORKERS", 7),
		},
	}
}

// Helper function to get environment variables with a fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper function to get an integer environment variable with a fallback
func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
