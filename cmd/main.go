package main

import (
	"MusicLibrary_Test/config"
	_ "MusicLibrary_Test/docs"
	"MusicLibrary_Test/internal/db"
	"MusicLibrary_Test/internal/handlers"
	"MusicLibrary_Test/internal/repository"
	"MusicLibrary_Test/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"net/http"
)

// @title Music Library API
// @version 1.0
// @description API Server for Music Library

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	cfg := config.LoadConfig()

	//settings logger
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	//connect to db
	dbpool, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		logrus.Fatal("Can't connect to database: ", err)
	}
	defer dbpool.Close()

	//migrations start
	err = db.RunMigrations(cfg.DBUrl)
	if err != nil {
		logrus.Fatal("Failed to run migrations: ", err)
	}

	//bundle layers
	repo := repository.NewSongRepository(dbpool)
	apiService := service.NewAPIService()
	handler := handlers.NewSongHandler(repo, apiService)

	//routers
	router := gin.Default()
	router.GET("/song", handler.GetSongs)
	router.GET("/songs:id/text", handler.GetSongsText)
	router.POST("/songs", handler.AddSong)
	router.PUT("/songs/:id", handler.UpdateSong)
	router.DELETE("/songs/:id", handler.DeleteSong)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//start service
	srv := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}
	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatal("Failed to start service")
	}
	logrus.Infof("Server is running at %s", cfg.ServerPort)
}
