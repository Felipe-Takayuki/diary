package config

import "os"

const (
	defaultPort     = "8080"
	defaultMetasDir = "./metas"
)

// Config holds the application configuration.
type Config struct {
	Port     string
	MetasDir string
}

// Load reads application configuration from environment variables or returns defaults.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	metasDir := os.Getenv("GOALS_DIR")
	if metasDir == "" {
		metasDir = os.Getenv("METAS_DIR")
	}
	if metasDir == "" {
		metasDir = defaultMetasDir
	}

	return &Config{
		Port:     port,
		MetasDir: metasDir,
	}
}

// Addr returns the network address string to listen on.
func (c *Config) Addr() string {
	return ":" + c.Port
}
