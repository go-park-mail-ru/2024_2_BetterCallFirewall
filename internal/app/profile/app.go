package profile

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/profile_api"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/models"
	profileRepository "github.com/2024_2_BetterCallFirewall/internal/profile/repository"
	profileService "github.com/2024_2_BetterCallFirewall/internal/profile/service"
	"github.com/2024_2_BetterCallFirewall/pkg/start_postgres"
)

type profileManager interface {
	GetHeader(ctx context.Context, userID uint32) (*models.Header, error)
	GetFriendsID(ctx context.Context, userID uint32) ([]uint32, error)
	Create(ctx context.Context, user *models.User) (uint32, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

func getGRPC(profile profileManager) *grpc.Server {
	server := grpc.NewServer()
	profile_api.RegisterProfileServiceServer(server, profile_api.New(profile))
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

	profRepo := profileRepository.NewProfileRepo(postgresDB)
	profService := profileService.NewProfileHelper(profRepo)

	serv := getGRPC(profService)
	return serv, nil
}
