package theme_test

import (
	"os"
	"path/filepath"
	"testing"

	"diary/internal/adapter/theme"
)

func TestOmarchyService_FallbackWhenMissing(t *testing.T) {
	svc := theme.NewOmarchyService("/non/existent/path")
	info := svc.GetCurrentTheme()

	if info.Source != "fallback" {
		t.Errorf("expected source fallback, got %s", info.Source)
	}
	if info.ThemeName != "Everforest" {
		t.Errorf("expected theme Everforest, got %s", info.ThemeName)
	}
	if info.Colors["accent"] != "#7fbbb3" {
		t.Errorf("expected accent #7fbbb3, got %s", info.Colors["accent"])
	}
}

func TestOmarchyService_LoadCustomTheme(t *testing.T) {
	tempDir := t.TempDir()
	themeDir := filepath.Join(tempDir, "theme")
	if err := os.MkdirAll(themeDir, 0755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	colorsContent := `
# Sample theme
mode = "light"

accent = "#89b4fa"
background = "#1e1e2e"
foreground = "#cdd6f4"
green = "#a6e3a1"
`
	if err := os.WriteFile(filepath.Join(themeDir, "colors.toml"), []byte(colorsContent), 0644); err != nil {
		t.Fatalf("failed to write colors.toml: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tempDir, "theme.name"), []byte("tokyo-night\n"), 0644); err != nil {
		t.Fatalf("failed to write theme.name: %v", err)
	}

	svc := theme.NewOmarchyService(tempDir)
	info := svc.GetCurrentTheme()

	if info.Source != "omarchy" {
		t.Errorf("expected source omarchy, got %s", info.Source)
	}
	if info.ThemeName != "Tokyo Night" {
		t.Errorf("expected Tokyo Night, got %s", info.ThemeName)
	}
	if info.Mode != "light" {
		t.Errorf("expected mode light, got %s", info.Mode)
	}
	if info.Colors["accent"] != "#89b4fa" {
		t.Errorf("expected accent #89b4fa, got %s", info.Colors["accent"])
	}
}
