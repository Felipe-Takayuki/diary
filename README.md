<p align="center">
  <img src="assets/logo.svg" alt="Diary Logo" width="540">
</p>

<p align="center">
  <strong>Minimalist local-first daily goals manager in Go with native Omarchy aesthetic and desktop integration.</strong>
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-7fbbb3.svg" alt="License: MIT"></a>
  <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.22+-83c092.svg" alt="Go 1.22+"></a>
  <img src="https://img.shields.io/badge/Dependencies-Zero_External_Runtime-d3c6aa.svg" alt="Zero External Runtime Dependencies">
  <img src="https://img.shields.io/badge/Theme-Omarchy_Adaptive-7fbbb3.svg" alt="Omarchy Adaptive">
</p>

---

Your daily goals are stored directly on your computer as standard Markdown files (`./metas/DD-MM-YYYY.md`). No external databases, no cloud subscriptions, and zero tracking.

Built on **Clean Architecture**, Diary provides a native graphical desktop interface tailored for **Omarchy / Hyprland**, along with an optional headless web mode embedded directly into the binary.

---

## Key Features

- **Local-First Markdown Persistence:** Each day has its own human-readable Markdown file formatted with standard checklists (`- [ ] ` and `- [x] `), instantly editable with any text editor or Obsidian.
- **Native Desktop GUI:** Lightweight graphical desktop app built with Fyne and customized with Omarchy design tokens and active theme sync.
- **Zero External Runtime Dependencies:** Pure Go core with embedded static assets via `go:embed`.
- **Omarchy Theme Synchronization:** Automatically detects and applies active system themes (Everforest, Catppuccin, Gruvbox, Nord, Rose Pine, Tokyo Night, and light variants) from `colors.toml`.
- **Accessible & High-Contrast:** Strict compliance with WCAG AAA contrast, JetBrains Mono typography, clean geometry, and instant visual progress indicators.
- **Fluid Keyboard Navigation:** Shortcuts for quick task entry (`/`), changing dates (`←` / `→`), and jumping straight to today (`T`).
- **Safe Concurrent Access:** File access guarded by `sync.RWMutex` with path traversal protections.

---

## Project Structure

```
diary/
├── assets/                            # Brand assets (horizontal logo, vector icons, .desktop launcher)
│   ├── logo.svg
│   ├── icon.svg
│   ├── favicon.svg
│   └── com.omarchy.diary.desktop
├── cmd/
│   └── diary/
│       └── main.go                    # Standard application entrypoint
├── internal/
│   ├── app/
│   │   └── app.go                     # Dependency injection & bootstrap
│   ├── domain/                        # Pure domain rules & business entities
│   │   ├── date.go                    # Canonical Date value object (DD-MM-YYYY)
│   │   ├── errors.go                  # Domain error definitions
│   │   ├── item.go                    # Goal Item entity (validation & toggle)
│   │   ├── meta.go                    # DailyGoal aggregate root (progress & metrics)
│   │   └── repository.go              # GoalRepository contract interface
│   ├── usecase/                       # Application orchestration
│   │   └── daily_goal.go              # DailyGoalUseCase (GetDailyGoals, SaveDailyGoals)
│   ├── adapter/                       # Inbound & outbound adapters
│   │   ├── gui/                       # Native Desktop GUI (Fyne + Omarchy Theme)
│   │   ├── handler/http/              # Headless HTTP API & embedded web UI
│   │   ├── repository/markdown/       # Concurrent Markdown file persistence
│   │   └── theme/                     # Omarchy colors.toml integration & theme watcher
│   └── config/
│       └── config.go                  # Environment configuration (PORT, GOALS_DIR/METAS_DIR)
├── web/
│   ├── embed.go                       # Embedded assets via go:embed
│   ├── static/                        # Vector and raster favicons
│   └── template/
│       └── index.html                 # Optional web interface
├── metas/                             # Directory where daily Markdown files are stored
├── Makefile                           # Development, test, build, and installation tasks
├── CONTRIBUTING.md                    # Guidelines for contributors
├── LICENSE                            # MIT License
├── go.mod
└── main.go                            # Convenient root entrypoint
```

