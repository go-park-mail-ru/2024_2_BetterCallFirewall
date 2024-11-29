package post

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/post_api"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/auth"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/community"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/profile"
	"github.com/2024_2_BetterCallFirewall/internal/metrics"
	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/post/controller"
	"github.com/2024_2_BetterCallFirewall/internal/post/repository/postgres"
	"github.com/2024_2_BetterCallFirewall/internal/post/service"
	"github.com/2024_2_BetterCallFirewall/internal/router"
	"github.com/2024_2_BetterCallFirewall/internal/router/post"
	"github.com/2024_2_BetterCallFirewall/pkg/start_postgres"
)

type postManager interface {
	GetAuthorsPosts(ctx context.Context, header *models.Header, userID uint32) ([]*models.Post, error)
}

func GetHTTPServer(cfg *config.Config, postMetric *metrics.HttpMetrics) (*http.Server, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.POSTDB.Host,
		cfg.POSTDB.Port,
		cfg.POSTDB.User,
		cfg.POSTDB.Pass,
		cfg.POSTDB.DBName,
		cfg.POSTDB.SSLMode,
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

	repo := postgres.NewAdapter(postgresDB)
	profileProvider, err := ext_grpc.GetGRPCProvider(cfg.PROFILEGRPC.Host, cfg.PROFILEGRPC.Port)
	if err != nil {
		return nil, err
	}
	pp := profile.New(profileProvider)
	communityProvider, err := ext_grpc.GetGRPCProvider(cfg.COMMUNITYGRPC.Host, cfg.COMMUNITYGRPC.Port)
	if err != nil {
		return nil, err
	}
	cp := community.New(communityProvider)

	postService := service.NewPostServiceImpl(repo, pp, cp)
	postController := controller.NewPostController(postService, responder)

	rout := post.NewRouter(postController, sm, logger, postMetric)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.POST.Port),
		Handler:      rout,
		ReadTimeout:  cfg.POST.ReadTimeout,
		WriteTimeout: cfg.POST.WriteTimeout,
	}

	return server, nil
}

func getGRPC(post postManager, metr *middleware.GrpcMiddleware) *grpc.Server {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(metr.GrpcMetricsInterceptor))
	post_api.RegisterPostServiceServer(server, post_api.New(post))
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
		cfg.POSTDB.Host,
		cfg.POSTDB.Port,
		cfg.POSTDB.User,
		cfg.POSTDB.Pass,
		cfg.POSTDB.DBName,
		cfg.POSTDB.SSLMode,
	)

	postgresDB, err := start_postgres.StartPostgres(connStr, logger)
	if err != nil {
		return nil, err
	}

	repo := postgres.NewAdapter(postgresDB)
	postHelper := service.NewPostProfileImpl(repo)

	metricsmw := middleware.NewGrpcMiddleware(grpcMetrics)
	serv := getGRPC(postHelper, metricsmw)
	return serv, nil
}
