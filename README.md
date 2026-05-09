# 🎮 Video Games REST API - Go

A production-grade RESTful API system for managing video games built with Go (Golang). This implementation demonstrates idiomatic Go practices, clean architecture, and comprehensive CRUD operations.

## Architecture Overview

### Design Principles

This API follows several key architectural principles:

1. **Separation of Concerns**: Code is organized into distinct packages:
   - `models/`: Data structures and domain objects
   - `handlers/`: HTTP request/response handling
   - `storage/`: Data persistence layer (in-memory)
   - `routes/`: Route definitions and server setup

2. **Thread Safety**: The in-memory store uses `sync.RWMutex` for concurrent access protection

3. **RESTful Design**: Proper use of HTTP methods and status codes
   - GET for retrieval
   - POST for creation (201 Created)
   - PUT for full updates (200 OK)
   - PATCH for partial updates (200 OK)
   - DELETE for removal (204 No Content)

4. **Error Handling**: Comprehensive error handling with meaningful error messages

5. **JSON Serialization**: Automatic marshaling/unmarshaling with proper tag annotations

## Project Structure

```
apirestgo-ia/
├── main.go                 # Entry point, server startup, data seeding
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums (auto-generated)
│
├── models/
│   └── game.go            # Game struct with JSON tags
│
├── handlers/
│   └── game_handler.go    # HTTP handlers for all CRUD operations
│
├── storage/
│   └── game_store.go      # In-memory data store with thread safety
│
├── routes/
│   └── routes.go          # Route definitions and router setup
│
└── README.md              # This file
```

## Data Model

The `Game` struct represents a video game with the following fields:

```go
type Game struct {
    ID          string    `json:"id"`
    GameName    string    `json:"game_name"`
    Publisher   string    `json:"publisher"`
    Developer   string    `json:"developer"`
    ReleaseDate time.Time `json:"release_date"`
    GameGenre   string    `json:"game_genre"`
}
```

## API Endpoints

### 1. Get All Games

**GET** `/api/games`

Retrieve all games in the system.

**Response Example:**

```json
{
  "success": true,
  "message": "Games retrieved successfully",
  "data": [
    {
      "id": "1",
      "game_name": "The Legend of Zelda: Breath of the Wild",
      "publisher": "Nintendo",
      "developer": "Nintendo EPD",
      "release_date": "2017-03-03T00:00:00Z",
      "game_genre": "Action-Adventure"
    }
  ]
}
```

### 2. Get Game by ID

**GET** `/api/games/{id}`

Retrieve a specific game by its ID.

**Response Example:**

```json
{
  "success": true,
  "message": "Game retrieved successfully",
  "data": {
    "id": "1",
    "game_name": "The Legend of Zelda: Breath of the Wild",
    "publisher": "Nintendo",
    "developer": "Nintendo EPD",
    "release_date": "2017-03-03T00:00:00Z",
    "game_genre": "Action-Adventure"
  }
}
```

### 3. Create a New Game

**POST** `/api/games`

Create a new game. All fields are required.

**Request Body:**

```json
{
  "id": "5",
  "game_name": "Hollow Knight",
  "publisher": "Team Cherry",
  "developer": "Team Cherry",
  "release_date": "2017-02-24T00:00:00Z",
  "game_genre": "Metroidvania"
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "message": "Game created successfully",
  "data": {
    "id": "5",
    "game_name": "Hollow Knight",
    "publisher": "Team Cherry",
    "developer": "Team Cherry",
    "release_date": "2017-02-24T00:00:00Z",
    "game_genre": "Metroidvania"
  }
}
```

### 4. Fully Update a Game (PUT)

**PUT** `/api/games/{id}`

Completely replace a game. All fields must be provided.

**Request Body:**

```json
{
  "id": "1",
  "game_name": "The Legend of Zelda: Tears of the Kingdom",
  "publisher": "Nintendo",
  "developer": "Nintendo EPD",
  "release_date": "2023-05-12T00:00:00Z",
  "game_genre": "Action-Adventure"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Game updated successfully",
  "data": {
    "id": "1",
    "game_name": "The Legend of Zelda: Tears of the Kingdom",
    "publisher": "Nintendo",
    "developer": "Nintendo EPD",
    "release_date": "2023-05-12T00:00:00Z",
    "game_genre": "Action-Adventure"
  }
}
```

### 5. Partially Update a Game (PATCH)

**PATCH** `/api/games/{id}`

Update specific fields only. Omitted fields are not changed.

**Request Body:**

```json
{
  "game_name": "Elden Ring: Shadow of the Erdtree",
  "game_genre": "Action RPG - DLC"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Game updated successfully",
  "data": {
    "id": "2",
    "game_name": "Elden Ring: Shadow of the Erdtree",
    "publisher": "Bandai Namco Entertainment",
    "developer": "FromSoftware",
    "release_date": "2022-02-25T00:00:00Z",
    "game_genre": "Action RPG - DLC"
  }
}
```

### 6. Delete a Game

**DELETE** `/api/games/{id}`

Remove a game from the system.

**Response (204 No Content):** Empty body with 204 status code

### 7. Health Check

**GET** `/api/health`

System health check endpoint.

**Response:**

