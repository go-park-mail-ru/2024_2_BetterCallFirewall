package stickers

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/auth"
	"github.com/2024_2_BetterCallFirewall/internal/router"
	"github.com/2024_2_BetterCallFirewall/internal/router/stickers"
	controller "github.com/2024_2_BetterCallFirewall/internal/stickers/controller"
	repository "github.com/2024_2_BetterCallFirewall/internal/stickers/repository"
	service "github.com/2024_2_BetterCallFirewall/internal/stickers/service"
	"github.com/2024_2_BetterCallFirewall/pkg/start_postgres"
)

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
	provider, err := ext_grpc.GetGRPCProvider(cfg.AUTHGRPC.Host, cfg.AUTHGRPC.Port)
	if err != nil {
		return nil, err
	}
	sm := auth.New(provider)

	repo := repository.NewStickerRepo(postgresDB)
	stickerService := service.NewStickerUsecase(repo)
	stickerController := controller.NewStickerController(stickerService, responder)

	rout := stickers.NewRouter(stickerController, sm, logger)
	server := &http.Server{
		Handler:      rout,
		Addr:         fmt.Sprintf(":%s", cfg.STCIKER.Port),
		ReadTimeout:  cfg.STCIKER.ReadTimeout,
		WriteTimeout: cfg.STCIKER.WriteTimeout,
	}

	return server, nil
}
