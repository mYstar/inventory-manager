# Inventory Manager

REST API server to manage inventory items.

## usage

- `docker compose -f docker/docker-compose.yml up`
- server is accepting requests on port 8080
  - example: `curl localhost:8080/items`

## configuration

- API server port can be configured via `API_PORT` environment variable
  - in `docker/.env`
  - or via: `API_PORT=8080 docker compose -f docker/docker-compose.yml up`

## requests

These are defined in `api/openapi.yaml`. Overview:

- `GET /items`: get all items in the inventory
- `GET /items/value`: calculate the total value of all items in the inventory
- `GET /items/value?id=1&id=2`: calculate the total value the given items
- `POST /item`: add an item to the inventory
  - body: `{"name": "Table", "quantity": 10, "price": 119.99}`
- `PATCH /item/1`: alter the quantity of the item with the given id
  - body: `{"quantity_delta": -5}`
- `DELETE /item/1`: remove the item with the given id from the inventory

## development

TODO