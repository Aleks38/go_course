# Examen final — Langage Go

**Durée : 2 h 00 à 2 h 30 — Épreuve sur table, sur votre PC**
**Modalité de rendu : un dépôt Git individuel**

---

## 1. Contexte et règles de l'épreuve

Cet examen prend la forme d'un **mini-projet pratique unique**. Vous développez, de bout en bout, un petit microservice Go nommé **URLWatch**. Le sujet est marqué sur la **concurrence** et le **développement backend**.

### Ressources autorisées

-  Vos notes de cours et la documentation officielle Go (`pkg.go.dev`, `go.dev`).
-  Internet.
-  L'assistance d'une IA générative — **à condition d'un usage réfléchi** (voir ci-dessous).

### Usage réfléchi de l'IA — lisez attentivement

L'IA est un **outil**, pas un sous-traitant. Vous restez l'architecte : vous devez comprendre, justifier et savoir défendre **chaque ligne** que vous rendez. En conséquence :

- Vous tenez un fichier **`JOURNAL_IA.md`** où vous notez brièvement *comment* vous avez utilisé l'IA, ce que vous avez **accepté, rejeté ou modifié**, et **pourquoi**.
- Plusieurs livrables (contrat d'API précis, `DESIGN.md` lié à *votre* code, tests qui doivent réellement passer, historique de commits) sont conçus pour qu'un copier-coller non réfléchi soit **insuffisant** et facilement détectable.

### Intégrité

Aucun secret (clé, token) ne doit être commité.

---

## 2. Le projet : URLWatch

**URLWatch** est un microservice de **vérification d'URLs en masse**. Un client lui envoie une liste d'URLs ; le service les interroge **en parallèle** (code HTTP, latence, disponibilité), **agrège** les résultats et les **expose** via une API REST. Chaque lot de vérifications (« batch ») est conservé et consultable a posteriori.

### Scénario fonctionnel

1. Un client envoie un lot d'URLs à vérifier, avec des options (niveau de parallélisme, délai d'expiration).
2. Le service vérifie les URLs **concurremment**, en bornant le nombre de requêtes simultanées, et en respectant un **timeout** et une **annulation** propres via `context`.
3. Le service calcule un résumé (combien sont accessibles / en échec, durée totale) et **persiste** le lot.
4. Le client peut **relire** un lot par son identifiant.

---

## 3. Prérequis techniques et structure du dépôt

- **Go 1.22 ou supérieur** (le paquet `log/slog` est requis, disponible depuis Go 1.21).
- Projet organisé en **module** (`go.mod`) avec un découpage en **packages** cohérent.
- Votre dépôt doit se construire et se tester avec des commandes standard :

```bash
go build ./...
go vet ./...
go test ./...
go run ./cmd/urlwatch
```

### Arborescence suggérée (vous pouvez l'adapter, en le justifiant)

```text
urlwatch/
├── go.mod
├── README.md              # build / run / test + exemples curl
├── DESIGN.md              # justification architecturale (voir §9)
├── JOURNAL_IA.md          # usage de l'IA (voir §1)
├── .gitignore
├── cmd/
│   └── urlwatch/
│       └── main.go        # point d'entrée, câblage des dépendances, démarrage du serveur
└── internal/
    ├── domain/            # types métier, erreurs, interfaces (Checker, Store)
    ├── checker/           # implémentation HTTP du Checker (+ mock pour les tests)
    ├── pool/              # cœur concurrent : worker pool, fan-out / fan-in
    ├── store/             # implémentation en mémoire (+ SQLite en bonus)
    └── api/               # handlers HTTP, routage, middleware, DTO JSON
```

> **Conseil :** placez vos types et interfaces dans `domain` pour que les autres packages en dépendent (inversion de dépendance), et gardez `cmd/urlwatch/main.go` mince (assemblage uniquement).

---

## 4. Partie 0 — Mise en place

1. Initialisez le module Go et l'arborescence de packages.
2. Initialisez le dépôt Git, ajoutez un `.gitignore` adapté à Go.
3. **Premier commit** (« init module + structure »).

> **Attendu transversal :** des **commits incrémentaux** et atomiques tout au long de l'épreuve, avec des messages clairs. Un dépôt avec un seul méga-commit final sera pénalisé.

---

## 5. Partie 1 — Modélisation du domaine

### 5.1 Types métier

Définissez au minimum les structures suivantes (champs à compléter, **avec les `struct tags` JSON**) :

