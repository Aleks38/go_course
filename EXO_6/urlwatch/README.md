# URLWatch

Microservice Go de **vérification d'URLs en masse**. Un client envoie une liste
d'URLs ; le service les interroge **en parallèle** (code HTTP, latence,
disponibilité), **agrège** les résultats et les **expose** via une API REST.
Chaque lot de vérifications (« batch ») est conservé en mémoire et consultable a
posteriori.

## Prérequis

- Go 1.22 ou supérieur (le routage `net/http` et `log/slog` sont utilisés).

## Build / Run / Test

```bash
go build ./...
go vet ./...
go test -race ./...
go run ./cmd/urlwatch
```

Le serveur écoute sur `:8080`.

## Configuration

| Variable | Valeurs | Défaut | Rôle |
|---|---|---|---|
| `LOG_LEVEL` | `debug`, `info`, `warn`, `error` | `info` | Niveau du logger `slog` (sortie JSON). |

```bash
LOG_LEVEL=debug go run ./cmd/urlwatch
```

## API

| Méthode | Chemin | Rôle |
|---|---|---|
| `POST` | `/v1/checks` | Crée et exécute un lot, le persiste, renvoie le résultat. |
| `GET` | `/v1/checks/{id}` | Renvoie un lot existant (ou `404`). |
| `GET` | `/healthz` | Sonde de vivacité (`200`). |

### Créer un lot

```bash
curl -s -X POST http://localhost:8080/v1/checks \
  -H 'Content-Type: application/json' \
  -d '{"urls":["https://go.dev","https://exemple.invalid"],"options":{"concurrency":4,"timeout_ms":2000}}'
```

Options (toutes optionnelles) : `concurrency` (défaut `8`, borné `1..50`),
`timeout_ms` (timeout par URL, défaut `5000`, borné `100..30000`). `urls` est
obligatoire (1 à 100 entrées, chacune en `http`/`https`).

Réponse `201 Created` :

```json
{
  "batch_id": "b_1781885985983253000",
  "created_at": "2026-06-19T15:21:50Z",
  "summary": { "total": 2, "up": 1, "down": 1, "duration_ms": 308 },
  "results": [
    { "url": "https://go.dev", "status_code": 200, "ok": true, "latency_ms": 120 },
    { "url": "https://exemple.invalid", "ok": false, "latency_ms": 13, "error": "dial tcp: lookup exemple.invalid: no such host" }
  ]
}
```

### Relire un lot

```bash
curl -s http://localhost:8080/v1/checks/b_1781885985983253000
```

### Sonde de vivacité

```bash
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080/healthz   # 200
```

### Contrat d'erreur

Toute erreur renvoie ce corps avec le bon code HTTP :

```json
{ "error": { "code": "batch_not_found", "message": "aucun lot avec l'id b_x" } }
```

Codes : `400 invalid_request` (corps ou validation), `404 batch_not_found`,
`405 method_not_allowed`, `500 internal`.

## Structure

```text
cmd/urlwatch/   point d'entrée (assemblage des dépendances + démarrage du serveur)
internal/
  domain/       types métier, erreurs, interfaces (Checker, Store), agrégation
  checker/      implémentation HTTP du Checker
  pool/         cœur concurrent : worker pool, fan-out / fan-in
  store/        persistance en mémoire des lots
  api/          handlers HTTP, routage, middleware de logging, DTO JSON, validation
```
