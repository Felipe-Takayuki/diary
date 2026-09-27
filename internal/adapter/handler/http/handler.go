package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"diary/internal/adapter/theme"
	"diary/internal/domain"
	"diary/internal/usecase"
)

// Handler serves HTTP requests for the daily goals monolith.
type Handler struct {
	useCase      usecase.DailyGoalUseCase
	webContent   []byte
	faviconSVG   []byte
	faviconICO   []byte
	themeService theme.Service
}

// NewHandler creates a new Handler with injected use case, web assets, and optional theme service.
func NewHandler(useCase usecase.DailyGoalUseCase, webContent []byte, themeService theme.Service) *Handler {
	if themeService == nil {
		themeService = theme.NewOmarchyService("")
	}
	return &Handler{
		useCase:      useCase,
		webContent:   webContent,
		themeService: themeService,
	}
}

// SetFavicon configures embedded favicon data for SVG and ICO formats.
func (h *Handler) SetFavicon(svg, ico []byte) {
	h.faviconSVG = svg
	h.faviconICO = ico
}

// RegisterRoutes registers the application routes on the given ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.ServeUI)
	mux.HandleFunc("/favicon.svg", h.ServeFaviconSVG)
	mux.HandleFunc("/favicon.ico", h.ServeFaviconICO)
	mux.HandleFunc("/api/metas", h.HandleMetas)
	mux.HandleFunc("/api/theme", h.HandleTheme)
}

// ServeFaviconSVG serves the vector favicon.
func (h *Handler) ServeFaviconSVG(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	if len(h.faviconSVG) == 0 {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.faviconSVG)
}

// ServeFaviconICO serves the ICO favicon.
func (h *Handler) ServeFaviconICO(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	if len(h.faviconICO) == 0 {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.faviconICO)
}

// ServeUI handles GET / serving the single-page application.
func (h *Handler) ServeUI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.webContent)
}

// HandleMetas dispatches GET and POST requests for /api/metas.
func (h *Handler) HandleMetas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getMetas(w, r)
	case http.MethodPost:
		h.saveMetas(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// HandleTheme handles GET /api/theme returning active theme colors.
func (h *Handler) HandleTheme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	h.writeJSON(w, http.StatusOK, h.themeService.GetCurrentTheme())
}

func (h *Handler) getMetas(w http.ResponseWriter, r *http.Request) {
	rawDate := r.URL.Query().Get("date")
	goal, err := h.useCase.GetDailyGoals(r.Context(), rawDate)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidDate) || errors.Is(err, domain.ErrEmptyDate) {
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Printf("Erro ao carregar metas (%s): %v", rawDate, err)
		h.writeJSONError(w, http.StatusInternalServerError, "Erro interno ao carregar metas")
		return
	}

	itemsDTO := make([]ItemDTO, 0, goal.TotalCount())
	for _, it := range goal.Items() {
		itemsDTO = append(itemsDTO, ItemDTO{
			Text: it.Text(),
			Done: it.Done(),
		})
	}

	h.writeJSON(w, http.StatusOK, MetasResponse{
		Date:  goal.Date().String(),
		Items: itemsDTO,
	})
}

func (h *Handler) saveMetas(w http.ResponseWriter, r *http.Request) {
	var payload SaveMetasRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	domainItems := make([]domain.Item, 0, len(payload.Items))
	for _, it := range payload.Items {
		item, err := domain.NewItem(it.Text, it.Done)
		if err != nil {
			// Skip empty items, consistent with domain rules
			continue
		}
		domainItems = append(domainItems, item)
	}

	_, err := h.useCase.SaveDailyGoals(r.Context(), payload.Date, domainItems)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidDate) || errors.Is(err, domain.ErrEmptyDate) {
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Printf("Erro ao salvar metas (%s): %v", payload.Date, err)
		h.writeJSONError(w, http.StatusInternalServerError, "Erro interno ao salvar metas")
		return
	}

	h.writeJSON(w, http.StatusOK, StatusResponse{Status: "ok"})
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Erro ao codificar resposta JSON: %v", err)
	}
}

func (h *Handler) writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
