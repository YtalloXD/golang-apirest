package routes

import (
	"net/http"

	"github.com/YtalloXD/apirestgo-ia/handlers"
	"github.com/YtalloXD/apirestgo-ia/repository"
	"github.com/gorilla/mux"
)

// SetupRoutes configures all API routes and returns the router
func SetupRoutes(store repository.GameRepository) *mux.Router {
	router := mux.NewRouter()

	// Create handlers
	gameHandler := handlers.NewGameHandler(store)

	// Define routes
	// GET /api/games - retrieve all games
	router.HandleFunc("/api/games", gameHandler.GetAllGames).Methods("GET")

	// POST /api/games - create a new game
	router.HandleFunc("/api/games", gameHandler.CreateGame).Methods("POST")

	// GET /api/games/{id} - retrieve a specific game
	router.HandleFunc("/api/games/{id}", gameHandler.GetGameByID).Methods("GET")

	// PUT /api/games/{id} - fully update an existing game
	router.HandleFunc("/api/games/{id}", gameHandler.UpdateGame).Methods("PUT")

	// PATCH /api/games/{id} - partially update a game
	router.HandleFunc("/api/games/{id}", gameHandler.PartialUpdateGame).Methods("PATCH")

	// DELETE /api/games/{id} - remove a game
	router.HandleFunc("/api/games/{id}", gameHandler.DeleteGame).Methods("DELETE")

	// Health check endpoint
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	return router
}
