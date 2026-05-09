# API Examples and Usage

This file contains comprehensive examples of how to use the Video Games REST API.

## Prerequisites

- Server running on `http://localhost:8080`
- `curl` command-line tool installed
- `jq` (optional) for formatted JSON output

## Quick Start

### 1. Start the Server

```bash
# In your terminal
go run main.go
```

Expected output:

```
🎮 Video Games API Server starting on http://localhost:8080
📝 Available endpoints:
   GET    /api/games          - Get all games
   POST   /api/games          - Create a new game
   GET    /api/games/{id}     - Get a specific game
   PUT    /api/games/{id}     - Update an entire game
   PATCH  /api/games/{id}     - Partially update a game
   DELETE /api/games/{id}     - Delete a game
   GET    /api/health         - Health check
```

## Detailed Examples

### Health Check

Test if the API is running and healthy.

```bash
curl -X GET http://localhost:8080/api/health

# Response:
# {"status":"healthy"}
```

### GET All Games

Retrieve all games in the system.

```bash
curl -X GET http://localhost:8080/api/games \
  -H "Content-Type: application/json" | jq '.'
```

**Expected Response:**

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
    },
    {
      "id": "2",
      "game_name": "Elden Ring",
      "publisher": "Bandai Namco Entertainment",
      "developer": "FromSoftware",
      "release_date": "2022-02-25T00:00:00Z",
      "game_genre": "Action RPG"
    }
  ]
}
```

### GET Single Game by ID

Retrieve a specific game using its ID.

```bash
curl -X GET http://localhost:8080/api/games/1 \
  -H "Content-Type: application/json" | jq '.'
```

**Expected Response:**

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

### POST - Create a New Game

Create a new game with all required fields. Note: The `id` field is required and must be unique.

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
  }' | jq '.'
```

**Expected Response (201 Created):**

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

### Create Multiple Games

```bash
# Game 6
curl -X POST http://localhost:8080/api/games \
  -H "Content-Type: application/json" \
  -d '{
    "id": "6",
    "game_name": "Stardew Valley",
    "publisher": "ConcernedApe",
    "developer": "ConcernedApe",
    "release_date": "2016-02-26T00:00:00Z",
    "game_genre": "Simulation"
  }'

# Game 7
curl -X POST http://localhost:8080/api/games \
  -H "Content-Type: application/json" \
  -d '{
    "id": "7",
    "game_name": "Celeste",
    "publisher": "Matt Makes Games",
    "developer": "Matt Makes Games",
    "release_date": "2018-01-25T00:00:00Z",
    "game_genre": "Platformer"
  }'
```

### PUT - Fully Update a Game

Replace an entire game. All fields must be provided.

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
  }' | jq '.'
```

**Expected Response (200 OK):**

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

### PATCH - Partially Update a Game

Update only specific fields. Omitted fields remain unchanged.

```bash
# Update just the game name
curl -X PATCH http://localhost:8080/api/games/2 \
  -H "Content-Type: application/json" \
  -d '{
    "game_name": "Elden Ring: Shadow of the Erdtree"
  }' | jq '.'
```

**Expected Response (200 OK):**

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
    "game_genre": "Action RPG"
  }
}
```

### Multiple Field Update with PATCH

```bash
curl -X PATCH http://localhost:8080/api/games/3 \
  -H "Content-Type: application/json" \
  -d '{
    "game_name": "Cyberpunk 2077: Phantom Liberty",
    "publisher": "CD Projekt",
    "game_genre": "Action RPG (Enhanced Edition)"
  }' | jq '.'
```

### DELETE - Remove a Game

Remove a game from the system.

```bash
curl -X DELETE http://localhost:8080/api/games/4
```

**Expected Response (204 No Content):**

```
(Empty body, just the 204 status code)
```

### Verify Game is Deleted

```bash
curl -X GET http://localhost:8080/api/games/4

# Response:
# {"success":false,"error":"Game not found"}
```

## Error Handling Examples

### Invalid JSON

