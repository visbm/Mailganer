package main

import (
	"context"
	"log"
	"mailganer/internal/handler"
	"mailganer/pkg/logger"
	"mailganer/pkg/server"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file. %s", err.Error())
	}

	logger := logger.GetLogger()

	/*db, err := postgres.NewPostgresDB(&postgres.PostgresDB{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
		Logger:   logger,
	})
	if err != nil {
		logger.Panicf("Error while initialisation database:%s", err)
	}
	repository := repository.New(db, logger)*/
	handler := handler.NewHandler(logger)

	server := server.NewServer(logger, *handler, os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	idleConnsClosed := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		err = server.Shutdown(context.Background())
		if err != nil {
			logger.Errorf("Error occured on server shutting down: %s", err.Error())
		}
		//	err = db.Close()
		if err != nil {
			logger.Errorf("Error occured on db connection close: %s", err.Error())
		}

		logger.Info("shutting down")
		os.Exit(0)
	}()

	if err := server.Start(); err != http.ErrServerClosed {
		logger.Panicf("Error while starting server:%s", err)
	}
	<-idleConnsClosed
}
