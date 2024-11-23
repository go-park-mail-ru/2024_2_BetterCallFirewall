package community

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/community_api"
	communityController "github.com/2024_2_BetterCallFirewall/internal/community/controller"
	communityRepository "github.com/2024_2_BetterCallFirewall/internal/community/repository"
	communityService "github.com/2024_2_BetterCallFirewall/internal/community/service"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/auth"
	"github.com/2024_2_BetterCallFirewall/internal/metrics"
	"github.com/2024_2_BetterCallFirewall/internal/router"
	"github.com/2024_2_BetterCallFirewall/internal/router/community"
	"github.com/2024_2_BetterCallFirewall/pkg/start_postgres"
)

type communityManager interface {
	CheckAccess(ctx context.Context, communityID, userID uint32) bool
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

	postgresDB, err := start_postgres.StartPostgres(connStr, logger)
	if err != nil {
		return nil, err
	}

	responder := router.NewResponder(logger)

	communityRepo := communityRepository.NewCommunityRepository(postgresDB)
	communityServ := communityService.NewCommunityService(communityRepo)
	communityControl := communityController.NewCommunityController(responder, communityServ)

	provider, err := auth.GetAuthProvider(cfg.AUTHGRPC.Host, cfg.AUTHGRPC.Port)
	if err != nil {
		return nil, err
	}
	sm := auth.New(provider)

	communityMetrics, err := metrics.NewHTTPMetrics("community")
	if err != nil {
		return nil, err
	}

	rout := community.NewRouter(communityControl, sm, logger, communityMetrics)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.COMMUNITY.Port),
		Handler:      rout,
		ReadTimeout:  cfg.COMMUNITY.ReadTimeout,
		WriteTimeout: cfg.COMMUNITY.WriteTimeout,
	}

	return server, nil
}

func getGRPC(community communityManager) *grpc.Server {
	server := grpc.NewServer()
	community_api.RegisterCommunityServiceServer(server, community_api.New(community))
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

	repo := communityRepository.NewCommunityRepository(postgresDB)
	communityHelper := communityService.NewServiceHelper(repo)

	serv := getGRPC(communityHelper)
	return serv, nil
}
