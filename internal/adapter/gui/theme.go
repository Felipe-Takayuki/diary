package gui

import (
	_ "embed"
	"image/color"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	omarchytheme "diary/internal/adapter/theme"
)

//go:embed icon.png
var AppIconBytes []byte

// OmarchyTheme implements fyne.Theme adapting to Omarchy system colors and tokens.
type OmarchyTheme struct {
	fallback fyne.Theme
	info     omarchytheme.ThemeInfo
	fontRes  fyne.Resource
}

// NewOmarchyTheme creates a theme adapted to the current Omarchy environment.
func NewOmarchyTheme(info omarchytheme.ThemeInfo) *OmarchyTheme {
	var fontRes fyne.Resource
	fontPaths := []string{
		"/usr/share/fonts/TTF/JetBrainsMonoNerdFont-Regular.ttf",
		"/usr/share/fonts/TTF/JetBrainsMono-Regular.ttf",
		"/usr/share/fonts/truetype/jetbrains-mono/JetBrainsMono-Regular.ttf",
	}

	for _, p := range fontPaths {
		if data, err := os.ReadFile(p); err == nil && len(data) > 0 {
			fontRes = fyne.NewStaticResource("JetBrainsMono.ttf", data)
			break
		}
	}

	return &OmarchyTheme{
		fallback: theme.DarkTheme(),
		info:     info,
		fontRes:  fontRes,
	}
}

func (t *OmarchyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	colors := t.info.Colors

	switch name {
	case theme.ColorNameBackground:
		if c, ok := colors["dark_background"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0x21, G: 0x27, B: 0x2c, A: 0xff})
		}
		if c, ok := colors["background"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0x2d, G: 0x35, B: 0x3b, A: 0xff})
		}
	case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground:
		if c, ok := colors["background"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0x2d, G: 0x35, B: 0x3b, A: 0xff})
		}
	case theme.ColorNameForeground:
		if c, ok := colors["foreground"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0xd3, G: 0xc6, B: 0xaa, A: 0xff})
		}
	case theme.ColorNamePrimary:
		if c, ok := colors["accent"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0x7f, G: 0xbb, B: 0xb3, A: 0xff})
		}
	case theme.ColorNameHover:
		if c, ok := colors["lighter_background"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0x34, G: 0x3f, B: 0x44, A: 0xff})
		}
	case theme.ColorNameFocus:
		if c, ok := colors["accent"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0x7f, G: 0xbb, B: 0xb3, A: 0xff})
		}
	case theme.ColorNameInputBackground:
		if c, ok := colors["darker_background"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0x18, G: 0x1d, B: 0x20, A: 0xff})
		}
	case theme.ColorNameInputBorder:
		if c, ok := colors["muted"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0x47, G: 0x52, B: 0x58, A: 0xff})
		}
	case theme.ColorNameButton:
		if c, ok := colors["background"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0x2d, G: 0x35, B: 0x3b, A: 0xff})
		}
	case theme.ColorNameSuccess:
		if c, ok := colors["green"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0xa7, G: 0xc0, B: 0x80, A: 0xff})
		}
	case theme.ColorNameError:
		if c, ok := colors["red"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0xe6, G: 0x7e, B: 0x80, A: 0xff})
		}
	case theme.ColorNameSelection:
		if c, ok := colors["selection"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0x3d, G: 0x48, B: 0x4d, A: 0xff})
		}
	case theme.ColorNameSeparator:
		if c, ok := colors["muted"]; ok {
			return parseHexColor(c, color.NRGBA{R: 0x47, G: 0x52, B: 0x58, A: 0xff})
		}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 40}
	}

	return color.NRGBA{R: 128, G: 128, B: 128, A: 255}
}

func (t *OmarchyTheme) Font(style fyne.TextStyle) fyne.Resource {
	if t.fontRes != nil {
		return t.fontRes
	}
	return t.fallback.Font(style)
}

func (t *OmarchyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.fallback.Icon(name)
}

func (t *OmarchyTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameButtonRadius, theme.SizeNameCardRadius, theme.SizeNameInputRadius, theme.SizeNameSelectionRadius, theme.SizeNameMenuRadius:
		return 5.0
	case theme.SizeNamePadding:
		return 8.0
	case theme.SizeNameInlineIcon:
		return 18.0
	case theme.SizeNameText:
		return 13.0
	case theme.SizeNameHeadingText:
		return 18.0
	case theme.SizeNameSubHeadingText:
		return 14.0
	case theme.SizeNameCaptionText:
		return 11.0
	}
	return t.fallback.Size(name)
}

func parseHexColor(s string, fallback color.Color) color.Color {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(s) == 3 {
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	}
	if len(s) != 6 && len(s) != 8 {
		return fallback
	}
	val, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return fallback
	}
	if len(s) == 6 {
		return color.NRGBA{
			R: uint8(val >> 16),
			G: uint8((val >> 8) & 0xFF),
			B: uint8(val & 0xFF),
			A: 0xFF,
		}
	}
	return color.NRGBA{
		R: uint8(val >> 24),
		G: uint8((val >> 16) & 0xFF),
		B: uint8((val >> 8) & 0xFF),
		A: uint8(val & 0xFF),
	}
}
