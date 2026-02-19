# Inventory Manager

REST API server to manage inventory items.

## Requirements

- [Docker Engine 17.05+](https://docs.docker.com/get-docker/)
- [Docker Compose v2](https://docs.docker.com/compose/install/)

## Usage

- `docker compose -f docker/docker-compose.yml up`
- server is accepting requests on port 8080
  - example: `curl localhost:8080/items`

### Testing

- to run all tests just use: `go test ./...`

## Configuration

- `API_PORT`: API server port (default: `8080`)
  - set in `docker/.env`, or:
  - `API_PORT=8081 docker compose -f docker/docker-compose.yml up`
- `DB_FILE`: SQLite database filename (default: `default.db`)
  - in `docker/.env`
  - or via: `DB_FILE=inventory.db docker compose -f docker/docker-compose.yml up`

## Requests

These are defined in `api/openapi.yaml`. Overview:

- `GET /items`: get all items in the inventory
- `GET /items/value`: calculate the total value of all items in the inventory
- `GET /items/value?ids=1&ids=2`: calculate the total value the given items
- `POST /item`: add an item to the inventory
  - body: `{"name": "Table", "quantity": 10, "price_cents": 11999}`
- `PATCH /item/1`: alter the quantity of the item with the given id
  - body: `{"quantity_delta": -5}`
- `DELETE /item/1`: remove the item with the given id from the inventory
