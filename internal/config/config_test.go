package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"diary/internal/config"
)

func TestConfig_Defaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("GOALS_DIR")
	os.Unsetenv("METAS_DIR")

	cfg := config.LoadFromArgs([]string{})
	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}

	home, _ := os.UserHomeDir()
	expectedDir := filepath.Join(home, "metas")
	if cfg.MetasDir != expectedDir {
		t.Errorf("expected default metas dir %s, got %s", expectedDir, cfg.MetasDir)
	}
	if cfg.WebMode {
		t.Errorf("expected WebMode=false by default")
	}
	if cfg.WithWeb {
		t.Errorf("expected WithWeb=false by default")
	}
}

func TestConfig_EnvOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("GOALS_DIR", "/custom/goals")

	cfg := config.LoadFromArgs([]string{})
	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Port)
	}
	if cfg.MetasDir != "/custom/goals" {
		t.Errorf("expected metas dir /custom/goals, got %s", cfg.MetasDir)
	}
}

func TestConfig_CLIArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantWeb     bool
		wantWithWeb bool
		wantPort    string
		wantDir     string
		wantHelp    bool
	}{
		{
			name:    "flag --web",
			args:    []string{"--web"},
			wantWeb: true,
		},
		{
			name:    "flag -w",
			args:    []string{"-w"},
			wantWeb: true,
		},
		{
			name:    "subcommand web",
			args:    []string{"web"},
			wantWeb: true,
		},
		{
			name:    "subcommand serve",
			args:    []string{"serve"},
			wantWeb: true,
		},
		{
			name:        "flag --with-web",
			args:        []string{"--with-web"},
			wantWithWeb: true,
		},
		{
			name:     "port flag -p",
			args:     []string{"--web", "-p", "3000"},
			wantWeb:  true,
			wantPort: "3000",
		},
		{
			name:     "port flag --port=",
			args:     []string{"--port=4000"},
			wantPort: "4000",
		},
		{
			name:    "dir flag -d",
			args:    []string{"-d", "/tmp/test-metas"},
			wantDir: "/tmp/test-metas",
		},
		{
			name:     "help flag -h",
			args:     []string{"-h"},
			wantHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.LoadFromArgs(tt.args)
			if tt.wantWeb && !cfg.WebMode {
				t.Errorf("expected WebMode=true")
			}
			if tt.wantWithWeb && !cfg.WithWeb {
				t.Errorf("expected WithWeb=true")
			}
			if tt.wantPort != "" && cfg.Port != tt.wantPort {
				t.Errorf("expected Port=%s, got %s", tt.wantPort, cfg.Port)
			}
			if tt.wantDir != "" && cfg.MetasDir != tt.wantDir {
				t.Errorf("expected MetasDir=%s, got %s", tt.wantDir, cfg.MetasDir)
			}
			if tt.wantHelp && !cfg.Help {
				t.Errorf("expected Help=true")
			}
		})
	}
}

func TestConfig_ExpandHome(t *testing.T) {
	home, _ := os.UserHomeDir()
	if home == "" {
		t.Skip("UserHomeDir not available")
	}

	expanded := config.ExpandHome("~/my-metas")
	expected := filepath.Join(home, "my-metas")
	if expanded != expected {
		t.Errorf("expected %s, got %s", expected, expanded)
	}

	expandedExact := config.ExpandHome("~")
	if expandedExact != home {
		t.Errorf("expected %s, got %s", home, expandedExact)
	}

	untouched := config.ExpandHome("/absolute/path")
	if untouched != "/absolute/path" {
		t.Errorf("expected /absolute/path, got %s", untouched)
	}
}
