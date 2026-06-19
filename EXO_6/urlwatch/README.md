# URLWatch

Microservice Go de **vérification d'URLs en masse**. Un client envoie une liste
d'URLs ; le service les interroge **en parallèle** (code HTTP, latence,
disponibilité), **agrège** les résultats et les **expose** via une API REST.
Chaque lot de vérifications (« batch ») est conservé et consultable a posteriori.

> Projet d'examen — en cours de construction. Cette page sera complétée
> (build / run / test + exemples `curl`) au fur et à mesure.

## Prérequis

- Go 1.22 ou supérieur.

## Build / Run / Test

```bash
go build ./...
go vet ./...
go test ./...           # (idéalement : go test -race ./...)
go run ./cmd/urlwatch
```

## Structure

```text
cmd/urlwatch/   point d'entrée (assemblage + démarrage du serveur)
internal/
  domain/       types métier, erreurs, interfaces (Checker, Store)
  checker/      implémentation HTTP du Checker (+ mock pour les tests)
  pool/         cœur concurrent : worker pool, fan-out / fan-in
  store/        persistance des lots (mémoire ; SQLite en bonus)
  api/          handlers HTTP, routage, middleware, DTO JSON
```

## API

À documenter (endpoints `POST /v1/checks`, `GET /v1/checks/{id}`, `GET /healthz`
+ exemples `curl`).
