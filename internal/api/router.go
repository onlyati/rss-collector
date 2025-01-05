package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/onlyati/rss-collector/internal/api/auth"
	"github.com/onlyati/rss-collector/internal/api/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type API struct {
	Db     *gorm.DB
	Router *gin.Engine
	App    *routes.App
	Config *APIConfig
}

func NewRouter(configYAML []byte) (*API, error) {
	config, err := newAPIConfigFromYAML(configYAML)
	if err != nil {
		return nil, err
	}
	db, err := connectToDatabase(config)
	if err != nil {
		return nil, err
	}

	authConf, err := auth.NewAuthentication(config.ApiOptions.AuthConfig)
	if err != nil {
		return nil, err
	}

	app := routes.App{Db: db}
	router := gin.Default()

	corsPolicy := cors.DefaultConfig()
	corsPolicy.AllowMethods = strings.Split(config.CorsConfig.Methods, ",")
	if config.CorsConfig.Origins == "*" {
		corsPolicy.AllowAllOrigins = true
	} else {
		corsPolicy.AllowOrigins = strings.Split(config.CorsConfig.Origins, ",")
	}
	router.Use(cors.New(corsPolicy))

	if os.Getenv("GIN_MODE") != "release" {
		router.StaticFile("/docs", "./openapi/index.html")
		router.StaticFile("/docs/index.html", "./openapi/index.html")
		router.StaticFile("/docs/oauth2-redirect.html", "./openapi/oauth2-redirect.html")
	}
	router.StaticFile("/docs/openapi.yaml", "./openapi/openapi.yaml")

	apiRSS := router.Group("/rss")
	apiRSS.Use(auth.AuthMiddleware(authConf))

	apiV1 := apiRSS.Group("/v1")
	apiV1.GET("", app.GetRSS)
	apiV1.GET("/item", app.GetItem)

	api := API{
		Db:     db,
		Router: router,
		App:    &app,
		Config: config,
	}
	return &api, nil
}

func (api *API) Listen() {
	addr := fmt.Sprintf("%s:%d", api.Config.ApiOptions.Hostname, api.Config.ApiOptions.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: api.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen failed", "error", err)
			panic("listen failed " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server stop failed", "error", err)
	}

	slog.Info("server exiting")
}

func connectToDatabase(config *APIConfig) (*gorm.DB, error) {
	slog.Info(
		"connect to database",
		"hostname", config.DatabaseOptions.Hostname,
		"port", config.DatabaseOptions.Port,
		"user", config.DatabaseOptions.UserName,
		"db_name", config.DatabaseOptions.DbName,
	)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.DatabaseOptions.Hostname,
		config.DatabaseOptions.UserName,
		config.DatabaseOptions.Password,
		config.DatabaseOptions.DbName,
		config.DatabaseOptions.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
