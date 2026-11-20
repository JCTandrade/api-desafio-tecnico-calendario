# Holidays API

API HTTP en Go que expone los feriados de Chile obtenidos de `https://api.boostr.cl/holidays.json`, con filtros y un frontend ligero en JavaScript vanilla.

Módulo: `github.com/JCTandrade/api-desafio-tecnico-calendario`

## Arquitectura

Hexagonal (`infrastructure -> application -> domain`):

- `internal/domain`: entidades (`Holiday`, `HolidayFilter`) y `ports` (contratos).
- `internal/application`: caso de uso `HolidayService`.
- `internal/infrastructure`: handlers HTTP, cliente de Boostr y configuración.
- `cmd/api`: composition root.
- `web/`: frontend estático.

## Ejecución

```
go mod tidy
go test ./...
go run ./cmd/api
```

Servidor en `http://localhost:8080`.

## Configuración

| Variable | Defecto |
|---|---|
| `PORT` | `8080` |
| `HOLIDAYS_API_URL` | `https://api.boostr.cl/holidays.json` |
| `HTTP_TIMEOUT` | `10s` |

## Endpoints

- `GET /api/holidays?year=&month=&type=&inalienable=` lista los feriados filtrados. Parámetros inválidos responden `400`; si falla la API externa, `502`.
- `GET /health` responde `{"status":"ok"}`.