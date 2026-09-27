package app

import (
	"fmt"
	"log"
	"net/http"

	deliveryhttp "diary/internal/adapter/handler/http"
	"diary/internal/adapter/repository/markdown"
	"diary/internal/config"
	"diary/internal/usecase"
	"diary/web"
)

// Run bootstraps the application dependencies and starts the HTTP server.
func Run() error {
	cfg := config.Load()

	// 1. Secondary Adapter (Persistence / Repository)
	repo, err := markdown.NewRepository(cfg.MetasDir)
	if err != nil {
		return fmt.Errorf("erro na inicialização do armazenamento: %w", err)
	}

	// 2. Use Case (Application Layer)
	dailyGoalUseCase := usecase.NewDailyGoalUseCase(repo)

	// 3. Primary Adapter (HTTP Handler)
	handler := deliveryhttp.NewHandler(dailyGoalUseCase, web.IndexHTML)

	// 4. Router setup
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	logBanner(cfg)

	// 5. HTTP Server startup
	server := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}

	return server.ListenAndServe()
}

func logBanner(cfg *config.Config) {
	log.Printf("==================================================")
	log.Printf(" Servidor de Metas Diárias iniciado com sucesso! ")
	log.Printf(" Formato de datas: Dia-Mês-Ano (DD-MM-YYYY)      ")
	log.Printf(" Acesse em: http://localhost:%s                  ", cfg.Port)
	log.Printf(" Diretório de armazenamento: %s                  ", cfg.MetasDir)
	log.Printf("==================================================")
}
