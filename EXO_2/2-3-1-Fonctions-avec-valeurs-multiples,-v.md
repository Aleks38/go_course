# Article 2-3-1 : Fonctions avec valeurs multiples, variadiques, fonctions comme variables en Go

## 2-3-Fondamentaux du langage – Fonctions

### Introduction

Les fonctions en Go offrent des fonctionnalités avancées tout en restant simples à utiliser. Parmi les caractéristiques puissantes figurent la capacité à retourner plusieurs valeurs, à accepter un nombre variable d’arguments (variadique) et la possibilité de manipuler les fonctions comme des valeurs (variables, arguments, retours). Cet article aborde ces concepts clés avec des exemples concrets.

---

## 1. Fonctions à valeurs multiples

Go permet à une fonction de retourner plusieurs résultats. C’est très utile pour retourner une valeur et une erreur, ou plusieurs résultats complémentaires.

**Exemple : fonction division retournant quotient et reste**

```go
func divMod(a, b int) (int, int) {
    return a / b, a % b
}
```

**Utilisation :**

```go
q, r := divMod(10, 3)
fmt.Println("Quotient:", q, "Reste:", r)
```

Cette capacité évite souvent le recours à des structures ou à l’usage d’exceptions.

---

## 2. Fonctions variadiques

Une fonction variadique accepte un nombre variable d’arguments du même type. Le dernier paramètre est précédé de `...`.

**Exemple : somme de plusieurs entiers**

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}
```

**Appel possible avec n arguments :**

```go
fmt.Println(sum(1, 2, 3))
fmt.Println(sum(10, 20))
```

On peut aussi passer explicitement une slice :

```go
vals := []int{5,6,7}
fmt.Println(sum(vals...))
```

---

## 3. Fonctions comme variables (fonctions de première classe)

En Go, les fonctions sont des valeurs. Elles peuvent être affectées à des variables, passées en paramètre, ou retournées depuis d’autres fonctions.

**Exemple :**

```go
func add(a, b int) int {
    return a + b
}

func main() {
    var f func(int, int) int
    f = add
    fmt.Println(f(3, 4))
}
```

**Passage de fonction en argument :**

```go
func apply(op func(int,int)int, x, y int) int {
    return op(x, y)
}

result := apply(add, 5, 6)
fmt.Println(result) // 11
```

---

## 4. Exemple complet mêlant les notions

```go
package main

import "fmt"

func divMod(a, b int) (int, int) {
    return a / b, a % b
}

func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

func apply(op func(int, int) int, x, y int) int {
    return op(x, y)
}

func add(a, b int) int {
    return a + b
}

func main() {
    q, r := divMod(20, 7)
    fmt.Println("divMod:", q, r)

    fmt.Println("sum variadique:", sum(1, 2, 3, 4))

    f := add
    fmt.Println("fonction comme variable:", f(10, 15))

    result := apply(add, 5, 6)
    fmt.Println("fonction en argument:", result)
}
```

---

## 5. Diagramme Mermaid : relations fonctionnelles

```mermaid
graph TD
    A[main] --> B[divMod(a,b) returns (quotient, reste)]
    A --> C[sum(nums...)]
    A --> D[f = add]
    A --> E[apply(op, x, y)]

    E --> F["op : func(int,int)int"]
    F --> G[add(a,b)]
```

---

## 6. Références

- [Go by Example - Multiple Return Values](https://gobyexample.com/multiple-return-values)
- [Go by Example - Variadic Functions](https://gobyexample.com/variadic-functions)
- [Go by Example - Functions as Values](https://gobyexample.com/functions)
- [Go Language Specification - Function types](https://golang.org/ref/spec#Function_types)
- [Effective Go - Functions](https://go.dev/doc/effective_go#functions)

---

La gestion des fonctions multiples, variadiques et comme valeurs permet à Go d’être à la fois simple et puissant, favorisant un style de programmation clair et flexible.