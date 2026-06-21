# DESIGN — URLWatch

## 1. Découpage en packages

`domain` contient les types métier (`CheckResult`, `Summary`, `Batch`), les erreurs
et les **interfaces** (`Checker`, `Store`). Les
frontières d'interface sont placées dans `domain` : `pool.Run` dépend de
`domain.Checker` (pas de l'implémentation HTTP) et `api.Server` de `domain.Checker` /
`domain.Store`. Conséquence : grâce à la **satisfaction implicite des interfaces**
(*duck typing*), les tests injectent des mocks, sans réseau ni base. `checker` et
`store` sont les implémentations concrètes. J'ai suivi l'arborescence
suggérée par le sujet, qui correspond bien à ce découpage.

**`net/http` plutôt que Gin.** J'ai choisi la stdlib pour deux raisons : éviter
**toute dépendance externe** (`go.mod` sans `require`), et parce que le middleware de
logging exigé (§7.3) reprend exactement le pattern d'enveloppe
`func(next) http.Handler` du cours — alors que `gin.Default()` embarque son propre
Logger/Recovery qui entrerait en conflit avec mon middleware `slog` JSON. Le routage
par méthode et l'enveloppe d'erreur restent ainsi sous mon contrôle direct.

## 2. Modèle de concurrence

`pool.Run` lance **exactement `concurrency` goroutines** workers (borne issue de la
requête, validée `1..50`) — jamais une goroutine par URL. La **taille du pool vient
donc du client**, qui arbitre lui-même débit contre charge. Le fan-out passe par le
channel `jobs`, le fan-in par `results`, tous deux **non bufferisés** et
**directionnels** côté worker.

Channels **non bufferisés** : chaque envoi attend un récepteur prêt (**synchronisation
stricte**), mémoire bornée (pas de buffer proportionnel au nombre d'URLs) ; un buffer
n'apporterait rien puisque le collecteur consomme `results` en continu. La
**directionnalité** (`<-chan` / `chan<-`) documente l'intention et empêche un mauvais
usage à la compilation.

**Échecs partiels** : un échec d'URL (timeout, DNS, code ≥ 400) n'est pas une erreur
de programme mais un **résultat normal** (`CheckResult.OK=false`) — d'où le fait que
`Checker.Check` ne renvoie pas d'`error` (pas de canal d'erreurs séparé). Le lot
aboutit toujours, `domain.Summarize` agrège `up`/`down`, et l'ordre d'entrée est
préservé via `indexedResult`.

## 3. Fuites de goroutines

Risques : un channel jamais fermé, un worker bloqué en envoi. Évités ainsi : le
producteur ferme `jobs` ; une goroutine dédiée fait `wg.Wait()` puis `close(results)` ;
le collecteur draine `results` jusqu'à fermeture → toutes les goroutines se terminent.
La garde `concurrency >= 1` empêche un producteur bloqué faute de worker. En cas
d'annulation ou de timeout, `check` court-circuite sur `ctx.Done()` (pré-test `select`
+ `context.WithTimeout` par URL). C'est **prouvé** par `go test -race ./...` et les
tests d'annulation et de timeout du pool.

## 4. Stratégie d'erreurs

Deux mécanismes complémentaires, traduits en HTTP dans `api` :

- **Sentinelle** `domain.ErrBatchNotFound`, renvoyée par `MemoryStore.Get`, **enroulée**
  avec `%w` (*wrap*) → le handler la détecte avec `errors.Is` → `404 batch_not_found`.
- **Type personnalisé** `domain.ValidationError{Field, Message}`, renvoyé par
  `validate` → détecté avec `errors.As` → `400 invalid_request`, avec le champ fautif.

Toute erreur sort sous la même enveloppe JSON `{ "error": { "code", "message" } }`.

## 5. Philosophie Go

Trois raisons, liées à cette implémentation, pour lesquelles Go est ici un meilleur
choix que Java, Python ou Rust :

1. **Concurrence native.** Le worker pool, le fan-out/fan-in et l'annulation
   s'écrivent directement avec goroutines, channels et `context` (`pool.Run`), sans
   bibliothèque tierce — là où Java demande un `ExecutorService` plus lourd, Python
   bute sur le GIL pour du vrai parallélisme, et Rust offrirait ce parallélisme mais
   au prix d'une gestion mémoire (ownership, lifetimes) surdimensionnée pour ce service.
2. **Bibliothèque standard suffisante.** `net/http`, `log/slog`, `encoding/json` et
   `testing` couvrent tout le projet : le `go.mod` n'a **aucune dépendance externe**,
   là où l'équivalent Java (Spring) ou Python (FastAPI + libs) tirerait un arbre de
   dépendances bien plus large.
3. **Interfaces implicites + outillage.** De petites interfaces (`Checker`, `Store`)
   rendent le code testable par mocks immédiats, et `go test -race` détecte les data
   races sans configuration — simplicité sans la cérémonie de Java, sûreté concurrente
   sans la complexité du borrow checker de Rust.

**Une limite ressentie.** Go privilégie l'explicite au concis, ce qui alourdit
parfois le code : la gestion d'erreurs répétitive (`if err != nil` partout dans le
checker et les handlers), et des contournements pour des besoins simples — par
exemple, pour lire le code HTTP de réponse dans le middleware de logging, j'ai dû
**envelopper `http.ResponseWriter`** dans `statusRecorder`, faute d'accès direct au
statut.
