package app

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx"
	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/auth/controller"
	"github.com/2024_2_BetterCallFirewall/internal/auth/repository/postgres"
	redismy "github.com/2024_2_BetterCallFirewall/internal/auth/repository/redis"
	"github.com/2024_2_BetterCallFirewall/internal/auth/service"
	ChatController "github.com/2024_2_BetterCallFirewall/internal/chat/controller"
	chatRepository "github.com/2024_2_BetterCallFirewall/internal/chat/repository/postgres"
	chatService "github.com/2024_2_BetterCallFirewall/internal/chat/service"
	communityController "github.com/2024_2_BetterCallFirewall/internal/community/controller"
	communityRepository "github.com/2024_2_BetterCallFirewall/internal/community/repository"
	communityService "github.com/2024_2_BetterCallFirewall/internal/community/service"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	filecontrol "github.com/2024_2_BetterCallFirewall/internal/fileService/controller"
	fileRepo "github.com/2024_2_BetterCallFirewall/internal/fileService/repository"
	fileservis "github.com/2024_2_BetterCallFirewall/internal/fileService/service"
	postController "github.com/2024_2_BetterCallFirewall/internal/post/controller"
	postgresPost "github.com/2024_2_BetterCallFirewall/internal/post/repository/postgres"
	postServ "github.com/2024_2_BetterCallFirewall/internal/post/service"
	profileController "github.com/2024_2_BetterCallFirewall/internal/profile/controller"
	profileRepository "github.com/2024_2_BetterCallFirewall/internal/profile/repository"
	profileService "github.com/2024_2_BetterCallFirewall/internal/profile/service"
	"github.com/2024_2_BetterCallFirewall/internal/router"
)

func Run() error {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	}

	confPath := flag.String("c", "./.env", "path to config file")
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

	postgresDB, err := startPostgres(connStr, logger)
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
	profileRepo := profileRepository.NewProfileRepo(postgresDB)
	authServ := service.NewAuthServiceImpl(repo)

	responder := router.NewResponder(logger)
	sessionRepo := redismy.NewSessionRedisRepository(redisPool)
	sessionManager := service.NewSessionManager(sessionRepo)
	control := controller.NewAuthController(responder, authServ, sessionManager)

	postRepo := postgresPost.NewAdapter(postgresDB)
	chatRepo := chatRepository.NewChatRepository(postgresDB)

	fileRepository := fileRepo.NewFileRepo(postgresDB)
	fileServ := fileservis.NewFileService(fileRepository)
	fileController := filecontrol.NewFileController(fileServ, responder)

	postsHelper := postServ.NewPostProfileImpl(fileServ, postRepo)
	profileUsecase := profileService.NewProfileUsecase(profileRepo, postsHelper)
	profileControl := profileController.NewProfileController(profileUsecase, fileServ, responder)

	chatService := chatService.NewChatService(chatRepo)
	chatControl := ChatController.NewChatController(chatService, responder)
	defer close(chatControl.Messages)

	communityRepo := communityRepository.NewCommunityRepository(postgresDB)
	communityServ := communityService.NewService(communityRepo)
	communityControl := communityController.NewController(responder, communityServ)

	postService := postServ.NewPostServiceImpl(postRepo, profileUsecase, communityRepo)
	postControl := postController.NewPostController(postService, responder, fileServ)

	rout := router.NewRouter(control,
		profileControl,
		postControl,
		fileController,
		sessionManager,
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

func startPostgres(connStr string, logger *logrus.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}
	db.SetMaxOpenConns(10)

	retrying := 10
	i := 1
	logger.Infof("try ping postgresql:%v", i)
	for err = db.Ping(); err != nil; err = db.Ping() {
		if i >= retrying {
			return nil, fmt.Errorf("postgres connect: %w", err)
		}
		i++
		time.Sleep(1 * time.Second)
		logger.Infof("try ping postgresql: %v", i)
	}

	return db, nil
}
