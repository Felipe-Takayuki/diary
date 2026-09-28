package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	guiadapter "diary/internal/adapter/gui"
	deliveryhttp "diary/internal/adapter/handler/http"
	"diary/internal/adapter/repository/markdown"
	"diary/internal/adapter/theme"
	"diary/internal/config"
	"diary/internal/usecase"
	"diary/web"
)

// Run bootstraps application dependencies and launches the application.
func Run() error {
	cfg := config.Load()
	if cfg.Help {
		config.PrintHelp()
		return nil
	}

	// 1. Secondary Adapter (Persistence / Repository)
	repo, err := markdown.NewRepository(cfg.MetasDir)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// 2. Use Case (Application Layer)
	dailyGoalUseCase := usecase.NewDailyGoalUseCase(repo)

	// 3. Theme Service Adapter
	themeService := theme.NewOmarchyService("")

	// Check if user specifically requested web server mode
	if cfg.WebMode {
		return runWebServer(cfg, dailyGoalUseCase, themeService)
	}

	// Start background web server alongside desktop GUI if --with-web was passed
	if cfg.WithWeb {
		go func() {
			if err := runWebServer(cfg, dailyGoalUseCase, themeService); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Printf("Background web server error: %v", err)
			}
		}()
	}

	logDesktopBanner(cfg)

	// 4. Primary Adapter (Native Desktop GUI)
	desktopApp := guiadapter.NewDesktopApp(dailyGoalUseCase, themeService, cfg.MetasDir)
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
	if cfg.WithWeb {
		log.Printf(" Web server running concurrently on: http://localhost:%s", cfg.Port)
	}
	log.Printf("==================================================")
}

func logWebBanner(cfg *config.Config) {
	log.Printf("==================================================")
	log.Printf(" Web server started successfully (Headless mode)")
	log.Printf(" Available at: http://localhost:%s", cfg.Port)
	log.Printf(" Local storage: %s", cfg.MetasDir)
	log.Printf("==================================================")
}
