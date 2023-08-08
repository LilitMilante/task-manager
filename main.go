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

const (
	port     = "8080"
	dbHost   = "localhost"
	dbPort   = 5432
	user     = "postgres"
	password = "your-password"
	dbName   = "task-manager"
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

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, user, password, dbName)

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
	s := service.NewService(repo)
	h := api.NewHandler(l, s)

	server := api.NewServer(port, h)
	l.Infof("server started on %s", port)

	err = server.Start()
	if err != nil {
		l.Fatal("start server error")
	}
}