- `CheckResult` — résultat de la vérification d'une URL : l'URL, le code HTTP, un booléen de disponibilité, la latence en millisecondes, et un éventuel message d'erreur.
- `Batch` — un lot : identifiant, date de création, la liste des `CheckResult`, et un **résumé** agrégé (total, nombre d'accessibles, nombre en échec, durée totale).

L'agrégation du résumé doit s'appuyer sur des **slices** et/ou des **maps** de manière idiomatique.

### 5.2 Interfaces (contrat minimal)

Le service doit reposer sur des **interfaces** afin d'être testable et évolutif. Contrat minimal suggéré (vous pouvez enrichir/ajuster en le justifiant dans `DESIGN.md`) :

```go
// Vérifie une URL unique. L'implémentation par défaut fait un vrai appel HTTP ;
// une implémentation mock (déterministe) sera utilisée dans les tests.
type Checker interface {
    Check(ctx context.Context, url string) CheckResult
}

// Persiste et relit les lots.
type Store interface {
    Save(ctx context.Context, b Batch) error
    Get(ctx context.Context, id string) (Batch, error)
}
```

### 5.3 Gestion des erreurs idiomatique

- Définissez une **erreur sentinelle** `ErrBatchNotFound` retournée par `Store.Get` quand l'identifiant est inconnu.
- Définissez au moins **une erreur personnalisée** (type implémentant `error`) — par exemple une erreur de validation portant le nom du champ fautif.
- Utilisez le **wrapping** (`fmt.Errorf("...: %w", err)`) et exploitez `errors.Is` / `errors.As` là où c'est pertinent (notamment pour traduire une erreur métier en code HTTP dans la couche API).

---

## 6. Partie 2 — Cœur concurrent — **partie centrale**

Implémentez le moteur qui vérifie un lot d'URLs **en parallèle**.

### Exigences fonctionnelles

1. **Worker pool borné** : le nombre de requêtes HTTP **simultanées** ne doit **jamais dépasser** la valeur `concurrency` demandée. (Une goroutine par URL *sans limite* est **refusée**.)
2. **Fan-out / fan-in** via des **channels** : distribution des URLs aux workers, collecte des résultats. Choisissez channels **bufferisés ou non**, **directionnels** — et **justifiez** ce choix dans `DESIGN.md`.
3. **`context`** :
   - Un **timeout global** pour le traitement d'un lot **et** un **timeout par URL**.
   - À l'expiration du délai ou en cas d'annulation, les vérifications en cours sont **interrompues proprement**.
4. **Synchronisation** : usage correct de `sync.WaitGroup` pour attendre la fin des workers ; `sync.Mutex` ou `sync.Once` si votre design le justifie.
5. **Aucune fuite de goroutine** et **aucun blocage** : tous les channels sont correctement fermés, toutes les goroutines se terminent.

### Anti-patterns refusés

- Lancer une goroutine par URL sans borne.
- Ignorer `ctx.Done()`.
- Channel jamais fermé provoquant un *deadlock*, ou `WaitGroup` mal géré.
- *Data race* (votre code doit passer `go test -race`).

---

## 7. Partie 3 — API REST, JSON et logging

Vous pouvez utiliser **`net/http`** (stdlib) **ou** **Gin**. Ce choix est une **décision d'architecture** : justifiez-le dans `DESIGN.md`.

### 7.1 Endpoints

| Méthode | Chemin | Rôle |
|---|---|---|
| `POST` | `/v1/checks` | Crée et exécute un lot, le persiste, renvoie le résultat. |
| `GET` | `/v1/checks/{id}` | Renvoie un lot existant (ou `404`). |
| `GET` | `/healthz` | Sonde de vivacité. |

### 7.2 Contrat JSON (à respecter précisément)

**Requête** `POST /v1/checks` :

```json
{
  "urls": ["https://go.dev", "https://exemple.invalid"],
  "options": { "concurrency": 4, "timeout_ms": 2000 }
}
```

Règles de validation :
- `urls` : obligatoire, **1 à 100** entrées, chacune une URL `http`/`https` valide.
- `options.concurrency` : optionnel, défaut `8`, borné `1..50`.
- `options.timeout_ms` : optionnel (timeout par URL), défaut `5000`, borné `100..30000`.

**Réponse** `201 Created` :

```json
{
  "batch_id": "b_4f3c1a",
  "created_at": "2026-06-18T10:00:00Z",
  "summary": { "total": 2, "up": 1, "down": 1, "duration_ms": 812 },
  "results": [
    { "url": "https://go.dev", "status_code": 200, "ok": true, "latency_ms": 120 },
    { "url": "https://exemple.invalid", "ok": false, "error": "dns: no such host", "latency_ms": 2001 }
  ]
}
```

**Contrat d'erreur** — toute erreur renvoie ce corps, avec le bon code HTTP :

```json
{ "error": { "code": "batch_not_found", "message": "aucun lot avec l'id b_x" } }
```

- `400 invalid_request` (corps invalide / validation), `404 batch_not_found`, `405`, `500 internal`…
- Les `struct tags` doivent produire **exactement** ces noms de champs (`snake_case`), avec `omitempty` là où c'est pertinent (ex. `error`, `status_code` absents si non significatifs).

### 7.3 Logging structuré (`log/slog`)

- Configurez un logger **`slog`** avec un **handler JSON**.
- Un **middleware** journalise chaque requête avec au minimum : `method`, `path`, `status`, `duration_ms`, et `batch_id` lorsqu'il est connu.
- Le niveau de log est **configurable** (ex. variable d'environnement `LOG_LEVEL`).
- `/healthz` ne doit pas polluer les logs applicatifs (à vous de gérer).

> Bonus apprécié : un middleware de **recovery** qui transforme une `panic` en réponse `500` propre (et la journalise).

---

## 8. Partie 4 — Tests

`go test ./...` doit **passer**. Écrivez des tests **pertinents**, pas de remplissage :

1. **Tests table-driven** sur la logique métier : validation des options, agrégation du résumé, et le **worker pool** alimenté par un **`Checker` mock déterministe** (sans réseau).
2. **`net/http/httptest`** pour tester les handlers : au minimum un `POST /v1/checks` réussi **et** un `GET /v1/checks/{id}` renvoyant `404`.
3. **Mocking par interfaces** : `Checker` et `Store` sont remplacés par des mocks dans les tests (aucun appel réseau ni base réelle requis).
4. **Concurrence** : au moins un test vérifie l'**annulation / le timeout** par `context`. Votre suite doit être propre sous **`go test -race`**.

---

## 9. Livrables de réflexion — **fortement valorisés**

### `DESIGN.md` (1 page maximum)

Justifiez vos décisions, **en référence à votre propre code** :

1. **Découpage** : pourquoi ce partage en packages ? Où avez-vous placé les frontières d'**interface**, et pourquoi ?
2. **Modèle de concurrence** : pourquoi cette taille de pool / cette bufferisation des channels ? Comment gérez-vous les **échecs partiels** d'un lot ?
3. **Fuites de goroutines** : quels risques dans votre design, et comment les évitez-vous concrètement ?
4. **Erreurs** : votre stratégie (sentinelles, types, wrapping) et son lien avec les codes HTTP.
5. **Philosophie Go** : citez **3 arguments**, liés à *votre* implémentation, pour lesquels Go est un bon choix ici plutôt que Java, Python ou Rust — et **une limite** que vous avez ressentie.

### `JOURNAL_IA.md`

Quelques lignes honnêtes : où l'IA vous a aidé, ce que vous avez **rejeté ou corrigé**, et pourquoi.

---

## 10. Partie 5 — Extensions bonus *(si le temps le permet)*

À traiter **seulement après** un socle fonctionnel et testé. Chaque extension reste **derrière vos interfaces** existantes.

- **Persistance SQLite** via `database/sql` (ou `sqlx` / `GORM`) : une implémentation `Store` sur SQLite, avec création de schéma. Le choix `database/sql` vs ORM doit être justifié.
- **Arrêt gracieux** : `http.Server` + capture de `SIGINT`/`SIGTERM` (`os/signal`) + `context` ; les requêtes et goroutines en cours **se terminent proprement** avant l'arrêt.
- **Asynchrone** : `POST` renvoie `202 Accepted` + `batch_id` immédiatement, le lot se calcule en tâche de fond, `GET` reflète l'avancement (`pending` / `done`).
- **Pagination / filtre** sur une route de liste `GET /v1/checks`.
- **Dockerfile** multi-stage produisant une image minimale.

---

## 11. Ce que vous devez rendre (checklist)

- [ ] Dépôt **Git individuel** avec **historique de commits incrémental**.
- [ ] Code qui compile : `go build ./...` sans erreur.
- [ ] `go test ./...` (idéalement `-race`) **vert**.
- [ ] API conforme au **contrat JSON**.
- [ ] `README.md` (build / run / test + exemples `curl`), `DESIGN.md`, `JOURNAL_IA.md`.
- [ ] Aucun secret commité ; `.gitignore` présent.

---

## 12. Critères de qualité

Aucune répartition de points ni de durée n'est imposée : **organisez vous-mêmes votre temps et vos priorités**. La qualité est appréciée sur :

- **Absence de *data race*** (`go test -race`) et de fuite de goroutine.
- **Lisibilité et architecture** : packages cohérents, interfaces **minimales**, `main` mince.
- **Pertinence des tests** : couverture des **chemins d'erreur**, pas seulement du cas nominal.

*Bon courage. Privilégiez un service **simple, correct et bien testé** à un service ambitieux mais cassé.*
