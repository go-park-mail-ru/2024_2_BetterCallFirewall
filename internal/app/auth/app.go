package auth

import (
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/auth_api"
	"github.com/2024_2_BetterCallFirewall/internal/auth/controller"
	redismy "github.com/2024_2_BetterCallFirewall/internal/auth/repository/redis"
	"github.com/2024_2_BetterCallFirewall/internal/auth/service"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/profile"
	metrics "github.com/2024_2_BetterCallFirewall/internal/metrics"
	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/router"
	"github.com/2024_2_BetterCallFirewall/internal/router/auth"
)

type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

func GetHTTPServer(cfg *config.Config) (*http.Server, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	}

	redisPool := &redis.Pool{
		MaxIdle:   cfg.REDIS.MaxIdle,
		MaxActive: cfg.REDIS.MaxActive,
		Dial: func() (redis.Conn, error) {
			addr := fmt.Sprintf("%s:%s", cfg.REDIS.Host, cfg.REDIS.Port)
			return redis.Dial("tcp", addr)
		},
	}

	profileProvider, err := ext_grpc.GetGRPCProvider(cfg.PROFILEGRPC.Host, cfg.PROFILEGRPC.Port)
	if err != nil {
		return nil, err
	}
	authMetrics, err := metrics.NewHTTPMetrics("auth")
	if err != nil {
		return nil, err
	}
	prof := profile.New(profileProvider)

	authServ := service.NewAuthServiceImpl(prof)
	responder := router.NewResponder(logger)
	sessionRepo := redismy.NewSessionRedisRepository(redisPool)
	sessionManager := service.NewSessionManager(sessionRepo)
	control := controller.NewAuthController(responder, authServ, sessionManager)

	rout := auth.NewRouter(control, sessionManager, logger, authMetrics)

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AUTH.Port),
		Handler:      rout,
		ReadTimeout:  cfg.AUTH.ReadTimeout,
		WriteTimeout: cfg.AUTH.WriteTimeout,
	}

	return &server, nil
}

func getGRPC(auth SessionManager, metr *middleware.GrpcMiddleware) *grpc.Server {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(metr.GrpcMetricsInterceptor))
	auth_api.RegisterAuthServiceServer(server, auth_api.New(auth))
	return server
}

func GetGRPCServer(cfg *config.Config) (*grpc.Server, error) {
	redisPool := &redis.Pool{
		MaxIdle:   cfg.REDIS.MaxIdle,
		MaxActive: cfg.REDIS.MaxActive,
		Dial: func() (redis.Conn, error) {
			addr := fmt.Sprintf("%s:%s", cfg.REDIS.Host, cfg.REDIS.Port)
			return redis.Dial("tcp", addr)
		},
	}
	sessionRepo := redismy.NewSessionRedisRepository(redisPool)
	sessionManager := service.NewSessionManager(sessionRepo)

	grpcMetrics, err := metrics.NewGrpcMetrics("auth")
	if err != nil {
		return nil, err
	}
	metricsmw := middleware.NewGrpcMiddleware(grpcMetrics)
	grpcServer := getGRPC(sessionManager, metricsmw)

	return grpcServer, nil
}
