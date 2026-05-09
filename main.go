package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/YtalloXD/apirestgo-ia/models"
	"github.com/YtalloXD/apirestgo-ia/routes"
	"github.com/YtalloXD/apirestgo-ia/storage"
)

const (
	port = ":8080"
)

func main() {
	// Initialize the in-memory store
	gameStore := storage.NewGameStore()

	// Seed initial data for demonstration
	seedInitialData(gameStore)

	// Setup routes
	router := routes.SetupRoutes(gameStore)

	// Add middleware for logging
	router.Use(loggingMiddleware)

	// Start server
	fmt.Printf("🎮 Video Games API Server starting on http://localhost%s\n", port)
	fmt.Println("📝 Available endpoints:")
	fmt.Println("   GET    /api/games          - Get all games")
	fmt.Println("   POST   /api/games          - Create a new game")
	fmt.Println("   GET    /api/games/{id}     - Get a specific game")
	fmt.Println("   PUT    /api/games/{id}     - Update an entire game")
	fmt.Println("   PATCH  /api/games/{id}     - Partially update a game")
	fmt.Println("   DELETE /api/games/{id}     - Delete a game")
	fmt.Println("   GET    /api/health         - Health check")

	log.Fatal(http.ListenAndServe(port, router))
}

// seedInitialData populates the store with sample data
func seedInitialData(store *storage.GameStore) {
	games := []models.Game{
		{
			ID:          "1",
			GameName:    "The Legend of Zelda: Breath of the Wild",
			Publisher:   "Nintendo",
			Developer:   "Nintendo EPD",
			ReleaseDate: time.Date(2017, 3, 3, 0, 0, 0, 0, time.UTC),
			GameGenre:   "Action-Adventure",
		},
		{
			ID:          "2",
			GameName:    "Elden Ring",
			Publisher:   "Bandai Namco Entertainment",
			Developer:   "FromSoftware",
			ReleaseDate: time.Date(2022, 2, 25, 0, 0, 0, 0, time.UTC),
			GameGenre:   "Action RPG",
		},
		{
			ID:          "3",
			GameName:    "Cyberpunk 2077",
			Publisher:   "CD Projekt",
			Developer:   "CD Projekt Red",
			ReleaseDate: time.Date(2020, 12, 10, 0, 0, 0, 0, time.UTC),
			GameGenre:   "Action RPG",
		},
		{
			ID:          "4",
			GameName:    "Hades",
			Publisher:   "Supergiant Games",
			Developer:   "Supergiant Games",
			ReleaseDate: time.Date(2020, 9, 17, 0, 0, 0, 0, time.UTC),
			GameGenre:   "Roguelike",
		},
	}

	for i := range games {
		if err := store.Create(&games[i]); err != nil {
			log.Printf("Error seeding game %s: %v", games[i].ID, err)
		}
	}

	fmt.Println("✅ Seeded initial data with 4 games")
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
