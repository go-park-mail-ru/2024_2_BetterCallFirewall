package csat

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/csat_api"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/pkg/start_postgres"
)

type Service interface {
	NewLike(id uint32)
	NewFriend(id uint32)
	NewMessage(id uint32)
}

func GetHTTPServer(cfg *config.Config) (*http.Server, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
	)

	DB, err := start_postgres.StartPostgres(connStr, logger)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func getGRPC(serv Service) *grpc.Server {
	server := grpc.NewServer()
	csat_api.RegisterCsatServiceServer(server, csat_api.NewAdapter(serv))
	return server
}

func GetGRPCServer(cfg *config.Config) (*grpc.Server, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	}

	repo := NewRepo(DB)
	serv := NewService(repo)
	grpcServer := getGRPC(serv)

	return grpcServer, nil
}
