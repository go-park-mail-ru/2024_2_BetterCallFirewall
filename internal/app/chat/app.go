package chat

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	ChatController "github.com/2024_2_BetterCallFirewall/internal/chat/controller"
	chatRepository "github.com/2024_2_BetterCallFirewall/internal/chat/repository/postgres"
	chatService "github.com/2024_2_BetterCallFirewall/internal/chat/service"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/auth"
	"github.com/2024_2_BetterCallFirewall/internal/metrics"
	"github.com/2024_2_BetterCallFirewall/internal/router"
	"github.com/2024_2_BetterCallFirewall/internal/router/chat"
	"github.com/2024_2_BetterCallFirewall/pkg/start_postgres"
)

func GetServer(cfg *config.Config) (*http.Server, error) {
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

	chatRepo := chatRepository.NewChatRepository(postgresDB)
	chatServ := chatService.NewChatService(chatRepo)
	chatControl := ChatController.NewChatController(chatServ, responder)
	//defer close(chatControl.Messages)

	provider, err := ext_grpc.GetGRPCProvider(cfg.AUTHGRPC.Host, cfg.AUTHGRPC.Port)
	if err != nil {
		return nil, err
	}
	sm := auth.New(provider)

	chatMetrics, err := metrics.NewHTTPMetrics("chat")
	if err != nil {
		return nil, err
	}

	rout := chat.NewRouter(chatControl, sm, logger, chatMetrics)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.CHAT.Port),
		Handler:      rout,
		ReadTimeout:  cfg.CHAT.ReadTimeout,
		WriteTimeout: cfg.CHAT.WriteTimeout,
	}

	return server, nil
}
