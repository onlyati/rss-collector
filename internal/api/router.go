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
	//
	// ===> Load configurations
	//

	// Read configuration from YAML, this hold every config that is requires in next
	config, err := newAPIConfigFromYAML(configYAML)
	if err != nil {
		return nil, err
	}

	// Try to connect for database, if failed during initialization, then exit
	db, err := connectToDatabase(config)
	if err != nil {
		return nil, err
	}

	// This setup authentication: configuration URL then save the endpoints
	// It also download the RSA public keys
	authConf, err := auth.NewAuthentication(config.ApiOptions.AuthConfig)
	if err != nil {
		return nil, err
	}

	//
	// ===> Inject data to endpoints
	//

	// Data is injected into every single endpoint via this struct because
	// endpoints are method of this structure
	app := routes.App{
		Db:          db,
		Hostname:    config.ApiOptions.Hostname,
		Port:        config.ApiOptions.Port,
		AuthOptions: authConf.Links,
	}

	// Create a new router
	router := gin.Default()

	//
	// ===> CORS policy
	//

	// Setup CORS policy based on what is specified in the configuration file
	corsPolicy := cors.DefaultConfig()
	corsPolicy.AllowMethods = strings.Split(config.CorsConfig.Methods, ",")
	if config.CorsConfig.Origins == "*" {
		corsPolicy.AllowAllOrigins = true
	} else {
		corsPolicy.AllowOrigins = strings.Split(config.CorsConfig.Origins, ",")
	}
	router.Use(cors.New(corsPolicy))

	//
	// ===> Swagger UI
	//

	// Swagger UI is provided but not in production version
	router.LoadHTMLGlob("openapi/*")
	swagger := router.Group("/docs")
	if os.Getenv("GIN_MODE") != "release" {
		swagger.GET("/", app.GetSwaggerUI)
		swagger.GET("/index.html", app.GetSwaggerUI)
		swagger.GET("oauth2-redirect.html", app.GetRedirect)
	}
	// The YAML file always available independent from the version
	swagger.GET("/openapi.yaml", app.GetSwaggerYAML)

	//
	// ===> Endpoint /rss
	//

	apiRSS := router.Group("/rss")
	apiRSS.Use(auth.AuthMiddleware(authConf))

	// ===> Endpoints /rss/v1
	apiRSSv1 := apiRSS.Group("/v1")
	apiRSSv1.GET("", app.GetRSS)
	apiRSSv1.GET("/item", app.GetItem)
	apiRSSv1.GET("/item-category", app.GetCategories)

	//
	// ===> Endpoints /user
	//
	apiUser := router.Group("/user")
	apiUser.Use(auth.AuthMiddleware(authConf))

	// ===> Endpoints /user/v1
	apiUserV1 := apiUser.Group("/v1")

	// User stuff
	apiUserV1.GET("", app.ListUsers)
	apiUserV1.POST("", app.CreateUser)
	apiUserV1.DELETE("", app.DeleteUser)
	apiUserV1.DELETE("/hard", app.DeleteUserHard)

	// Category stuff
	apiUserV1.GET("/favorite", app.ListCategories)
	apiUserV1.POST("/favorite", app.CreateCategory)
	apiUserV1.DELETE("/favorite", app.DeleteCategory)
	apiUserV1.DELETE("/favorite/hard", app.DeleteCategoryHard)

	//
	// ===> Finish it!
	//
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
