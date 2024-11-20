package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DBConnect struct {
	Port    string
	Host    string
	DBName  string
	User    string
	Pass    string
	SSLMode string
}

type Redis struct {
	Host      string
	Port      string
	MaxIdle   int
	MaxActive int
}

type Server struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type GRPCServer struct {
	Port string
	Host string
}

type Config struct {
	DB            DBConnect
	REDIS         Redis
	SERVER        Server
	AUTH          Server
	FILE          Server
	CHAT          Server
	AUTHGRPC      GRPCServer
	PROFILEGRPC   GRPCServer
	POSTGRPC      GRPCServer
	COMMUNITY     Server
	COMMUNITYGRPC GRPCServer
}

func GetConfig(configFilePath string) (*Config, error) {
	err := godotenv.Load(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	return &Config{
			DB: DBConnect{
				Port:    os.Getenv("DB_PORT"),
				Host:    os.Getenv("DB_HOST"),
				User:    os.Getenv("DB_USER"),
				Pass:    os.Getenv("DB_PASSWORD"),
				DBName:  os.Getenv("DB_NAME"),
				SSLMode: os.Getenv("DB_SSLMODE"),
			},
			REDIS: Redis{
				Host:      os.Getenv("REDIS_HOST"),
				Port:      os.Getenv("REDIS_PORT"),
				MaxIdle:   getIntEnv("REDIS_MAX_IDLE"),
				MaxActive: getIntEnv("REDIS_MAX_ACTIVE"),
			},
			SERVER: Server{
				Port:         os.Getenv("HTTP_PORT"),
				ReadTimeout:  time.Duration(getIntEnv("SERVER_READ_TIMEOUT")) * time.Second,
				WriteTimeout: time.Duration(getIntEnv("SERVER_WRITE_TIMEOUT")) * time.Second,
			},
			AUTH: Server{
				Port:         os.Getenv("AUTH_HTTP_PORT"),
				ReadTimeout:  time.Duration(getIntEnv("SERVER_READ_TIMEOUT")) * time.Second,
				WriteTimeout: time.Duration(getIntEnv("SERVER_WRITE_TIMEOUT")) * time.Second,
			},
			AUTHGRPC: GRPCServer{
				Port: os.Getenv("AUTH_GRPC_PORT"),
				Host: os.Getenv("AUTH_GRPC_HOST"),
			},
			FILE: Server{
				Port:         os.Getenv("FILE_HTTP_PORT"),
				ReadTimeout:  time.Duration(getIntEnv("SERVER_READ_TIMEOUT")) * time.Second,
				WriteTimeout: time.Duration(getIntEnv("SERVER_WRITE_TIMEOUT")) * time.Second,
			},
			CHAT: Server{
				Port:         os.Getenv("CHAT_HTTP_PORT"),
				ReadTimeout:  time.Duration(getIntEnv("SERVER_READ_TIMEOUT")) * time.Second,
				WriteTimeout: time.Duration(getIntEnv("SERVER_WRITE_TIMEOUT")) * time.Second,
			},
			PROFILEGRPC: GRPCServer{
				Port: os.Getenv("PROFILE_GRPC_PORT"),
				Host: os.Getenv("PROFILE_GRPC_HOST"),
			},
			POSTGRPC: GRPCServer{
				Port: os.Getenv("POST_GRPC_PORT"),
				Host: os.Getenv("POST_GRPC_HOST"),
			},
			COMMUNITY: Server{
				Port:         os.Getenv("COMMUNITY_HTTP_PORT"),
				ReadTimeout:  time.Duration(getIntEnv("SERVER_READ_TIMEOUT")) * time.Second,
				WriteTimeout: time.Duration(getIntEnv("SERVER_WRITE_TIMEOUT")) * time.Second,
			},
			COMMUNITYGRPC: GRPCServer{
				Port: os.Getenv("COMMUNITY_GRPC_PORT"),
				Host: os.Getenv("COMMUNITY_GRPC_HOST"),
			},
		},
		nil
}

func getIntEnv(key string) int {
	c := os.Getenv(key)
	res, err := strconv.Atoi(c)
	if err != nil {
		panic("Invalid data in key: " + key)
	}
	return res
}
