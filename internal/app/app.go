package app

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx"

	ChatController "github.com/2024_2_BetterCallFirewall/internal/chat/controller"
	chatRepository "github.com/2024_2_BetterCallFirewall/internal/chat/repository/postgres"
	chatService "github.com/2024_2_BetterCallFirewall/internal/chat/service"
	communityController "github.com/2024_2_BetterCallFirewall/internal/community/controller"
	communityRepository "github.com/2024_2_BetterCallFirewall/internal/community/repository"
	communityService "github.com/2024_2_BetterCallFirewall/internal/community/service"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/adapter/auth"
	filecontrol "github.com/2024_2_BetterCallFirewall/internal/fileService/controller"
	fileservis "github.com/2024_2_BetterCallFirewall/internal/fileService/service"
	postController "github.com/2024_2_BetterCallFirewall/internal/post/controller"
	postgresPost "github.com/2024_2_BetterCallFirewall/internal/post/repository/postgres"
	postServ "github.com/2024_2_BetterCallFirewall/internal/post/service"
	profileController "github.com/2024_2_BetterCallFirewall/internal/profile/controller"
	profileRepository "github.com/2024_2_BetterCallFirewall/internal/profile/repository"
	profileService "github.com/2024_2_BetterCallFirewall/internal/profile/service"
	"github.com/2024_2_BetterCallFirewall/internal/router"
	"github.com/2024_2_BetterCallFirewall/pkg/start_postgres"
)

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

	profileRepo := profileRepository.NewProfileRepo(postgresDB)

	responder := router.NewResponder(logger)

	postRepo := postgresPost.NewAdapter(postgresDB)
	chatRepo := chatRepository.NewChatRepository(postgresDB)

	fileServ := fileservis.NewFileService()
	fileController := filecontrol.NewFileController(fileServ, responder)

	postsHelper := postServ.NewPostProfileImpl(postRepo)
	profileUsecase := profileService.NewProfileUsecase(profileRepo, postsHelper)
	profileControl := profileController.NewProfileController(profileUsecase, responder)

	chatService := chatService.NewChatService(chatRepo)
	chatControl := ChatController.NewChatController(chatService, responder)
	defer close(chatControl.Messages)

	communityRepo := communityRepository.NewCommunityRepository(postgresDB)
	communityServ := communityService.NewService(communityRepo)
	communityControl := communityController.NewController(responder, communityServ)

	postService := postServ.NewPostServiceImpl(postRepo, profileUsecase, communityRepo)
	postControl := postController.NewPostController(postService, responder)

	provider, err := auth.GetAuthProvider(string(cfg.AUTHGRPC))
	if err != nil {
		return err
	}

	sm := auth.New(provider)
	rout := router.NewRouter(
		profileControl,
		postControl,
		fileController,
		sm,
		chatControl,
		communityControl,
		logger,
	)

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.SERVER.Port),
		Handler:      rout,
		ReadTimeout:  cfg.SERVER.ReadTimeout,
		WriteTimeout: cfg.SERVER.WriteTimeout,
	}

	logger.Infof("Starting server on port %s", cfg.SERVER.Port)
	return server.ListenAndServe()
}
