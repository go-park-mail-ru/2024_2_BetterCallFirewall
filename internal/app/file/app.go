package file

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/auth"
	filecontrol "github.com/2024_2_BetterCallFirewall/internal/fileService/controller"
	fileservis "github.com/2024_2_BetterCallFirewall/internal/fileService/service"
	"github.com/2024_2_BetterCallFirewall/internal/router"
	"github.com/2024_2_BetterCallFirewall/internal/router/file"
)

func GetServer(cfg *config.Config) (*http.Server, error) {
	provider, err := auth.GetAuthProvider(string(cfg.AUTHGRPC))
	if err != nil {
		return nil, err
	}
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

	sm := auth.New(provider)

	rout := file.NewRouter(fileController, sm, logger)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.FILE.Port),
		Handler:      rout,
		ReadTimeout:  cfg.SERVER.ReadTimeout,
		WriteTimeout: cfg.SERVER.WriteTimeout,
	}

	return server, nil
}