```bash
curl -X POST http://localhost:8080/api/games \
  -H "Content-Type: application/json" \
  -d 'invalid json'
```

**Response (400 Bad Request):**

```json
{
  "success": false,
  "error": "Invalid JSON format"
}
```

### Missing Required Field

```bash
curl -X POST http://localhost:8080/api/games \
  -H "Content-Type: application/json" \
  -d '{
    "id": "8",
    "game_name": "Test Game"
    // missing other required fields
  }'
```

**Response (400 Bad Request):**

```json
{
  "success": false,
  "error": "publisher is required"
}
```

### Game Not Found

```bash
curl -X GET http://localhost:8080/api/games/999
```

**Response (404 Not Found):**

```json
{
  "success": false,
  "error": "Game not found"
}
```

### Update Non-existent Game

```bash
curl -X PUT http://localhost:8080/api/games/999 \
  -H "Content-Type: application/json" \
  -d '{
    "id": "999",
    "game_name": "Ghost Game",
    "publisher": "Publisher",
    "developer": "Developer",
    "release_date": "2023-01-01T00:00:00Z",
    "game_genre": "Mystery"
  }'
```

**Response (404 Not Found):**

```json
{
  "success": false,
  "error": "Game not found"
}
```

## Batch Operations Example

Create a test script to perform multiple operations:

```bash
#!/bin/bash
# File: batch_test.sh

BASE_URL="http://localhost:8080/api"

# Create 3 new games
echo "Creating new games..."
for i in {1..3}; do
  curl -X POST "$BASE_URL/games" \
    -H "Content-Type: application/json" \
    -d "{
      \"id\": \"test_$i\",
      \"game_name\": \"Test Game $i\",
      \"publisher\": \"Test Publisher\",
      \"developer\": \"Test Developer\",
      \"release_date\": \"2023-01-01T00:00:00Z\",
      \"game_genre\": \"Test Genre\"
    }"
  echo ""
done

# Get all games
echo "Fetching all games..."
curl -X GET "$BASE_URL/games" | jq '.data | length'

# Update a game
echo "Updating test_1..."
curl -X PATCH "$BASE_URL/games/test_1" \
  -H "Content-Type: application/json" \
  -d '{"game_name": "Updated Test Game 1"}'

echo ""
```

Run the script:

```bash
chmod +x batch_test.sh
./batch_test.sh
```

## Performance Testing with Apache Bench

```bash
# Install Apache Bench (ab)
# Windows: Download from Apache site
# macOS: brew install httpd
# Linux: sudo apt-get install apache2-utils

# Test GET all games (100 requests, 10 concurrent)
ab -n 100 -c 10 http://localhost:8080/api/games

# Expected output shows:
# - Requests per second
# - Time per request
# - Failed requests
```

## Using Postman

1. Create a new collection
2. Add requests for each endpoint
3. Set up environment variables:
   - `base_url`: http://localhost:8080
   - `game_id`: 1

4. Use in requests:
   - GET `{{base_url}}/api/games`
   - GET `{{base_url}}/api/games/{{game_id}}`
   - POST `{{base_url}}/api/games`

## Using Thunder Client (VS Code)

1. Install Thunder Client extension
2. Click the Thunder Client icon in VS Code
3. Create new requests:
   - Method: GET/POST/PUT/PATCH/DELETE
   - URL: http://localhost:8080/api/games or variants
   - Body: JSON for POST/PUT/PATCH
4. Click Send

## Testing Concurrency

```bash
# Run multiple concurrent requests
for i in {1..10}; do
  curl -X GET http://localhost:8080/api/games &
done
wait
```

This tests the thread-safe implementation of the in-memory store.

## Notes

- All timestamps use ISO 8601 format: `YYYY-MM-DDTHH:MM:SSZ`
- The `id` field is required when creating games
- GET requests return appropriate JSON responses
- DELETE requests return 204 No Content (empty body)
- All errors include meaningful error messages
- The API is thread-safe for concurrent operations
