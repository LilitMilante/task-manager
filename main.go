package main

import (
	"database/sql"
	"fmt"
	"log"

	"task-manager/internal/api"
	"task-manager/internal/app"
	"task-manager/internal/repository"
	"task-manager/internal/service"

	"go.uber.org/zap"
)

func main() {
	baseLogger, _ := zap.NewDevelopment()

	defer func() {
		err := baseLogger.Sync()
		if err != nil {
			log.Println(err)
		}
	}()

	l := baseLogger.Sugar()

	cfg, err := app.NewConfig()
	if err != nil {
		l.Panicf("get config: %s", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresName)

	db, err := app.ConnectToPostgres(dsn)
	if err != nil {
		l.Panicf("connect to postgres: %s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			l.Warnw("close DB", zap.Error(err))
		}
	}(db)

	repo := repository.NewRepository(l, db)
	auth := service.NewAuthService(repo)
	taskService := service.NewTaskService(repo, auth)
	authService := service.NewAuthService(repo)

	taskHandler := api.NewTaskHandler(l, taskService)
	authHandler := api.NewAuthHandler(l, authService)

	server := api.NewServer(l, cfg.Port, taskHandler, authHandler)
	l.Infof("server started on %s", cfg.Port)

	err = server.Start()
	if err != nil {
		l.Fatal("start server error")
	}
}
