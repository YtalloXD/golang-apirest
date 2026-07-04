package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/YtalloXD/apirestgo-ia/database"
	"github.com/YtalloXD/apirestgo-ia/models"
	"github.com/YtalloXD/apirestgo-ia/repository"
	"github.com/YtalloXD/apirestgo-ia/routes"
)

const port = ":8080"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := database.LoadConfigFromEnv()
	pool, err := database.NewPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	gameRepository := repository.NewPostgresGameRepository(pool)
	if err := gameRepository.Migrate(ctx); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	if err := seedInitialData(ctx, gameRepository); err != nil {
		log.Printf("seed skipped: %v", err)
	}

	router := routes.SetupRoutes(gameRepository)
	router.Use(loggingMiddleware)

	fmt.Printf("Video Games API Server starting on http://localhost%s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("   GET    /api/games          - Get all games")
	fmt.Println("   POST   /api/games          - Create a new game")
	fmt.Println("   GET    /api/games/{id}     - Get a specific game")
	fmt.Println("   PUT    /api/games/{id}     - Update an entire game")
	fmt.Println("   PATCH  /api/games/{id}     - Partially update a game")
	fmt.Println("   DELETE /api/games/{id}     - Delete a game")
	fmt.Println("   GET    /api/health         - Health check")

	log.Fatal(http.ListenAndServe(port, router))
}

func seedInitialData(ctx context.Context, repo *repository.PostgresGameRepository) error {
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

	return repo.Seed(ctx, games)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
