package file

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/auth"
	filecontrol "github.com/2024_2_BetterCallFirewall/internal/fileService/controller"
	fileservis "github.com/2024_2_BetterCallFirewall/internal/fileService/service"
	"github.com/2024_2_BetterCallFirewall/internal/metrics"
	"github.com/2024_2_BetterCallFirewall/internal/router"
	"github.com/2024_2_BetterCallFirewall/internal/router/file"
)

func GetServer(cfg *config.Config) (*http.Server, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	}

	responder := router.NewResponder(logger)
	fileServ := fileservis.NewFileService()
	fileController := filecontrol.NewFileController(fileServ, responder)

	provider, err := ext_grpc.GetGRPCProvider(cfg.AUTHGRPC.Host, cfg.AUTHGRPC.Port)
	if err != nil {
		return nil, err
	}
	sm := auth.New(provider)

	fileMetrics, err := metrics.NewFileMetrics("file")
	if err != nil {
		return nil, err
	}

	rout := file.NewRouter(fileController, sm, logger, fileMetrics)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.FILE.Port),
		Handler:      rout,
		ReadTimeout:  cfg.FILE.ReadTimeout,
		WriteTimeout: cfg.FILE.WriteTimeout,
	}

	return server, nil
}
