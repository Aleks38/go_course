# Article 2-5-1 : Tableaux, slices et maps en Go - Manipulation, dimensionnement, allocation, pièges et optimisations

## 2-5-Fondamentaux du langage – Tableaux, slices et maps

### Introduction

Go offre trois types fondamentaux pour gérer des collections de données : les **tableaux (arrays)**, les **slices (tranches)** et les **maps**. Chacun a ses particularités en termes de manipulation, allocation mémoire et performances. Comprendre leurs différences, leur dimensionnement, ainsi que les pièges classiques et optimisations associées est indispensable pour écrire du code efficace et sûr.

---

## 1. Tableaux (Arrays)

- Taille fixe déterminée à la compilation.
- Syntaxe : `[taille]T` où T est le type élémentaire.
- Stockage contigu en mémoire.
- Peu utilisés directement car rigides (taille fixe).

**Exemple :**

```go
var arr [3]int = [3]int{1, 2, 3}
fmt.Println(arr)
```

- Copiés par valeur : l’affectation ou passage en argument crée une copie entière du tableau.

---

## 2. Slices (tranches)

- Vue flexible sur un tableau, dynamique en taille.
- Représentées par une structure interne avec un pointeur sur tableau, une longueur `len` et une capacité `cap`.
- Syntaxe : `[]T` (sans taille).
- Partagent la mémoire sous-jacente avec d'autres slices.

**Créer une slice :**

```go
s := []int{1, 2, 3}
```

**Création via `make` avec capacité :**

```go
s := make([]int, 0, 5) // longueur 0, capacité 5
```

**Append :**  

La fonction `append` permet d'ajouter des éléments. Si la capacité est dépassée, Go alloue un nouveau tableau plus grand (souvent double) et copie les données.

```go
s = append(s, 4, 5)
```

**Attention aux pièges :**

- Copier une slice copie la structure (pointeur, len, cap) mais pas les données, donc plusieurs slices peuvent partager la même mémoire.
- Modifier un élément via une slice modifie le tableau sous-jacent visible par d’autres slices pointant dessus.
- Redimensionner (append au-delà de la capacité) crée une nouvelle allocation : la slice est dissociée.

---

## 3. Maps

- Structure de données clé-valeur.
- Syntaxe : `map[KeyType]ValueType`.
- Allocation via `make` :

```go
m := make(map[string]int)
m["age"] = 30
```

- Opérations courantes : insert, accès, suppression (`delete`), test d’existence.

---

## 4. Dimensionnement et allocation

### Tableaux

- Allocation statique, taille fixe.
- Usage adapté aux données de taille connue et constante.

### Slices

- Allocation dynamique.
- Capacité et longueur contrôlent la croissance.
- Pour optimiser, utiliser `make` avec une capacité estimée pour limiter les réallocations fréquentes.

**Exemple :**

```go
s := make([]int, 0, 100) // capacité prévue à 100 éléments
```

### Maps

- La capacité initiale ne peut être spécifiée directement, mais elle est ajustée dynamiquement.
- Pour grandes maps, peut être utile d’estimer la taille pour limiter les réallocations internes.

---

## 5. Pièges courants

- Modifier une slice partagée sans vouloir impacter toutes les vues.
- Append sur une slice peut invalider les autres slices partageant le même tableau.
- Accès à une map avec clé absente retourne la valeur zéro du type, peut induire en erreur si non contrôlé. Utiliser la syntaxe avec deuxième valeur booléenne :

```go
v, ok := m["key"]
if ok {
    // clé présente
}
```

---

## 6. Optimisations

- Pré-allouer slices avec `make` pour éviter trop de copies à chaque append.
- Minimiser les copies de slices larges (travailler avec pointeurs quand nécessaire).
- Utiliser des types simples en clé de map pour optimiser performance.
- Eviter les reallocations en prévoyant la taille si possible.

---

## 7. Exemples combinés

```go
package main

import "fmt"

func main() {
    // Tableau
    var arr [3]int = [3]int{10, 20, 30}
    fmt.Println("Array:", arr)

    // Slice avec make
    s := make([]int, 0, 5)
    s = append(s, 1, 2, 3)
    fmt.Printf("Slice: len=%d cap=%d %v\n", len(s), cap(s), s)

    // Map
    m := make(map[string]int)
    m["alice"] = 25
    m["bob"] = 30
    fmt.Println("Map:", m)

    // Test existence clé
    if age, ok := m["alice"]; ok {
        fmt.Println("Alice a", age, "ans")
    }
}
```

---

## 8. Diagramme Mermaid : relations entre tableaux et slices

```mermaid
graph TD
    Array[Array (taille fixe)]
    Slice[Slice]
    Slice -->|pointeur| Array
    Slice -->|len| Length
    Slice -->|cap| Capacity

    Append -->|si cap dépassée| NewAlloc[Nouvelle allocation + copie]
    Append -->|sinon| ExtendLength[Extension longueur]
```

---

## 9. Sources

- [Go Blog - Slices: usage and internals](https://blog.golang.org/slices-intro)
- [Go by Example - Arrays](https://gobyexample.com/arrays)
- [Go by Example - Slices](https://gobyexample.com/slices)
- [Go by Example - Maps](https://gobyexample.com/maps)
- [Golang Docs - Effective Go: Slices](https://go.dev/doc/effective_go#slices)

---

Cet article clarifie les mécanismes de base pour manipuler tableaux, slices et maps, en soulignant les points critiques à surveiller et les meilleures pratiques pour optimiser vos programmes Go.