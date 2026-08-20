English · [Русский](README.ru.md)

# warehouse-api

REST API for tracking warehouses, items and stock levels.

## Quick start

```
cp .env.example .env
make up
```

Brings up Postgres and the API in Docker, migrations run automatically on startup. API listens on `:8080`, docs at `http://localhost:8080/swagger/index.html`.

Tear down:

```
make down
```

## API

### Warehouses

| Method | Path | Description |
|---|---|---|
| GET | `/warehouses` | List warehouses |
| GET | `/warehouses/{id}` | Get warehouse by id |
| POST | `/warehouses` | Create warehouse |
| DELETE | `/warehouses/{id}` | Delete warehouse |

### Items

| Method | Path | Description |
|---|---|---|
| GET | `/items` | List items |
| GET | `/items/{id}` | Get item by id |
| POST | `/items` | Create item |

### Stock

| Method | Path | Description |
|---|---|---|
| GET | `/stock` | All stock records |
| GET | `/warehouses/{id}/stock` | Stock for a given warehouse |
| GET | `/items/{id}/stock` | Stock for a given item across all warehouses |
| POST | `/stock` | Create a stock record (warehouse + item + quantity) |
| PATCH | `/stock` | Adjust quantity by a delta (restock or write-off, can be negative) |
| POST | `/stock/transfer` | Move stock between two warehouses in one transaction |
