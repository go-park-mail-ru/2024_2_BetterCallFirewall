package post

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/post_api"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/post/repository/postgres"
	"github.com/2024_2_BetterCallFirewall/internal/post/service"
	"github.com/2024_2_BetterCallFirewall/pkg/start_postgres"
)

type postManager interface {
	GetAuthorsPosts(ctx context.Context, header *models.Header) ([]*models.Post, error)
}

func getGRPC(post postManager) *grpc.Server {
	server := grpc.NewServer()
	post_api.RegisterPostServiceServer(server, post_api.New(post))
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

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
	)

	postgresDB, err := start_postgres.StartPostgres(connStr, logger)
	if err != nil {
		return nil, err
	}

	repo := postgres.NewAdapter(postgresDB)
	postHelper := service.NewPostProfileImpl(repo)

	serv := getGRPC(postHelper)
	return serv, nil
}
