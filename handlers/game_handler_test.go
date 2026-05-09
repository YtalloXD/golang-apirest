package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/YtalloXD/apirestgo-ia/models"
	"github.com/YtalloXD/apirestgo-ia/storage"
)

// TestCreateGame tests the CreateGame handler
func TestCreateGame(t *testing.T) {
	store := storage.NewGameStore()
	handler := NewGameHandler(store)

	game := models.Game{
		ID:          "1",
		GameName:    "Test Game",
		Publisher:   "Test Publisher",
		Developer:   "Test Developer",
		ReleaseDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		GameGenre:   "RPG",
	}

	body, _ := json.Marshal(game)
	req := httptest.NewRequest("POST", "/api/games", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateGame(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response Response
	json.NewDecoder(w.Body).Decode(&response)

	if !response.Success {
		t.Errorf("Expected success response, got failure")
	}
}

// TestGetAllGames tests the GetAllGames handler
func TestGetAllGames(t *testing.T) {
	store := storage.NewGameStore()
	handler := NewGameHandler(store)

	// Add test data
	game := models.Game{
		ID:          "1",
		GameName:    "Test Game",
		Publisher:   "Test Publisher",
		Developer:   "Test Developer",
		ReleaseDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		GameGenre:   "RPG",
	}
	store.Create(&game)

	req := httptest.NewRequest("GET", "/api/games", nil)
	w := httptest.NewRecorder()

	handler.GetAllGames(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response Response
	json.NewDecoder(w.Body).Decode(&response)

	if !response.Success {
		t.Errorf("Expected success response, got failure")
	}

	games := response.Data.([]interface{})
	if len(games) != 1 {
		t.Errorf("Expected 1 game, got %d", len(games))
	}
}

// TestDeleteGame tests the DeleteGame handler
func TestDeleteGame(t *testing.T) {
	store := storage.NewGameStore()
	handler := NewGameHandler(store)

	// Add test data
	game := models.Game{
		ID:          "1",
		GameName:    "Test Game",
		Publisher:   "Test Publisher",
		Developer:   "Test Developer",
		ReleaseDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		GameGenre:   "RPG",
	}
	store.Create(&game)

	req := httptest.NewRequest("DELETE", "/api/games/1", nil)
	w := httptest.NewRecorder()

	handler.DeleteGame(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Verify game is deleted
	_, err := store.GetByID("1")
	if err != storage.ErrGameNotFound {
		t.Errorf("Expected game to be deleted")
	}
}

// TestGameNotFound tests the 404 error handling
func TestGameNotFound(t *testing.T) {
	store := storage.NewGameStore()
	handler := NewGameHandler(store)

	req := httptest.NewRequest("GET", "/api/games/999", nil)
	w := httptest.NewRecorder()

	handler.GetGameByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
