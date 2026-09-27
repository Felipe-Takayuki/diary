package gui

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/theme"

	omarchytheme "diary/internal/adapter/theme"
)

func TestParseHexColor(t *testing.T) {
	fallback := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

	tests := []struct {
		input    string
		wantR    uint8
		wantG    uint8
		wantB    uint8
		wantA    uint8
		fallback bool
	}{
		{input: "#7fbbb3", wantR: 0x7f, wantG: 0xbb, wantB: 0xb3, wantA: 0xff},
		{input: "7fbbb3", wantR: 0x7f, wantG: 0xbb, wantB: 0xb3, wantA: 0xff},
		{input: "#fff", wantR: 0xff, wantG: 0xff, wantB: 0xff, wantA: 0xff},
		{input: "invalid", fallback: true},
		{input: "", fallback: true},
	}

	for _, tt := range tests {
		got := parseHexColor(tt.input, fallback)
		nrgba, ok := got.(color.NRGBA)
		if !ok {
			t.Fatalf("expected color.NRGBA, got %T", got)
		}

		if tt.fallback {
			if nrgba != fallback {
				t.Errorf("input %q: got %v, want fallback %v", tt.input, nrgba, fallback)
			}
			continue
		}

		if nrgba.R != tt.wantR || nrgba.G != tt.wantG || nrgba.B != tt.wantB || nrgba.A != tt.wantA {
			t.Errorf("input %q: got RGBA(%d, %d, %d, %d), want RGBA(%d, %d, %d, %d)",
				tt.input, nrgba.R, nrgba.G, nrgba.B, nrgba.A, tt.wantR, tt.wantG, tt.wantB, tt.wantA)
		}
	}
}

func TestOmarchyTheme_ColorsAndSizes(t *testing.T) {
	info := omarchytheme.ThemeInfo{
		ThemeName: "Everforest",
		Colors: map[string]string{
			"accent":          "#7fbbb3",
			"background":      "#2d353b",
			"dark_background": "#21272c",
			"foreground":      "#d3c6aa",
			"muted":           "#475258",
			"green":           "#a7c080",
			"red":             "#e67e80",
		},
	}

	th := NewOmarchyTheme(info)

	// Test primary color matches accent
	primary := th.Color(theme.ColorNamePrimary, theme.VariantDark)
	nrgba, ok := primary.(color.NRGBA)
	if !ok {
		t.Fatalf("expected color.NRGBA, got %T", primary)
	}
	if nrgba.R != 0x7f || nrgba.G != 0xbb || nrgba.B != 0xb3 {
		t.Errorf("got primary color %v, want #7fbbb3", nrgba)
	}

	// Test corner radius is 5px (Omarchy standard)
	radius := th.Size(theme.SizeNameInputRadius)
	if radius != 5.0 {
		t.Errorf("got input radius %f, want 5.0", radius)
	}

	buttonRadius := th.Size(theme.SizeNameButtonRadius)
	if buttonRadius != 5.0 {
		t.Errorf("got button radius %f, want 5.0", buttonRadius)
	}
}
