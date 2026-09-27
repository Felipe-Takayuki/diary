package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	deliveryhttp "diary/internal/adapter/handler/http"
	"diary/internal/domain"
)

type mockGoalUseCase struct {
	getFn  func(ctx context.Context, rawDate string) (*domain.DailyGoal, error)
	saveFn func(ctx context.Context, rawDate string, items []domain.Item) (*domain.DailyGoal, error)
}

func (m *mockGoalUseCase) GetDailyGoals(ctx context.Context, rawDate string) (*domain.DailyGoal, error) {
	if m.getFn != nil {
		return m.getFn(ctx, rawDate)
	}
	date, _ := domain.ParseDate(rawDate)
	return domain.NewDailyGoal(date, nil), nil
}

func (m *mockGoalUseCase) SaveDailyGoals(ctx context.Context, rawDate string, items []domain.Item) (*domain.DailyGoal, error) {
	if m.saveFn != nil {
		return m.saveFn(ctx, rawDate, items)
	}
	date, _ := domain.ParseDate(rawDate)
	return domain.NewDailyGoal(date, items), nil
}

func TestHandler_ServeUI(t *testing.T) {
	htmlContent := []byte("<html><body>Mock App</body></html>")
	h := deliveryhttp.NewHandler(&mockGoalUseCase{}, htmlContent, nil)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	t.Run("GET / returns 200 with HTML", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
		}
		if !strings.Contains(rec.Header().Get("Content-Type"), "text/html") {
			t.Errorf("expected text/html content-type, got %s", rec.Header().Get("Content-Type"))
		}
		if !strings.Contains(rec.Body.String(), "Mock App") {
			t.Errorf("expected body to contain Mock App, got %s", rec.Body.String())
		}
	})

	t.Run("GET /unknown returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("POST / returns 405 Method Not Allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
		}
	})
}

func TestHandler_GetMetas(t *testing.T) {
	date := domain.MustParseDate("25-09-2026")
	item1, _ := domain.NewItem("Meta 1", true)
	item2, _ := domain.NewItem("Meta 2", false)

	mockUC := &mockGoalUseCase{
		getFn: func(_ context.Context, rawDate string) (*domain.DailyGoal, error) {
			if rawDate == "invalid" {
				return nil, domain.ErrInvalidDate
			}
			return domain.NewDailyGoal(date, []domain.Item{item1, item2}), nil
		},
	}

	h := deliveryhttp.NewHandler(mockUC, []byte(""), nil)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	t.Run("GET /api/metas with valid date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/metas?date=25-09-2026", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
		}

		var res deliveryhttp.MetasResponse
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Date != "25-09-2026" {
			t.Errorf("got date %s, want 25-09-2026", res.Date)
		}
		if len(res.Items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(res.Items))
		}
		if res.Items[0].Text != "Meta 1" || !res.Items[0].Done {
			t.Errorf("unexpected item 0: %+v", res.Items[0])
		}
	})

	t.Run("GET /api/metas with invalid date returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/metas?date=invalid", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}

func TestHandler_SaveMetas(t *testing.T) {
	mockUC := &mockGoalUseCase{
		saveFn: func(_ context.Context, rawDate string, items []domain.Item) (*domain.DailyGoal, error) {
			if rawDate == "invalid" {
				return nil, domain.ErrInvalidDate
			}
			d, _ := domain.ParseDate(rawDate)
			return domain.NewDailyGoal(d, items), nil
		},
	}

	h := deliveryhttp.NewHandler(mockUC, []byte(""), nil)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	t.Run("POST /api/metas with valid payload", func(t *testing.T) {
		payload := `{"date":"25-09-2026","items":[{"text":"Nova meta","done":false}]}`
		req := httptest.NewRequest(http.MethodPost, "/api/metas", bytes.NewBufferString(payload))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
		}

		var res deliveryhttp.StatusResponse
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if res.Status != "ok" {
			t.Errorf("got status %s, want ok", res.Status)
		}
	})

	t.Run("POST /api/metas with bad JSON returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/metas", bytes.NewBufferString("{bad json}"))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("DELETE /api/metas returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/metas", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
		}
	})
}

func TestHandler_GetTheme(t *testing.T) {
	h := deliveryhttp.NewHandler(&mockGoalUseCase{}, []byte(""), nil)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	t.Run("GET /api/theme returns 200 with theme payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/theme", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
		}

		var data map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&data); err != nil {
			t.Fatalf("failed to decode json: %v", err)
		}

		if _, ok := data["colors"]; !ok {
			t.Errorf("expected colors field in response, got %+v", data)
		}
	})

	t.Run("POST /api/theme returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/theme", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
		}
	})
}

func TestHandler_Favicons(t *testing.T) {
	h := deliveryhttp.NewHandler(&mockGoalUseCase{}, []byte(""), nil)
	h.SetFavicon([]byte("<svg></svg>"), []byte("fake-ico"))
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	t.Run("GET /favicon.svg returns 200 with SVG", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/favicon.svg", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
		}
		if !strings.Contains(rec.Header().Get("Content-Type"), "image/svg+xml") {
			t.Errorf("expected image/svg+xml, got %s", rec.Header().Get("Content-Type"))
		}
	})

	t.Run("GET /favicon.ico returns 200 with ICO", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
		}
		if !strings.Contains(rec.Header().Get("Content-Type"), "image/x-icon") {
			t.Errorf("expected image/x-icon, got %s", rec.Header().Get("Content-Type"))
		}
	})

	t.Run("POST /favicon.svg returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/favicon.svg", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
		}
	})
}