---

## Installation & Running

### Prerequisites
- [Go 1.22+](https://go.dev/dl/) installed.
- Linux graphical environment (Wayland / Hyprland or X11).

### 1. Run the Desktop App
Clone the repository and launch directly:

```bash
git clone https://github.com/Felipe-Takayuki/diary.git
cd diary
go run .
```

The native window will appear instantly with full keyboard and mouse support.

### 2. Install System-Wide (Omarchy Launcher)
To integrate Diary directly into system launchers (`rofi`, `walker`, or desktop menus):

```bash
make install
```

This compiles the binary to `~/.local/bin/diary`, installs the `.desktop` desktop file to `~/.local/share/applications/`, and registers the high-resolution vector and raster icons.

### 3. Headless Web Mode (Optional)
If you prefer running a local browser server:

```bash
# Using the installed binary from anywhere:
diary --web
# Or: diary web -p 8080

# Or via Makefile inside the project:
make web

# Available at http://localhost:8080
```

To run both the Desktop GUI and background Web server simultaneously:
```bash
diary --with-web
```

---

## Keyboard Shortcuts

| Shortcut | Description |
|---|---|
| `/` | Focus goal input field |
| `←` | Navigate to previous day |
| `→` | Navigate to next day |
| `T` | Jump to today's goals |
| `Enter` | Save and submit new goal |
| `Esc` | Unfocus input field |

---

## Configuration & Flags

All modes (Desktop, Web, and CLI) automatically share the same canonical data folder: `~/metas`.

### CLI Flags & Commands

| Command / Flag | Description |
|---|---|
| `diary` | Opens the native Desktop GUI |
| `diary --web` (or `diary web`, `diary -w`) | Starts headless HTTP web server |
| `diary --with-web` | Opens Desktop GUI and starts background web server |
| `-p, --port <port>` | Sets web server port (default: `8080`) |
| `-d, --dir <path>` | Sets Markdown storage directory (default: `~/metas`) |
| `-h, --help` | Displays command-line help |

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Port used by the HTTP server in web mode |
| `GOALS_DIR` (or `METAS_DIR`) | `~/metas` | Path to the directory where Markdown files are stored |

Example running with a custom data folder:
```bash
GOALS_DIR=~/Documents/MyGoals diary
```

---

## API Reference

### `GET /api/goals?date=DD-MM-YYYY`
Retrieves daily goals for the specified date. Also accessible via `/api/metas`.

**Success Response (`200 OK`):**
```json
{
  "date": "27-09-2026",
  "items": [
    {
      "text": "Review architecture diagrams",
      "done": true
    },
    {
      "text": "Write open source documentation",
      "done": false
    }
  ]
}
```

### `POST /api/goals`
Saves goal items to the corresponding daily Markdown file. Also accessible via `/api/metas`.

**Request Payload:**
```json
{
  "date": "27-09-2026",
  "items": [
    {
      "text": "Review architecture diagrams",
      "done": true
    }
  ]
}
```

### `GET /api/theme`
Returns the active theme palette and mode (Omarchy or fallback).

**Example Response (`200 OK`):**
```json
{
  "source": "omarchy",
  "themeName": "Everforest",
  "mode": "dark",
  "colors": {
    "accent": "#7fbbb3",
    "background": "#2d353b",
    "foreground": "#d3c6aa",
    "green": "#a7c080",
    "red": "#e67e80"
  }
}
```

---

## Development & Testing

Comprehensive unit and concurrency race tests are included:

```bash
# Run unit tests
make test

# Run tests with race detection
make test-race

# Generate code coverage report
make test-cover

# Check formatting and static analysis
make fmt
make vet
```

---

## Contributing

Contributions are very welcome! Please check our [Contributing Guide](file:///home/takayuki/Projects/diary/CONTRIBUTING.md) for details on code standards, branch naming, and pull request workflows.

---

## License

Distributed under the [MIT License](file:///home/takayuki/Projects/diary/LICENSE). See `LICENSE` for more details.
