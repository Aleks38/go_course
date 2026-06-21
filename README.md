# Mastère – Module Go

Dépôt de cours pour le module Go. Chaque dossier `src/EXO_N/` correspond à un
chapitre (leçons `.md` + exercices `main_*.go`).

## Exercices notés

Voici où trouver les exercices à évaluer, avec la commande pour les lancer.

| Exercice | Emplacement | Lancer |
|---|---|---|
| **4a** – Concurrence (goroutines, channels, worker pool) | `src/EXO_4/main_4a.go` | `go run src/EXO_4/main_4a.go` |
| **5a** – API REST avec `net/http` (CRUD en mémoire, JSON) | `src/EXO_5/main_5a.go` | `go run src/EXO_5/main_5a.go` |
| **5d** – Accès base de données (`database/sql`, `sqlx`, `GORM`) | `src/EXO_5/go-db-tp/` | `cd src/EXO_5/go-db-tp && go run .` |
| **Examen final** – Microservice URLWatch | `src/EXO_6/urlwatch/` | `cd src/EXO_6/urlwatch && go run ./cmd/urlwatch` |