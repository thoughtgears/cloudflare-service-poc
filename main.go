package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/thoughtgears/cloudflare-tunnels-poc/config"
	_ "github.com/thoughtgears/cloudflare-tunnels-poc/docs"
	"github.com/thoughtgears/cloudflare-tunnels-poc/router"
	"github.com/thoughtgears/cloudflare-tunnels-poc/services"
)

// cfg holds the application's configuration, loaded from environment variables
// during initialization via the init() function.
var cfg config.Config

// init performs initial setup before the main function runs.
// It loads configuration from environment variables into the global cfg variable
// using envconfig and initializes the global zerolog logger settings.
func init() {
	envconfig.MustProcess("", &cfg)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.LevelFieldName = "severity"

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

// @title			User Service
// @version		1.0
// @description	This is a sample server for managing users.
// @termsOfService	http://swagger.io/terms/

// @contact.name	API Support
// @contact.url	http://www.example.com/support
// @contact.email	support@example.com

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/
// @schemes	https
func main() {
	// --- Dependency Initialization ---
	userService := services.NewUserService()

	// --- Router Setup ---
	routerEngine := router.NewRouter(cfg, userService)

	// --- Init swagger Paths ---
	routerEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- Server Configuration ---
	// Listen on all interfaces in production/default, only localhost in debug
	host := "0.0.0.0"
	if cfg.Debug {
		host = "127.0.0.1"
	}

	// --- Start Server ---
	log.Info().Msgf("Starting server, listening on %s:%s", host, cfg.Port)
	// routerEngine.Run blocks until the server is stopped or fails.
	// log.Fatal will print the error and exit the application if Run returns an error.
	log.Fatal().Err(routerEngine.Run(host + ":" + cfg.Port))
}
