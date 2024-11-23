package csat

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/CSAT/controller"
	"github.com/2024_2_BetterCallFirewall/internal/CSAT/repository"
	"github.com/2024_2_BetterCallFirewall/internal/CSAT/service"
	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/csat_api"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/auth"
	"github.com/2024_2_BetterCallFirewall/internal/router"
	"github.com/2024_2_BetterCallFirewall/internal/router/csat"
	"github.com/2024_2_BetterCallFirewall/pkg/start_postgres"
)

type Service interface {
	NewLike(id uint32)
	NewFriend(id uint32)
	NewMessage(id uint32)
	TimeSpent(id uint32, dur time.Duration)
}

func GetServers(cfg *config.Config) (*http.Server, *grpc.Server, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.CSATDB.Host,
		cfg.CSATDB.Port,
		cfg.CSATDB.User,
		cfg.CSATDB.Pass,
		cfg.CSATDB.DBName,
		cfg.CSATDB.SSLMode,
	)

	DB, err := start_postgres.StartPostgres(connStr, logger)
	if err != nil {
		return nil, nil, err
	}

	responder := router.NewResponder(logger)

	repo := repository.NewCSATRepository(DB)
	serv := service.NewCSATServiceImpl(repo)
	grpcServer := getGRPC(serv)
	provider, err := auth.GetAuthProvider(cfg.AUTHGRPC.Host, cfg.AUTHGRPC.Port)
	if err != nil {
		return nil, nil, err
	}

	sm := auth.New(provider)
	control := controller.NewCSATController(serv, responder)
	rout := csat.NewRouter(control, sm, logger)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.CSAT.Port),
		Handler:      rout,
		ReadTimeout:  cfg.POST.ReadTimeout,
		WriteTimeout: cfg.POST.WriteTimeout,
	}
	return httpServer, grpcServer, nil
}

func getGRPC(serv Service) *grpc.Server {
	server := grpc.NewServer()
	csat_api.RegisterCsatServiceServer(server, csat_api.NewAdapter(serv))
	return server
}
