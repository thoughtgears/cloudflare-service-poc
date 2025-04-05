package config

// Config holds application configuration parameters, typically loaded from
// environment variables using a library like 'kelseyhightower/envconfig'.
// Struct tags define the corresponding environment variable names and default values.
type Config struct {
	// Port specifies the network port on which the HTTP server should listen.
	// It defaults to "8080" if the PORT environment variable is not set.
	// Loaded from env: PORT
	Port string `envconfig:"PORT" default:"8080"`

	// Debug enables or disables debug mode for the application.
	// When true, may enable more verbose logging, Gin debug mode, etc.
	// It defaults to false if the DEBUG environment variable is not set or invalid.
	// Loaded from env: DEBUG
	Debug bool `envconfig:"DEBUG" default:"false"`
}
