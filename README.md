# API-RiftRadar

API-RiftRadar is a Go-based web service that provides information about League of Legends players and leaderboards. It allows users to retrieve player data, create and manage leaderboards, and add players to leaderboards.

## Project Overview

This project follows a hexagonal architecture (ports and adapters) pattern to maintain separation of concerns and facilitate testing. It integrates with the Riot Games API to fetch player data while also maintaining local in-memory storage.

### Key Features

- Retrieve player information by name and tag
- List all players
- Create and manage leaderboards
- Add players to leaderboards
- RESTful API design

## Architecture

The project is structured according to hexagonal architecture principles:

- **Domain**: Contains the core business entities and port interfaces
- **Adapters**: 
  - **In**: HTTP handlers for incoming requests
  - **Out**: Implementations for external services (memory storage, Riot API)
- **Application**: Contains the business logic and use cases

## Setup Instructions

### Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose (optional, for containerized setup)
- Riot Games API key (for fetching player data from Riot API)

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
RIOT_API_KEY=your_riot_api_key_here
```

### Running Locally

1. Clone the repository:
   ```
   git clone https://github.com/tyryoxan/API-RiftRadar.git
   cd API-RiftRadar
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Run the application:
   ```
   go run main.go
   ```

The server will start on port 8080.

### Running with Docker

1. Clone the repository:
   ```
   git clone https://github.com/tyryoxan/API-RiftRadar.git
   cd API-RiftRadar
   ```

2. Create a `.env` file with your Riot API key

3. Start the application with Docker Compose:
   ```
   docker-compose up
   ```

The server will be available on port 8080.

## API Documentation

This section provides detailed information about the available API endpoints, including request parameters, response formats, and status codes.

### Player Endpoints

#### `GET /players`

Retrieves a list of all players in the system.

**Request:**
- Method: `GET`
- URL: `/players`
- Headers: None required

**Response:**
- Status Code: `200 OK`
- Content-Type: `application/json`
- Body: Array of player objects

**Example Response:**
```json
[
  {
    "name": "PlayerOne",
    "region": "EUW",
    "elo": "Diamond",
    "lp": 75
  },
  {
    "name": "PlayerTwo",
    "region": "NA",
    "elo": "Gold",
    "lp": 30
  }
]
```

**Error Responses:**
- None specific to this endpoint

---

#### `GET /players/{name}/{tag}`

Retrieves information about a specific player by their name and tag.

**Request:**
- Method: `GET`
- URL: `/players/{name}/{tag}`
- URL Parameters:
  - `name` (required): The player's name
  - `tag` (required): The player's tag

**Response:**
- Status Code: `200 OK`
- Content-Type: `application/json`
- Body: Player object

**Example Request:**
```
GET /players/PlayerOne/EUW1
```

**Example Response:**
```json
{
  "name": "PlayerOne",
  "region": "EUW",
  "elo": "Diamond",
  "lp": 75
}
```

**Error Responses:**
- `400 Bad Request`: If name or tag parameters are missing
  ```json
  {
    "error": "name and tag parameters are required"
  }
  ```
- `500 Internal Server Error`: If there's an error retrieving the player
  ```json
  {
    "error": "player not found"
  }
  ```

---

### Leaderboard Endpoints

#### `GET /leaderboard`

Retrieves the current leaderboard with all players.

**Request:**
- Method: `GET`
- URL: `/leaderboard`
- Headers: None required

**Response:**
- Status Code: `200 OK`
- Content-Type: `application/json`
- Body: Leaderboard object with players

**Example Response:**
```json
{
  "name": "Global Leaderboard",
  "players": [
    {
      "name": "PlayerOne",
      "region": "EUW",
      "elo": "Diamond",
      "lp": 75
    },
    {
      "name": "PlayerTwo",
      "region": "NA",
      "elo": "Gold",
      "lp": 30
    }
  ]
}
```

**Error Responses:**
- None specific to this endpoint

---

#### `POST /leaderboard`

Creates a new leaderboard.

**Request:**
- Method: `POST`
- URL: `/leaderboard`
- Headers:
  - Content-Type: `application/json`
- Body: JSON object with leaderboard name

**Example Request:**
```json
{
  "name": "Custom Leaderboard"
}
```

**Response:**
- Status Code: `200 OK`
- Content-Type: `application/json`
- Body: Created leaderboard object

**Example Response:**
```json
{
  "name": "Custom Leaderboard",
  "players": []
}
```

**Error Responses:**
- `400 Bad Request`: If the request payload is invalid
  ```json
  {
    "error": "Invalid request payload"
  }
  ```

---

#### `POST /leaderboard/player`

Adds a player to the leaderboard.

**Request:**
- Method: `POST`
- URL: `/leaderboard/player`
- Headers:
  - Content-Type: `application/json`
- Body: JSON object with player name and tag

**Example Request:**
```json
{
  "name": "PlayerOne",
  "tag": "EUW1"
}
```

**Response:**
- Status Code: `200 OK`
- Content-Type: `application/json`
- Body: Success message

**Example Response:**
```json
{
  "message": "Player added successfully"
}
```

**Error Responses:**
- `400 Bad Request`: If the request payload is invalid
  ```json
  {
    "error": "Invalid request payload"
  }
  ```
- `500 Internal Server Error`: If there's an error adding the player
  ```json
  {
    "error": "player not found"
  }
  ```

## Data Models

### Player

```json
{
  "name": "string",
  "region": "string",
  "elo": "string",
  "lp": 0
}
```

### Leaderboard

```json
{
  "name": "string",
  "players": [
    {
      "name": "string",
      "region": "string",
      "elo": "string",
      "lp": 0
    }
  ]
}
```

## Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux): HTTP router and URL matcher

## License

This project is licensed under the MIT License - see the LICENSE file for details.
