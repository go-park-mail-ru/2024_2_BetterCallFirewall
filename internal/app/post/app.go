package post

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/post_api"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/auth"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/community"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/profile"
	"github.com/2024_2_BetterCallFirewall/internal/metrics"
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

	provider, err := auth.GetAuthProvider(cfg.AUTHGRPC.Host, cfg.AUTHGRPC.Port)
	if err != nil {
		return nil, err
	}
	sm := auth.New(provider)

	repo := postgres.NewAdapter(postgresDB)
	profileProvider, err := profile.GetProfileProvider(cfg.PROFILEGRPC.Host, cfg.PROFILEGRPC.Port)
	if err != nil {
		return nil, err
	}
	pp := profile.New(profileProvider)
	communityProvider, err := community.GetCommunityProvider(cfg.COMMUNITYGRPC.Host, cfg.COMMUNITYGRPC.Port)
	if err != nil {
		return nil, err
	}
	cp := community.New(communityProvider)

	postService := service.NewPostServiceImpl(repo, pp, cp)
	postController := controller.NewPostController(postService, responder)

	postMetric, err := metrics.NewHTTPMetrics("post")
	if err != nil {
		return nil, err
	}

	rout := post.NewRouter(postController, sm, logger, postMetric)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.POST.Port),
		Handler:      rout,
		ReadTimeout:  cfg.POST.ReadTimeout,
		WriteTimeout: cfg.POST.WriteTimeout,
	}

	return server, nil
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
