package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultPort = "8080"
)

// Config holds the application configuration.
type Config struct {
	Port     string
	MetasDir string
	WebMode  bool
	WithWeb  bool
	Help     bool
}

// DefaultStorageDir determines the canonical directory for storing Markdown daily goals.
// It defaults to ~/metas to ensure consistency across Desktop GUI, CLI, and Web modes.
func DefaultStorageDir() string {
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return filepath.Join(home, "metas")
	}
	return "./metas"
}

// ExpandHome resolves "~/" prefix to the current user's home directory.
func ExpandHome(path string) string {
	if path == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
	} else if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

// Load reads application configuration from environment variables and CLI arguments.
func Load() *Config {
	return LoadFromArgs(os.Args[1:])
}

// LoadFromArgs parses configuration from environment variables and explicit CLI arguments.
func LoadFromArgs(args []string) *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	metasDir := os.Getenv("GOALS_DIR")
	if metasDir == "" {
		metasDir = os.Getenv("METAS_DIR")
	}
	if metasDir == "" {
		metasDir = DefaultStorageDir()
	} else {
		metasDir = ExpandHome(metasDir)
	}

	cfg := &Config{
		Port:     port,
		MetasDir: metasDir,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-h" || arg == "--help" || arg == "help":
			cfg.Help = true
		case arg == "-w" || arg == "--web" || arg == "-web" || arg == "web" || arg == "serve" || arg == "--server":
			cfg.WebMode = true
		case arg == "--with-web":
			cfg.WithWeb = true
		case arg == "-p" || arg == "--port":
			if i+1 < len(args) {
				cfg.Port = args[i+1]
				i++
			}
		case strings.HasPrefix(arg, "--port="):
			cfg.Port = strings.TrimPrefix(arg, "--port=")
		case strings.HasPrefix(arg, "-p="):
			cfg.Port = strings.TrimPrefix(arg, "-p=")
		case arg == "-d" || arg == "--dir" || arg == "--goals-dir":
			if i+1 < len(args) {
				cfg.MetasDir = ExpandHome(args[i+1])
				i++
			}
		case strings.HasPrefix(arg, "--dir="):
			cfg.MetasDir = ExpandHome(strings.TrimPrefix(arg, "--dir="))
		case strings.HasPrefix(arg, "-d="):
			cfg.MetasDir = ExpandHome(strings.TrimPrefix(arg, "-d="))
		case strings.HasPrefix(arg, "--goals-dir="):
			cfg.MetasDir = ExpandHome(strings.TrimPrefix(arg, "--goals-dir="))
		}
	}

	return cfg
}

// Addr returns the network address string to listen on.
func (c *Config) Addr() string {
	return ":" + c.Port
}

// PrintHelp outputs the command-line usage information.
func PrintHelp() {
	fmt.Print(`Diary - Minimalist local-first daily goals manager

Usage:
  diary [flags] [command]

Commands:
  (default)     Launch native desktop GUI
  web, serve    Launch headless web server mode

Flags:
  -w, --web         Launch headless web server mode
      --with-web    Launch desktop GUI and start background web server
  -p, --port string Web server port (default: 8080, or $PORT)
  -d, --dir string  Directory for daily Markdown files (default: ~/metas, or $GOALS_DIR)
  -h, --help        Show this help message

Environment Variables:
  PORT              Web server port (default: 8080)
  GOALS_DIR         Directory for storing goals (default: ~/metas)
  METAS_DIR         Alias for GOALS_DIR

Examples:
  diary                     # Open desktop GUI
  diary --web               # Start web server on http://localhost:8080
  diary web -p 3000         # Start web server on port 3000
  diary --with-web          # Open desktop GUI and start web server
`)
}
