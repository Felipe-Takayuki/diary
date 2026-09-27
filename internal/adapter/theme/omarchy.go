package theme

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ThemeInfo represents the active theme metadata and color definitions.
type ThemeInfo struct {
	Source    string            `json:"source"`
	ThemeName string            `json:"themeName"`
	Mode      string            `json:"mode"`
	Colors    map[string]string `json:"colors"`
}

// Service defines the interface for retrieving the active theme.
type Service interface {
	GetCurrentTheme() ThemeInfo
}

// OmarchyService loads theme colors from the local Omarchy environment.
type OmarchyService struct {
	basePath string
}

// NewOmarchyService creates a new theme service.
// If customBasePath is empty, it resolves to $HOME/.local/state/omarchy/current.
func NewOmarchyService(customBasePath string) *OmarchyService {
	return &OmarchyService{
		basePath: customBasePath,
	}
}

// GetCurrentTheme returns the active theme information.
func (s *OmarchyService) GetCurrentTheme() ThemeInfo {
	base := s.basePath
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return defaultFallbackTheme()
		}
		base = filepath.Join(home, ".local", "state", "omarchy", "current")
	}

	colorsFile := filepath.Join(base, "theme", "colors.toml")
	colors, mode, err := parseColorsTOML(colorsFile)
	if err != nil || len(colors) == 0 {
		return defaultFallbackTheme()
	}

	themeName := readThemeName(filepath.Join(base, "theme.name"))
	if themeName == "" {
		themeName = "Omarchy"
	}

	return ThemeInfo{
		Source:    "omarchy",
		ThemeName: themeName,
		Mode:      mode,
		Colors:    colors,
	}
}

func readThemeName(nameFile string) string {
	data, err := os.ReadFile(nameFile)
	if err != nil {
		return ""
	}
	raw := strings.TrimSpace(string(data))
	if raw == "" {
		return ""
	}

	parts := strings.Split(raw, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

func parseColorsTOML(filePath string) (map[string]string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	colors := make(map[string]string)
	mode := "dark"

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"'")

		if key == "mode" {
			if val == "light" {
				mode = "light"
			} else {
				mode = "dark"
			}
			continue
		}

		colors[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, "", err
	}

	return colors, mode, nil
}

func defaultFallbackTheme() ThemeInfo {
	return ThemeInfo{
		Source:    "fallback",
		ThemeName: "Everforest",
		Mode:      "dark",
		Colors: map[string]string{
			"accent":             "#7fbbb3",
			"selection":          "#3d484d",
			"muted":              "#475258",
			"background":         "#2d353b",
			"dark_background":    "#21272c",
			"darker_background":  "#181d20",
			"lighter_background": "#343f44",
			"foreground":         "#d3c6aa",
			"dark_foreground":    "#4f585e",
			"light_foreground":   "#9da9a0",
			"bright_foreground":  "#d3c6aa",
			"red":                "#e67e80",
			"yellow":             "#dbbc7f",
			"orange":             "#e09d7f",
			"green":              "#a7c080",
			"cyan":               "#83c092",
			"blue":               "#7fbbb3",
			"magenta":            "#d699b6",
			"brown":              "#704e3f",
		},
	}
}
