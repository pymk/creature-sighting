> [!NOTE]
> This over-engineered project serves as a comprehensive learning exercise exploring development workflows with **Claude Code**. It was built purely to learn development using Claude Code. Rather than writing code directly, I focused entirely on orchestrating Claude Code through prompting and project management.

---

# Creature Sighting: Fictional Creature Sighting Generator

A Go server that generates random fictional creature sightings around the world. Currently supports Kaiju sightings with an extensible architecture for adding more creature types.

Features both a web interface for viewing sightings and a REST API for programmatic access.

## Running the Server

```bash
make run
```

The server will start on port 8080.

## Web Interface

Visit `http://localhost:8080` to access the web interface:

- **Home** (`/`) - Welcome page with navigation and stats
- **Sightings** (`/sightings`) - Grid view of all creature sightings
- **Sighting Details** (`/sighting/{id}`) - Detailed view of individual sightings
- **Random Sighting** (`/sighting/random`) - Generate and view new sightings

The web interface uses minimal CSS styling and requires no JavaScript.

## API Endpoints

The REST API is available under `/api/` routes:

### Generate a Sighting
```bash
GET /api/sighting?category=kaiju
```

Example response:
```json
{
  "id": "kaiju-1736114523456789",
  "name": "Gorgozilla",
  "type": "Aquatic",
  "category": "kaiju",
  "location": {
    "latitude": 35.6762,
    "longitude": 139.6503,
    "city": "Tokyo",
    "country": "Japan",
    "region": "Asia"
  },
  "description": "A colossal Aquatic kaiju displaying aggressive behavior",
  "timestamp": "2025-01-04T15:55:23Z",
  "attributes": {
    "behavior": "aggressive",
    "height": "175 meters",
    "size": "colossal"
  }
}
```

### List Available Categories
```bash
GET /api/categories
```

Example response:
```json
{
  "categories": ["kaiju"]
}
```

## Adding New Creature Types

1. Create a new package in `internal/creatures/`
2. Implement the `sighting.Generator` interface
3. Register the generator in `cmd/server/main.go`

Example:
```go
type Generator interface {
    Generate() (*Sighting, error)
    Category() string
}
```

## Development

```bash
# Run tests
make test

# Format code
make fmt

# Run linters
make lint

# Clean build artifacts
make clean
```
