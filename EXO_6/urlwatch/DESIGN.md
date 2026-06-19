# DESIGN — URLWatch

> Document de justification architecturale (à compléter, 1 page max). Chaque
> section doit référencer *mon propre code*.

## 1. Découpage en packages

Pourquoi ce partage en packages ? Où sont placées les frontières d'interface, et
pourquoi ? *(à compléter)*

## 2. Modèle de concurrence

Taille du pool, bufferisation et directionnalité des channels, gestion des échecs
partiels d'un lot. *(à compléter)*

## 3. Fuites de goroutines

Risques dans ce design et comment ils sont évités concrètement. *(à compléter)*

## 4. Stratégie d'erreurs

Sentinelles, types personnalisés, wrapping (`%w`), et lien avec les codes HTTP
(`errors.Is` / `errors.As` dans la couche API). *(à compléter)*

## 5. Philosophie Go

Trois arguments (liés à mon implémentation) pour lesquels Go est un bon choix ici
plutôt que Java, Python ou Rust — et une limite ressentie. *(à compléter)*
