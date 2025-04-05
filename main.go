package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/thoughtgears/cloudflare-tunnels-poc/config"
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

// main is the entry point of the application.
// It initializes the necessary services, sets up the HTTP router and middleware,
// determines the host and port to listen on based on configuration, and starts
// the Gin HTTP server. It logs a fatal error if the server fails to start.
func main() {
	// --- Dependency Initialization ---
	userService := services.NewUserService()

	// --- Router Setup ---
	routerEngine := router.NewRouter(cfg, userService)

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