```json
{
  "status": "healthy"
}
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- curl (for testing) or Postman/Thunder Client

### Installation

1. **Clone/Navigate to the project:**

   ```bash
   cd c:\Program Files\VSCode projects\apirestgo-ia
   ```

2. **Download dependencies:**

   ```bash
   go mod download
   ```

3. **Build the application:**

   ```bash
   go build -o video-games-api.exe
   ```

4. **Run the server:**

   ```bash
   go run main.go
   ```

   Or run the compiled binary:

   ```bash
   .\video-games-api.exe
   ```

The server will start on `http://localhost:8080` with 4 pre-seeded games.

## Example API Requests using cURL

### Get All Games

```bash
curl -X GET http://localhost:8080/api/games \
  -H "Content-Type: application/json"
```

### Get Single Game

```bash
curl -X GET http://localhost:8080/api/games/1 \
  -H "Content-Type: application/json"
```

### Create a New Game

```bash
curl -X POST http://localhost:8080/api/games \
  -H "Content-Type: application/json" \
  -d '{
    "id": "5",
    "game_name": "Hollow Knight",
    "publisher": "Team Cherry",
    "developer": "Team Cherry",
    "release_date": "2017-02-24T00:00:00Z",
    "game_genre": "Metroidvania"
  }'
```

### Update Entire Game (PUT)

```bash
curl -X PUT http://localhost:8080/api/games/1 \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1",
    "game_name": "The Legend of Zelda: Tears of the Kingdom",
    "publisher": "Nintendo",
    "developer": "Nintendo EPD",
    "release_date": "2023-05-12T00:00:00Z",
    "game_genre": "Action-Adventure"
  }'
```

### Partially Update Game (PATCH)

```bash
curl -X PATCH http://localhost:8080/api/games/2 \
  -H "Content-Type: application/json" \
  -d '{
    "game_name": "Elden Ring: Enhanced Edition",
    "game_genre": "Action RPG - Enhanced"
  }'
```

### Delete a Game

```bash
curl -X DELETE http://localhost:8080/api/games/3
```

### Health Check

```bash
curl -X GET http://localhost:8080/api/health
```

## Key Implementation Features

### 1. Thread-Safe Storage

The `GameStore` uses `sync.RWMutex` for concurrent access:

```go
type GameStore struct {
    mu    sync.RWMutex
    games map[string]*models.Game
}
```

### 2. Proper HTTP Status Codes

- **200 OK**: Successful GET, PUT, PATCH
- **201 Created**: Successful POST
- **204 No Content**: Successful DELETE
- **400 Bad Request**: Invalid input or malformed JSON
- **404 Not Found**: Resource doesn't exist
- **500 Internal Server Error**: Server error

### 3. Validation

Each handler validates required fields before processing:

```go
func (h *GameHandler) validateGame(game *models.Game) error {
    // Comprehensive field validation
}
```

### 4. JSON Tags

Proper JSON tag annotations for serialization/deserialization:

```go
type Game struct {
    ID          string    `json:"id"`
    GameName    string    `json:"game_name"`
    // ...
}
```

### 5. Error Handling

Consistent error response format:

```json
{
  "success": false,
  "error": "Game not found"
}
```

## Go Idioms Applied

1. **Interface-based design** for flexibility and testability
2. **Pointer receivers** for methods that modify state
3. **Error as return value** instead of exceptions
4. **Composition over inheritance** with struct embedding
5. **Defer** for resource cleanup
6. **Goroutine-safe** synchronization with mutexes
7. **Unexported fields** (lowercase) with exported accessors
8. **Clear naming** following Go conventions

## Performance Considerations

- **In-memory storage** provides O(1) access for most operations
- **RWMutex** allows multiple concurrent readers
- **JSON streaming** with `json.Decoder` and `json.Encoder` for efficiency

## Testing the API

### Using Thunder Client (VS Code Extension)

1. Install the Thunder Client extension
2. Create requests for each endpoint
3. Save collections for future use

### Using Postman

1. Import the API endpoints into Postman
2. Create a collection with all CRUD operations
3. Test with various payloads

### Automated Testing

Add unit tests in `*_test.go` files:

```bash
go test ./...
```

## Extending the API

### Add Filtering

```go
// GET /api/games?genre=RPG
router.HandleFunc("/api/games", gameHandler.GetGamesByGenre).Methods("GET")
```

### Add Pagination

```go
// GET /api/games?page=1&limit=10
type PaginationParams struct {
    Page  int
    Limit int
}
```

### Add Database

Replace in-memory storage with PostgreSQL/MongoDB:

```bash
go get github.com/lib/pq          // PostgreSQL
go get go.mongodb.org/mongo-driver // MongoDB
```

## Production Considerations

1. **Error Logging**: Implement structured logging
2. **Request Validation**: Add input sanitization
3. **CORS**: Add middleware for cross-origin requests
4. **Rate Limiting**: Implement request throttling
5. **Authentication**: Add JWT or API key validation
6. **Database**: Replace in-memory storage with persistent DB
7. **Docker**: Containerize the application
8. **CI/CD**: Set up automated testing and deployment

## Project Highlights

✅ **Clean Architecture**: Well-organized package structure  
✅ **Idiomatic Go**: Follows Go best practices and conventions  
✅ **Thread-Safe**: Concurrent request handling with proper synchronization  
✅ **Comprehensive**: Full CRUD operations with proper HTTP semantics  
✅ **Error Handling**: Detailed error messages and status codes  
✅ **Readable Code**: Clear naming and meaningful comments  
✅ **Scalable**: Designed for easy extension and testing

## License

This project is open source and available under the MIT License.

## Author

Created as a comprehensive demonstration of RESTful API development in Go.
