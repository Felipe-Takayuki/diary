package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	guiadapter "diary/internal/adapter/gui"
	deliveryhttp "diary/internal/adapter/handler/http"
	"diary/internal/adapter/repository/markdown"
	"diary/internal/adapter/theme"
	"diary/internal/config"
	"diary/internal/usecase"
	"diary/web"
)

// Run bootstraps application dependencies and launches the native desktop app.
func Run() error {
	cfg := config.Load()

	// 1. Secondary Adapter (Persistence / Repository)
	repo, err := markdown.NewRepository(cfg.MetasDir)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// 2. Use Case (Application Layer)
	dailyGoalUseCase := usecase.NewDailyGoalUseCase(repo)

	// 3. Theme Service Adapter
	themeService := theme.NewOmarchyService("")

	// Check if user specifically requested web server mode via CLI flag
	if len(os.Args) > 1 && (os.Args[1] == "--web" || os.Args[1] == "-web" || os.Args[1] == "serve") {
		return runWebServer(cfg, dailyGoalUseCase, themeService)
	}

	logDesktopBanner(cfg)

	// 4. Primary Adapter (Native Desktop GUI)
	desktopApp := guiadapter.NewDesktopApp(dailyGoalUseCase, themeService)
	desktopApp.Run()

	return nil
}

func runWebServer(cfg *config.Config, useCase usecase.DailyGoalUseCase, themeService theme.Service) error {
	handler := deliveryhttp.NewHandler(useCase, web.IndexHTML, themeService)
	handler.SetFavicon(web.FaviconSVG, web.FaviconICO)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	logWebBanner(cfg)

	server := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}

	return server.ListenAndServe()
}

func logDesktopBanner(cfg *config.Config) {
	log.Printf("==================================================")
	log.Printf(" Diary - Native Desktop Application (Omarchy / Hyprland)")
	log.Printf(" Local storage: %s", cfg.MetasDir)
	log.Printf("==================================================")
}

func logWebBanner(cfg *config.Config) {
	log.Printf("==================================================")
	log.Printf(" Web server started successfully (Headless mode)")
	log.Printf(" Available at: http://localhost:%s", cfg.Port)
	log.Printf(" Local storage: %s", cfg.MetasDir)
	log.Printf("==================================================")
}
