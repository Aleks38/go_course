# Article 1-4-1 : Packages, modules, workspace avec go.mod et commandes associées

## 1-4-Structure d'un projet Go

### Introduction

Go a introduit depuis la version 1.11 un système de gestion de dépendances basé sur les **modules** pour remplacer la gestion classique via GOPATH. Cette nouvelle organisation simplifie grandement la gestion des dépendances, la versionnage et la reproductibilité des builds. Cet article explore la notion de **package**, la structure des **modules**, le fichier `go.mod`, les workspaces, ainsi que les principales commandes pour gérer un projet Go moderne.

---

## 1. Packages : unité de compilation et d’organisation du code

- Un **package** représente un ensemble de fichiers `.go` regroupés sous un même nom.
- Il sert à organiser le code source en unités fonctionnelles réutilisables.
- Chaque fichier `.go` commence par la ligne :

```go
package nomDuPackage
```

- Le package `main` est un cas spécial qui produit un exécutable avec une fonction `main()` d’entrée.

### Exemple minimal (package utils)

```go
package utils

func Add(a, b int) int {
    return a + b
}
```

Un autre fichier dans un projet peut importer ce package `utils` :

```go
import "monprojet/utils"
```

---

## 2. Modules : gestion des dépendances via go.mod

### Qu'est-ce qu'un module ?

- Un **module** est un ensemble de packages versionné, défini par un répertoire racine contenant un fichier `go.mod`.
- Le fichier `go.mod` décrit :
  - Le nom (path) du module.
  - Les dépendances externes (autres modules) avec leurs versions.
  - Les contraintes de versionnage.
- Il permet à Go de résoudre les dépendances et versions à utiliser.

### Initialisation d’un module

Dans le répertoire racine du projet, lancez :

```bash
go mod init monmodule/nomprojet
```

Cela crée un fichier `go.mod` :

```go
module monmodule/nomprojet

go 1.20
```

---

## 3. Fichier `go.mod` - Exemple détaillé

```go
module github.com/user/monprojet

go 1.20

require (
    github.com/some/dependency v1.2.3
    golang.org/x/text v0.4.0
)
```

- Liste des modules requis et leurs versions exactes.
- Ce fichier est auto-généré et mis à jour par les commandes Go.

---

## 4. Workspaces (`go.work`) : collaboration multi-modules

Introduit pour gérer plusieurs modules localement dans un même workspace, par exemple lors du développement simultané de plusieurs projets liés.

### Création d’un workspace

```bash
go work init ./moduleA ./moduleB
```

Génère un fichier `go.work` qui lie plusieurs modules :

```go
go 1.20

use (
    ./moduleA
    ./moduleB
)
```

---

## 5. Commandes clés pour gérer modules et packages

| Commande              | Description                                         |
|----------------------|-----------------------------------------------------|
| `go mod init`         | Initialise un nouveau module Go                     |
| `go mod tidy`         | Nettoie les dépendances inutilisées ou manquantes. |
| `go mod download`     | Télécharge les dépendances décrites dans go.mod    |
| `go get [paquet@version]` | Ajoute ou met à jour une dépendance               |
| `go mod verify`       | Vérifie l’intégrité des modules téléchargés        |
| `go work init [dirs]` | Crée un workspace liant plusieurs modules          |

---

## 6. Exemple : organisation d’un projet simple

```plaintext
monprojet/
├── go.mod
├── main.go
├── utils/
│   └── utils.go
└── go.work  (optionnel, pour plusieurs modules)
```

- `main.go` contient `package main` et une fonction `main()`.
- `utils/` est un package importable comme `"monprojet/utils"`.

---

## 7. Exemple simple de go.mod et import

- `main.go`

```go
package main

import (
    "fmt"
    "monprojet/utils"
)

func main() {
    result := utils.Add(2, 3)
    fmt.Println("Resultat:", result)
}
```

- `utils/utils.go`

```go
package utils

func Add(a, b int) int {
    return a + b
}
```

---

## 8. Diagramme Mermaid : Organisation classique d’un module Go

```mermaid
graph TD
    A[Répertoire racine du module]
    A --> B[go.mod]
    A --> C[main.go (package main)]
    A --> D[utils (package utils)]
    D --> D1[utils.go]

    subgraph Modules externes
    E[github.com/some/lib v1.2.3]
    end

    B -->|dépend| E
```

---

## Conclusion

Le système de modules Go avec `go.mod` a simplifié la gestion des dépendances et amélioré la modularité des projets. Comprendre les packages, les modules et le workspace permet d’organiser efficacement son code, gérer les versions des dépendances et garantir à la fois portabilité et reproductibilité des builds.

---

## Sources

- [Go Modules Reference](https://go.dev/ref/mod)
- [Go Blog - Using Go Modules](https://blog.golang.org/using-go-modules)
- [Go Command Documentation](https://go.dev/doc/cmd/go)
- [Working with Go Modules - Official](https://go.dev/doc/modules/managing-dependencies)
- [Go Workspace Proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/45716-workspaces.md)