package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/YtalloXD/apirestgo-ia/models"
	"github.com/YtalloXD/apirestgo-ia/storage"
	"github.com/gorilla/mux"
)

// GameHandler contains all HTTP handlers for game operations
type GameHandler struct {
	store *storage.GameStore
}

// NewGameHandler creates and returns a new GameHandler
func NewGameHandler(store *storage.GameStore) *GameHandler {
	return &GameHandler{store: store}
}

// Response is a generic API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse sends an error response with appropriate status code
func (h *GameHandler) ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
	})
}

// SuccessResponse sends a success response with data
func (h *GameHandler) SuccessResponse(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// GetAllGames handles GET / - retrieve all games
func (h *GameHandler) GetAllGames(w http.ResponseWriter, r *http.Request) {
	games := h.store.GetAll()
	h.SuccessResponse(w, http.StatusOK, games, "Games retrieved successfully")
}

// GetGameByID handles GET /{id} - retrieve a specific game
func (h *GameHandler) GetGameByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	game, err := h.store.GetByID(id)
	if err != nil {
		if err == storage.ErrGameNotFound {
			h.ErrorResponse(w, http.StatusNotFound, "Game not found")
		} else {
			h.ErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	h.SuccessResponse(w, http.StatusOK, game, "Game retrieved successfully")
}

// CreateGame handles POST / - create a new game
func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate required fields
	if err := h.validateGame(&game); err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create the game
	if err := h.store.Create(&game); err != nil {
		h.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SuccessResponse(w, http.StatusCreated, game, "Game created successfully")
}

// UpdateGame handles PUT /{id} - fully update an existing game
func (h *GameHandler) UpdateGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var game models.Game

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate required fields
	if err := h.validateGame(&game); err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update the game
	if err := h.store.Update(id, &game); err != nil {
		if err == storage.ErrGameNotFound {
			h.ErrorResponse(w, http.StatusNotFound, "Game not found")
		} else {
			h.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	game.ID = id
	h.SuccessResponse(w, http.StatusOK, game, "Game updated successfully")
}

// PartialUpdateGame handles PATCH /{id} - partially update a game
func (h *GameHandler) PartialUpdateGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updates map[string]interface{}

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Apply partial updates
	game, err := h.store.PartialUpdate(id, updates)
	if err != nil {
		if err == storage.ErrGameNotFound {
			h.ErrorResponse(w, http.StatusNotFound, "Game not found")
		} else {
			h.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	h.SuccessResponse(w, http.StatusOK, game, "Game updated successfully")
}

// DeleteGame handles DELETE /{id} - remove a game
func (h *GameHandler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.store.Delete(id); err != nil {
		if err == storage.ErrGameNotFound {
			h.ErrorResponse(w, http.StatusNotFound, "Game not found")
		} else {
			h.ErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	h.SuccessResponse(w, http.StatusNoContent, nil, "Game deleted successfully")
}

// validateGame checks if a game has all required fields
func (h *GameHandler) validateGame(game *models.Game) error {
	if game.ID == "" {
		return CustomError("game ID is required")
	}
	if game.GameName == "" {
		return CustomError("game_name is required")
	}
	if game.Publisher == "" {
		return CustomError("publisher is required")
	}
	if game.Developer == "" {
		return CustomError("developer is required")
	}
	if game.ReleaseDate.IsZero() {
		return CustomError("release_date is required")
	}
	if game.GameGenre == "" {
		return CustomError("game_genre is required")
	}
	return nil
}

// CustomError is a simple custom error type
type CustomError string

func (ce CustomError) Error() string {
	return string(ce)
}
