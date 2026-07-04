package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/YtalloXD/apirestgo-ia/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrGameNotFound = errors.New("game not found")
	ErrEmptyID      = errors.New("game id cannot be empty")
)

type GameRepository interface {
	GetAll(ctx context.Context) ([]*models.Game, error)
	GetByID(ctx context.Context, id string) (*models.Game, error)
	Create(ctx context.Context, game *models.Game) error
	Update(ctx context.Context, id string, game *models.Game) error
	PartialUpdate(ctx context.Context, id string, updates map[string]interface{}) (*models.Game, error)
	Delete(ctx context.Context, id string) error
}

type PostgresGameRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresGameRepository(pool *pgxpool.Pool) *PostgresGameRepository {
	return &PostgresGameRepository{pool: pool}
}

func (r *PostgresGameRepository) Migrate(ctx context.Context) error {
	const query = `
		CREATE TABLE IF NOT EXISTS games (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			publisher TEXT NOT NULL,
			developer TEXT NOT NULL,
			release_date DATE NOT NULL,
			genre TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`

	if _, err := r.pool.Exec(ctx, query); err != nil {
		return fmt.Errorf("create games table: %w", err)
	}
	return nil
}

func (r *PostgresGameRepository) GetAll(ctx context.Context) ([]*models.Game, error) {
	const query = `
		SELECT id, name, publisher, developer, release_date, genre
		FROM games
		ORDER BY release_date DESC, name ASC;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select games: %w", err)
	}
	defer rows.Close()

	games := make([]*models.Game, 0)
	for rows.Next() {
		game, err := scanGame(rows)
		if err != nil {
			return nil, fmt.Errorf("scan game: %w", err)
		}
		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate games: %w", err)
	}

	return games, nil
}

func (r *PostgresGameRepository) GetByID(ctx context.Context, id string) (*models.Game, error) {
	if id == "" {
		return nil, ErrEmptyID
	}

	const query = `
		SELECT id, name, publisher, developer, release_date, genre
		FROM games
		WHERE id = $1;
	`

	game, err := scanGame(r.pool.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrGameNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("select game by id: %w", err)
	}
	return game, nil
}

func (r *PostgresGameRepository) Create(ctx context.Context, game *models.Game) error {
	if game.ID == "" {
		return ErrEmptyID
	}

	const query = `
		INSERT INTO games (id, name, publisher, developer, release_date, genre)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := r.pool.Exec(ctx, query, game.ID, game.GameName, game.Publisher, game.Developer, game.ReleaseDate, game.GameGenre)
	if err != nil {
		return fmt.Errorf("insert game: %w", err)
	}
	return nil
}

func (r *PostgresGameRepository) Update(ctx context.Context, id string, game *models.Game) error {
	if id == "" {
		return ErrEmptyID
	}

	const query = `
		UPDATE games
		SET name = $2, publisher = $3, developer = $4, release_date = $5, genre = $6, updated_at = NOW()
		WHERE id = $1;
	`

	tag, err := r.pool.Exec(ctx, query, id, game.GameName, game.Publisher, game.Developer, game.ReleaseDate, game.GameGenre)
	if err != nil {
		return fmt.Errorf("update game: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrGameNotFound
	}

	game.ID = id
	return nil
}

func (r *PostgresGameRepository) PartialUpdate(ctx context.Context, id string, updates map[string]interface{}) (*models.Game, error) {
	if id == "" {
		return nil, ErrEmptyID
	}

	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name, ok := updates["game_name"].(string); ok {
		current.GameName = name
	}
	if publisher, ok := updates["publisher"].(string); ok {
		current.Publisher = publisher
	}
	if developer, ok := updates["developer"].(string); ok {
		current.Developer = developer
	}
	if genre, ok := updates["game_genre"].(string); ok {
		current.GameGenre = genre
	}
	if value, ok := updates["release_date"].(string); ok {
		releaseDate, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, fmt.Errorf("invalid release_date: %w", err)
		}
		current.ReleaseDate = releaseDate
	}

	if err := r.Update(ctx, id, current); err != nil {
		return nil, err
	}
	return current, nil
}

func (r *PostgresGameRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrEmptyID
	}

	tag, err := r.pool.Exec(ctx, `DELETE FROM games WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("delete game: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrGameNotFound
	}
	return nil
}

func (r *PostgresGameRepository) Seed(ctx context.Context, games []models.Game) error {
	const query = `
		INSERT INTO games (id, name, publisher, developer, release_date, genre)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO NOTHING;
	`

	batch := &pgx.Batch{}
	for i := range games {
		batch.Queue(query, games[i].ID, games[i].GameName, games[i].Publisher, games[i].Developer, games[i].ReleaseDate, games[i].GameGenre)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for range games {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("seed games: %w", err)
		}
	}
	return nil
}

type gameScanner interface {
	Scan(dest ...interface{}) error
}

func scanGame(row gameScanner) (*models.Game, error) {
	var game models.Game
	if err := row.Scan(&game.ID, &game.GameName, &game.Publisher, &game.Developer, &game.ReleaseDate, &game.GameGenre); err != nil {
		return nil, err
	}
	return &game, nil
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
