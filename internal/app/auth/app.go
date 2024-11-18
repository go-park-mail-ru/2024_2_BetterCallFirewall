package auth

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/auth_api"
	"github.com/2024_2_BetterCallFirewall/internal/auth/controller"
	"github.com/2024_2_BetterCallFirewall/internal/auth/repository/postgres"
	redismy "github.com/2024_2_BetterCallFirewall/internal/auth/repository/redis"
	"github.com/2024_2_BetterCallFirewall/internal/auth/service"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/router"
	"github.com/2024_2_BetterCallFirewall/internal/router/auth"
	"github.com/2024_2_BetterCallFirewall/pkg/start_postgres"
)

type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

func Run() error {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	}

	confPath := flag.String("c", ".env", "path to config file")
	flag.Parse()

	cfg, err := config.GetConfig(*confPath)
	if err != nil {
		return err
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
		return err
	}
	defer postgresDB.Close()

	redisPool := &redis.Pool{
		MaxIdle:   cfg.REDIS.MaxIdle,
		MaxActive: cfg.REDIS.MaxActive,
		Dial: func() (redis.Conn, error) {
			addr := fmt.Sprintf("%s:%s", cfg.REDIS.Host, cfg.REDIS.Port)
			return redis.Dial("tcp", addr)
		},
	}

	repo := postgres.NewAdapter(postgresDB)
	authServ := service.NewAuthServiceImpl(repo)
	responder := router.NewResponder(logger)
	sessionRepo := redismy.NewSessionRedisRepository(redisPool)
	sessionManager := service.NewSessionManager(sessionRepo)
	control := controller.NewAuthController(responder, authServ, sessionManager)

	rout := auth.NewRouter(control, sessionManager, logger)

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AUTH.Port),
		Handler:      rout,
		ReadTimeout:  cfg.AUTH.ReadTimeout,
		WriteTimeout: cfg.AUTH.WriteTimeout,
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.AUTHGRPC))
	if err != nil {
		return err
	}
	go startGRPC(l, logger, sessionManager)
	logger.Infof("Listening on :%s with protocol gRPC", cfg.AUTHGRPC)

	logger.Infof("Starting server on port %s", cfg.AUTH.Port)
	return server.ListenAndServe()
}

func startGRPC(l net.Listener, logger *logrus.Logger, auth SessionManager) {
	server := grpc.NewServer()
	auth_api.RegisterAuthServiceServer(server, auth_api.New(auth))
	if err := server.Serve(l); err != nil {
		logger.Fatal(err)
	}
}
