package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"diary/internal/domain"
	"diary/internal/usecase"
)

// Handler serves HTTP requests for the daily goals monolith.
type Handler struct {
	useCase    usecase.DailyGoalUseCase
	webContent []byte
}

// NewHandler creates a new Handler with injected use case and web assets.
func NewHandler(useCase usecase.DailyGoalUseCase, webContent []byte) *Handler {
	return &Handler{
		useCase:    useCase,
		webContent: webContent,
	}
}

// RegisterRoutes registers the application routes on the given ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.ServeUI)
	mux.HandleFunc("/api/metas", h.HandleMetas)
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
