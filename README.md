# Chile Holidays API

API en Go (Gin) que expone los feriados de Chile con filtros por tipo y rango de fechas, en JSON o XML.

- Los feriados se obtienen de <https://api.boostr.cl/holidays.json> **una sola vez, al iniciar la aplicación**, y se mantienen en memoria. No se vuelve a llamar al servicio por cada request. Si la carga inicial falla, la aplicación no arranca.
- Arquitectura hexagonal, logs estructurados (`log/slog`, JSON), tests unitarios y especificación OpenAPI 3.

## Requisitos

- Go 1.26, o
- Docker Engine 19.03.5+ (y docker-compose 1.25+) para ejecutarlo en contenedor.

## Ejecución

### Docker Compose

```bash
docker-compose up --build -d     # con Docker Compose v2: docker compose up --build -d
curl http://localhost:8080/health
```

### Docker

```bash
docker build -t holidays-api .
docker run --rm -p 8080:8080 holidays-api
```

### Local

```bash
go run ./cmd/api
```

### Configuración (variables de entorno)

| Variable           | Por defecto                              | Descripción                          |
|--------------------|------------------------------------------|--------------------------------------|
| `PORT`             | `8080`                                   | Puerto HTTP                          |
| `HOLIDAYS_API_URL` | `https://api.boostr.cl/holidays.json`    | Servicio base de feriados            |
| `HTTP_TIMEOUT`     | `10s`                                    | Timeout hacia el servicio base       |
| `GIN_MODE`         | `debug` (`release` en Docker)            | Modo de Gin                          |

## Endpoints

| Método | Ruta             | Descripción                                   |
|--------|------------------|-----------------------------------------------|
| GET    | `/api/holidays`  | Lista feriados filtrados (JSON o XML)         |
| GET    | `/health`        | Health check                                  |
| GET    | `/openapi.yaml`  | Especificación OpenAPI 3 ([api/openapi.yaml](api/openapi.yaml)) |
| GET    | `/docs`          | Swagger UI                                    |

### Parámetros de `/api/holidays` (todos opcionales)

| Parámetro     | Ejemplo        | Descripción                                                        |
|---------------|----------------|--------------------------------------------------------------------|
| `type`        | `Religioso`    | Tipo de feriado (sin distinguir mayúsculas), p. ej. `Civil`, `Religioso` |
| `from`        | `2026-01-01`   | Inicio del rango de fechas (inclusivo, `YYYY-MM-DD`)               |
| `to`          | `2026-06-30`   | Fin del rango de fechas (inclusivo, `YYYY-MM-DD`)                  |
| `year`        | `2026`         | Año                                                                |
| `month`       | `5`            | Mes (1-12)                                                         |
| `inalienable` | `true`         | Filtra por irrenunciable                                           |
| `format`      | `json` / `xml` | Formato de respuesta; tiene prioridad sobre el header `Accept`     |

El formato se negocia con el header `Accept` (`application/json` por defecto, `application/xml`) o con `format`.

### Campos de la respuesta

`date` (fecha), `title` (título), `type` (tipo), `inalienable` (irrenunciable) y `extra`.

> El enunciado menciona "Teléfono", pero el servicio base de boostr no entrega ese dato (sus campos son `date`, `title`, `type`, `inalienable` y `extra`), por lo que no se incluye.

## Ejemplos (curl)

```bash
# Todos los feriados (JSON)
curl "http://localhost:8080/api/holidays"

# Solo religiosos
curl "http://localhost:8080/api/holidays?type=Religioso"

# Civiles entre dos fechas
curl "http://localhost:8080/api/holidays?type=Civil&from=2026-05-01&to=2026-06-30"

# Irrenunciables de un año
curl "http://localhost:8080/api/holidays?year=2026&inalienable=true"

# Respuesta en XML (header Accept)
curl -H "Accept: application/xml" "http://localhost:8080/api/holidays?from=2026-05-01&to=2026-05-31"

# Respuesta en XML (parámetro format)
curl "http://localhost:8080/api/holidays?format=xml"

# Respuesta en JSON explícita
curl -H "Accept: application/json" "http://localhost:8080/api/holidays?month=9"

# Error 400: rango invertido
curl -i "http://localhost:8080/api/holidays?from=2026-06-01&to=2026-01-01"

# Health check y especificación
curl http://localhost:8080/health
curl http://localhost:8080/openapi.yaml
```

Ejemplo de respuesta JSON:

```json
[{"date":"2026-05-01","title":"Día Nacional del Trabajo","type":"Civil","inalienable":true,"extra":"Civil e Irrenunciable"}]
```

Ejemplo de respuesta XML:

```xml
<holidays><holiday><date>2026-05-01</date><title>Día Nacional del Trabajo</title><type>Civil</type><inalienable>true</inalienable><extra>Civil e Irrenunciable</extra></holiday></holidays>
```

Códigos de error: `400` parámetro inválido o rango invertido, `500` error interno.

## Arquitectura hexagonal

```
cmd/api                          composición (main): conecta puertos y adaptadores
api/                             especificación OpenAPI embebida
internal/
  domain/                        entidades y reglas (Holiday, HolidayFilter) sin dependencias externas
  domain/ports/                  puertos: HolidayService (entrada) y HolidayRepository (salida)
  application/                   casos de uso (HolidayService)
  infrastructure/
    http/                        adaptador de entrada: Gin, negociación JSON/XML, logging de requests
    repository/                  adaptadores de salida: cliente boostr y repositorio en memoria (carga única)
    config/                      configuración por variables de entorno
```

El dominio y la aplicación no dependen de la infraestructura; la dependencia siempre apunta hacia adentro.

## Tests

```bash
go test ./...
```
