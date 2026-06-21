# JOURNAL_IA — URLWatch

## Méthode générale

L'IA m'a servi à **me challenger** et
à **revoir mon implémentation quand il le fallait** : remettre en question un
choix, proposer une alternative, pointer un point faible. J'ai avancé avec cette philosophie, **partie par partie** — je relisais et je validais chaque étape avant de passer à la suivante, et je la reprenais quand elle anticipait trop ou sortait du périmètre demandé.

Je me suis aussi appuyé sur l'IA pour **rédiger le `DESIGN.md`**. Il lui est
arrivé de **proposer des façons de faire qu'on n'avait pas vues tout au long du
cours** — je les ai alors évaluées avant de les retenir ou de les écarter.

## Ce que l'IA m'a fait découvrir

Deux éléments de la bibliothèque standard ne sont pas couverts par le cours, et c'est
l'IA qui me les a fait découvrir :

- **`http.NewRequestWithContext`** : m'a permis de faire de vraies requêtes HTTP
  sortantes qui **respectent le `context`**. C'est grâce à ça que le timeout par URL
  et l'annulation décidés par le pool **interrompent réellement** une requête en
  cours (un simple `http.Get` n'aurait pas permis ça).
- **`net/http/httptest`** : m'a permis de **tester mes handlers en isolation**, sans
  démarrer un vrai serveur — par exemple vérifier qu'un `POST /v1/checks` renvoie
  `201` et qu'un `GET /v1/checks/{id}` inconnu renvoie `404`.
