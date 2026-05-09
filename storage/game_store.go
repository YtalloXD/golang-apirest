package storage

import (
	"errors"
	"sync"

	"github.com/YtalloXD/apirestgo-ia/models"
)

var (
	ErrGameNotFound = errors.New("game not found")
	ErrEmptyID      = errors.New("game id cannot be empty")
)

// GameStore manages in-memory storage of games with thread-safe operations
type GameStore struct {
	mu    sync.RWMutex
	games map[string]*models.Game
}

// NewGameStore creates and returns a new GameStore instance
func NewGameStore() *GameStore {
	return &GameStore{
		games: make(map[string]*models.Game),
	}
}

// GetAll returns all games in the store
func (gs *GameStore) GetAll() []*models.Game {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	games := make([]*models.Game, 0, len(gs.games))
	for _, game := range gs.games {
		games = append(games, game)
	}
	return games
}

// GetByID retrieves a game by its ID
func (gs *GameStore) GetByID(id string) (*models.Game, error) {
	if id == "" {
		return nil, ErrEmptyID
	}

	gs.mu.RLock()
	defer gs.mu.RUnlock()

	game, exists := gs.games[id]
	if !exists {
		return nil, ErrGameNotFound
	}
	return game, nil
}

// Create adds a new game to the store
func (gs *GameStore) Create(game *models.Game) error {
	if game.ID == "" {
		return ErrEmptyID
	}

	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.games[game.ID] = game
	return nil
}

// Update replaces an entire game (PUT operation)
func (gs *GameStore) Update(id string, game *models.Game) error {
	if id == "" {
		return ErrEmptyID
	}

	gs.mu.Lock()
	defer gs.mu.Unlock()

	if _, exists := gs.games[id]; !exists {
		return ErrGameNotFound
	}

	game.ID = id // Ensure ID consistency
	gs.games[id] = game
	return nil
}

// PartialUpdate updates specific fields of a game (PATCH operation)
func (gs *GameStore) PartialUpdate(id string, updates map[string]interface{}) (*models.Game, error) {
	if id == "" {
		return nil, ErrEmptyID
	}

	gs.mu.Lock()
	defer gs.mu.Unlock()

	game, exists := gs.games[id]
	if !exists {
		return nil, ErrGameNotFound
	}

	// Apply partial updates
	if name, ok := updates["game_name"].(string); ok {
		game.GameName = name
	}
	if publisher, ok := updates["publisher"].(string); ok {
		game.Publisher = publisher
	}
	if developer, ok := updates["developer"].(string); ok {
		game.Developer = developer
	}
	if genre, ok := updates["game_genre"].(string); ok {
		game.GameGenre = genre
	}

	return game, nil
}

// Delete removes a game from the store
func (gs *GameStore) Delete(id string) error {
	if id == "" {
		return ErrEmptyID
	}

	gs.mu.Lock()
	defer gs.mu.Unlock()

	if _, exists := gs.games[id]; !exists {
		return ErrGameNotFound
	}

	delete(gs.games, id)
	return nil
}
