package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/YtalloXD/apirestgo-ia/models"
	"github.com/YtalloXD/apirestgo-ia/repository"
	"github.com/gorilla/mux"
)

type fakeGameRepository struct {
	games map[string]*models.Game
}

func newFakeGameRepository() *fakeGameRepository {
	return &fakeGameRepository{games: make(map[string]*models.Game)}
}

func (r *fakeGameRepository) GetAll(ctx context.Context) ([]*models.Game, error) {
	games := make([]*models.Game, 0, len(r.games))
	for _, game := range r.games {
		games = append(games, game)
	}
	return games, nil
}

func (r *fakeGameRepository) GetByID(ctx context.Context, id string) (*models.Game, error) {
	game, ok := r.games[id]
	if !ok {
		return nil, repository.ErrGameNotFound
	}
	return game, nil
}

func (r *fakeGameRepository) Create(ctx context.Context, game *models.Game) error {
	r.games[game.ID] = game
	return nil
}

func (r *fakeGameRepository) Update(ctx context.Context, id string, game *models.Game) error {
	if _, ok := r.games[id]; !ok {
		return repository.ErrGameNotFound
	}
	game.ID = id
	r.games[id] = game
	return nil
}

func (r *fakeGameRepository) PartialUpdate(ctx context.Context, id string, updates map[string]interface{}) (*models.Game, error) {
	game, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if name, ok := updates["game_name"].(string); ok {
		game.GameName = name
	}
	return game, nil
}

func (r *fakeGameRepository) Delete(ctx context.Context, id string) error {
	if _, ok := r.games[id]; !ok {
		return repository.ErrGameNotFound
	}
	delete(r.games, id)
	return nil
}

func TestCreateGame(t *testing.T) {
	store := newFakeGameRepository()
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

func TestGetAllGames(t *testing.T) {
	store := newFakeGameRepository()
	handler := NewGameHandler(store)

	game := models.Game{
		ID:          "1",
		GameName:    "Test Game",
		Publisher:   "Test Publisher",
		Developer:   "Test Developer",
		ReleaseDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		GameGenre:   "RPG",
	}
	store.Create(context.Background(), &game)

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

func TestDeleteGame(t *testing.T) {
	store := newFakeGameRepository()
	handler := NewGameHandler(store)

	game := models.Game{
		ID:          "1",
		GameName:    "Test Game",
		Publisher:   "Test Publisher",
		Developer:   "Test Developer",
		ReleaseDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		GameGenre:   "RPG",
	}
	store.Create(context.Background(), &game)

	req := httptest.NewRequest("DELETE", "/api/games/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteGame(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	_, err := store.GetByID(context.Background(), "1")
	if err != repository.ErrGameNotFound {
		t.Errorf("Expected game to be deleted")
	}
}

func TestGameNotFound(t *testing.T) {
	store := newFakeGameRepository()
	handler := NewGameHandler(store)

	req := httptest.NewRequest("GET", "/api/games/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.GetGameByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
