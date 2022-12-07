package main

import (
	"ToDo/database"
	"ToDo/server"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutDownTimeOut = 10 * time.Second

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//utils.FirebaseClient = config.SetupFirebase()

	// create server instance
	srv := server.SetupRoutes()
	if err := database.ConnectAndMigrate(
		"localhost",
		"5433",
		"local",
		"local",
		"local",
		database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}
	logrus.Print("migration successful!!")

	// configure firebase
	//firebaseAuth := config.SetupFirebase()
	//fbCtx := context.WithValue(r.Context(), authContext, firebaseAuth)

	go func() {
		if err := srv.Run(":8080"); err != nil && err != http.ErrServerClosed {
			logrus.Panicf("Failed to run server with error: %+v", err)
		}
	}()
	logrus.Print("Server started at :8080")

	<-done

	logrus.Info("shutting down server")
	if err := database.ShutdownDatabase(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}
	if err := srv.Shutdown(shutDownTimeOut); err != nil {
		logrus.WithError(err).Panic("failed to gracefully shutdown server")
	}
}
