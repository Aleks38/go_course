# Article 1-3-2 : Premier programme 'Hello World' en Go

## 1-3-Installation et premiers programmes

### Introduction

Le traditionnel programme **Hello World** est le point de départ pour découvrir un nouveau langage de programmation. En Go, ce programme est à la fois simple, lisible et illustre clairement des notions essentielles comme la structure d’un package, l’importation de bibliothèques et la fonction d’entrée `main`.

---

## 1. Structure minimale d’un programme Go

Un programme Go démarre avec une déclaration de package, généralement `package main` pour les exécutables :

```go
package main
```

Puis viennent les imports, c’est-à-dire les packages standard ou tiers utilisés, par exemple `fmt` pour l’affichage :

```go
import "fmt"
```

Enfin, la fonction `main` sert d’entrée au programme :

```go
func main() {
    fmt.Println("Hello, World!")
}
```

---

## 2. Code complet du premier programme Hello World

Voici le code complet du fichier `hello.go` :

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

- `package main` indique que ce fichier est une application exécutable.
- `import "fmt"` permet d’utiliser la bibliothèque standard `fmt` (formatage).
- `func main()` est la fonction obligatoire d’exécution.
- `fmt.Println` affiche le texte suivi d’un retour à la ligne.

---

## 3. Compilation et exécution

### Compilation

Dans le terminal, placez-vous dans le dossier contenant `hello.go` puis exécutez :

```bash
go build hello.go
```

Cette commande génère un exécutable nommé `hello` (ou `hello.exe` sur Windows) que vous pouvez lancer :

```bash
./hello
```

Vous verrez alors la sortie :

```
Hello, World!
```

### Exécution directe

Alternativement, vous pouvez exécuter sans compilation explicite :

```bash
go run hello.go
```

Cela compile et lance le programme en une seule étape.

---

## 4. Organisation du programme et fonctionnement interne

```mermaid
flowchart TD
    A[Fichier hello.go] --> B[package main]
    A --> C[import "fmt"]
    A --> D[func main()]
    D --> E[fmt.Println]
    E --> F[Affichage "Hello, World!"]
```

- Le compilateur Go identifie la fonction `main` comme point d’entrée.
- Le package `fmt` fournit la fonction `Println` pour afficher du texte.
- L’exécutable est linké et généré grâce à `go build`.

---

## 5. Remarques supplémentaires

- Chaque programme en Go qui produit un exécutable doit définir un package `main`.
- La fonction `main` ne prend pas d’arguments ni ne renvoie de valeur.
- `fmt.Println` ajoute automatiquement un saut de ligne après le texte.
- Le système de paquets de Go invite à structurer proprement de plus gros projets.

---

## 6. Exemple d’extension simple

Pour afficher plusieurs messages en Go, vous pouvez écrire :

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    fmt.Println("Bienvenue dans le langage Go.")
}
```

Exécution :

```
Hello, World!
Bienvenue dans le langage Go.
```

---

## Sources

- [Documentation officielle Go - Hello World](https://go.dev/doc/tutorial/hello)
- [Tour of Go - Basics](https://go.dev/tour/basics/1)
- [Go by Example - Hello World](https://gobyexample.com/hello-world)
- [Go Command Documentation](https://go.dev/doc/cmd/go)

---

Ce premier programme Go vous permet d’appréhender les bases syntaxiques et l’environnement de développement. Il pose aussi les fondations pour écrire rapidement des applications plus complexes.