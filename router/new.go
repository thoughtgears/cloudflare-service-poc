package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thoughtgears/cloudflare-tunnels-poc/config"
	"github.com/thoughtgears/cloudflare-tunnels-poc/handlers"
	"github.com/thoughtgears/cloudflare-tunnels-poc/router/middleware"
	"github.com/thoughtgears/cloudflare-tunnels-poc/services"
)

// NewRouter creates and configures a new Gin engine instance with middleware and API routes.
//
// It initializes a new Gin engine using gin.New() (instead of gin.Default() to allow
// for explicit middleware selection). It sets the Gin mode to ReleaseMode if
// config.Debug is false.
//
// Middleware added includes:
//   - A custom structured logger (via middleware.Logger()).
//   - Gin's default recovery middleware to handle panics gracefully.
//
// It clears any default trusted proxies using SetTrustedProxies(nil), which is often
// suitable when running behind a known reverse proxy or load balancer.
//
// Routes defined:
//   - GET /health: A simple health check endpoint.
//   - /users group: CRUD endpoints for user management, handled by the UserHandler.
//   - GET /: Retrieves all users.
//   - POST /: Creates a new user.
//   - GET /:id: Retrieves a specific user by ID.
//   - PUT /:id: Updates a specific user by ID.
//   - DELETE /:id: Deletes a specific user by ID.
//
// Parameters:
//   - config: The application's configuration settings, used here to set the Gin mode.
//   - userService: An instance of the UserService, which will be injected into the user handlers.
//
// Returns:
//   - A pointer to the configured *gin.Engine instance, ready to be run.
func NewRouter(config config.Config, userService services.UserService) *gin.Engine {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(middleware.Logger())
	engine.Use(gin.Recovery())

	// Explicitly clear trusted proxies (important for security depending on deployment)
	// If behind a trusted proxy (like Cloudflare), you might configure this differently.
	_ = engine.SetTrustedProxies(nil)

	userHandler := handlers.NewUserHandler(userService)

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	userRoutes := engine.Group("/users")
	{
		userRoutes.GET("", userHandler.GetUsers)          // GET /users
		userRoutes.POST("", userHandler.CreateUser)       // POST /users
		userRoutes.GET("/:id", userHandler.GetUserByID)   // GET /users/:id
		userRoutes.PUT("/:id", userHandler.UpdateUser)    // PUT /users/:id
		userRoutes.DELETE("/:id", userHandler.DeleteUser) // DELETE /users/:id
	}

	return engine
}
