package profile

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/profile_api"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/auth"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/post"
	"github.com/2024_2_BetterCallFirewall/internal/metrics"
	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/profile/controller"
	"github.com/2024_2_BetterCallFirewall/internal/profile/repository"
	"github.com/2024_2_BetterCallFirewall/internal/profile/service"
	"github.com/2024_2_BetterCallFirewall/internal/router"
	"github.com/2024_2_BetterCallFirewall/internal/router/profile"
	"github.com/2024_2_BetterCallFirewall/pkg/start_postgres"
)

type profileManager interface {
	GetHeader(ctx context.Context, userID uint32) (*models.Header, error)
	GetFriendsID(ctx context.Context, userID uint32) ([]uint32, error)
	Create(ctx context.Context, user *models.User) (uint32, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

func GetHTTPServer(cfg *config.Config, metric *metrics.HttpMetrics) (*http.Server, error) {
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
	provider, err := ext_grpc.GetGRPCProvider(cfg.AUTHGRPC.Host, cfg.AUTHGRPC.Port)
	if err != nil {
		return nil, err
	}
	sm := auth.New(provider)

	postProvider, err := ext_grpc.GetGRPCProvider(cfg.POSTGRPC.Host, cfg.POSTGRPC.Port)
	if err != nil {
		return nil, err
	}
	pp := post.New(postProvider)

	repo := repository.NewProfileRepo(postgresDB)
	profileService := service.NewProfileUsecase(repo, pp)
	profileController := controller.NewProfileController(profileService, responder)

	rout := profile.NewRouter(profileController, sm, logger, metric)
	server := &http.Server{
		Handler:      rout,
		Addr:         fmt.Sprintf(":%s", cfg.PROFILE.Port),
		ReadTimeout:  cfg.PROFILE.ReadTimeout,
		WriteTimeout: cfg.PROFILE.WriteTimeout,
	}

	return server, nil
}

func getGRPC(profile profileManager, metr *middleware.GrpcMiddleware) *grpc.Server {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(metr.GrpcMetricsInterceptor))
	profile_api.RegisterProfileServiceServer(server, profile_api.New(profile))
	return server
}

func GetGRPCServer(cfg *config.Config, grpcMetrics *metrics.GrpcMetrics) (*grpc.Server, error) {
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

	profRepo := repository.NewProfileRepo(postgresDB)
	profService := service.NewProfileHelper(profRepo)

	metricsmw := middleware.NewGrpcMiddleware(grpcMetrics)
	serv := getGRPC(profService, metricsmw)
	return serv, nil
}
