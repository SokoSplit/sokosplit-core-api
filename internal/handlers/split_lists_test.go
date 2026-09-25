package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sokosplit/sokosplit-core-api/internal/events"
	"github.com/sokosplit/sokosplit-core-api/internal/handlers"
	"github.com/sokosplit/sokosplit-core-api/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTestRouter wires a real handler against an in-memory SQLite DB, and
// injects a fixed user_id instead of running real JWT auth, so these tests
// exercise handler logic without needing Postgres or a signed token.
func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB, uuid.UUID) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.SplitList{}, &models.Recipient{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	ownerID := uuid.New()
	bus := &events.Bus{}
	h := handlers.NewSplitListHandler(db, bus)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", ownerID)
		c.Next()
	})
	r.POST("/split-lists", h.Create)
	r.GET("/split-lists", h.List)
	r.GET("/split-lists/:id", h.Get)

	return r, db, ownerID
}

func TestCreateSplitList_ValidBps(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	body := map[string]any{
		"split_id": "job_1",
		"token":    "GABC...",
		"amount":   1000,
		"recipients": []map[string]any{
			{"address": "GRECIPIENTA", "bps": 7000},
			{"address": "GRECIPIENTB", "bps": 3000},
		},
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/split-lists", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateSplitList_InvalidBpsSum(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	body := map[string]any{
		"split_id": "job_2",
		"token":    "GABC...",
		"amount":   1000,
		"recipients": []map[string]any{
			{"address": "GRECIPIENTA", "bps": 5000}, // doesn't sum to 10000
		},
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/split-lists", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}
