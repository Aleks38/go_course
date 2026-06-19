# Article 2-4-1 : Usage des pointeurs en Go, passage par référence, limites et différences avec C

## 2-4-Fondamentaux du langage – Pointeurs en Go

### Introduction

Le concept de pointeur est fondamental en langage système comme Go. Cependant, Go adopte une approche simplifiée et sécurisée comparé à C, limitant certains usages tout en facilitant la gestion mémoire et évitant des erreurs classiques. Cet article détaille l’usage des pointeurs en Go, le passage par référence grâce aux pointeurs et les différences clés avec le langage C.

---

## 1. Définition et usage des pointeurs en Go

Un pointeur est une variable qui contient l’adresse mémoire d’une autre variable.

- Le type pointeur s’écrit `*T` où `T` est un type.
- L’opérateur `&` récupère l’adresse d’une variable.
- L’opérateur `*` permet d’accéder à la valeur pointée (déréférencement).

**Exemple :**

```go
package main

import "fmt"

func main() {
    x := 42
    p := &x           // p est un pointeur vers x
    fmt.Println(*p)   // déférencement : affiche 42

    *p = 21           // modifie la valeur pointée (x)
    fmt.Println(x)    // affiche 21
}
```

---

## 2. Passage par référence avec pointeurs

Par défaut, Go passe les arguments aux fonctions par copie (valeur). Pour modifier la variable d’origine, on utilise un pointeur.

**Exemple :**

```go
func increment(n *int) {
    *n = *n + 1
}

func main() {
    x := 5
    increment(&x)
    fmt.Println(x)    // affiche 6
}
```

---

## 3. Limites des pointeurs en Go comparé à C

| Aspect                 | Go                                      | C                                |
|------------------------|-----------------------------------------|---------------------------------|
| Arithmétique de pointeur | Non supportée — pas d’opérations sur les adresses | Supportée                       |
| Pointeurs sur pointeurs | Supporté                               | Supporté                       |
| Pointeurs `void*`      | Non existant, Go est typé strict        | Existe, pointeur générique      |
| Manipulation mémoire brute | Restreinte, pas d’accès direct à la mémoire arbitraire | Libre, mais risque de corruption |
| Garbage Collector       | Oui — gestion automatique               | Non — gestion manuelle          |

---

## 4. Conseils d’usage

- Utilisez les pointeurs pour éviter les copies coûteuses (ex : gros structs).
- Utilisez-les pour modifier les valeurs dans les fonctions.
- Évitez l’arithmétique de pointeurs en Go, privilégiez les structures et slices.
- La gestion mémoire est sécurisée, pas de libres ou malloc explicite.

---

## 5. Exemple complet : échange de valeurs à l’aide de pointeurs

```go
func swap(a, b *int) {
    temp := *a
    *a = *b
    *b = temp
}

func main() {
    x, y := 3, 7
    swap(&x, &y)
    fmt.Println(x, y)  // affiche 7 3
}
```

---

## 6. Diagramme Mermaid : fonctionnement d’un pointeur et passage par référence

```mermaid
graph LR
    subgraph Fonction main
        x[Variable x = 10]
        p[Pointeur p]
    end

    x -->|Adresse de x (&x)| p

    subgraph Fonction increment
        np[Pointeur n]
    end

    p -- Passage de l'adresse --> np
    np -- Valeur pointée --> x
    np --> |*n = *n+1| xUpdated(x = 11)
```

---

## 7. Résumé des différences majeures avec C

- **Pas d’arithmétique de pointeurs** : Go simplifie la sécurité mémoire.
- **Gestion automatique** via Garbage Collector : évite fuites liées aux pointeurs.
- **Typage plus strict** : pas de pointeurs génériques `void*`.
- **Syntaxe claire et sécurisée** utilisant `&` et `*` uniquement.

---

## 8. Sources

- [Go Blog - Pointers](https://go.dev/blog/pointers)
- [Go Language Spec - Pointers](https://golang.org/ref/spec#Pointer_types)
- [Go by Example - Pointers](https://gobyexample.com/pointers)
- [Effective Go - Pointers](https://go.dev/doc/effective_go#pointers)
- [C vs Go Pointers Discussion](https://stackoverflow.com/questions/38830480/difference-between-pointers-in-c-and-go)

---

Go propose ainsi une approche moderne et simplifiée des pointeurs, limitant les erreurs tout en offrant la puissance nécessaire pour la programmation système et la performance. Les pointeurs restent un outil indispensable quand on veut passer par référence ou éviter des copies inutiles.