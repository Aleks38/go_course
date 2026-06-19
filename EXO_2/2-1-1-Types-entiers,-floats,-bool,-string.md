# Article 2-1-1 : Types entiers, floats, bool, string, inférence automatique en Go

## 2-1-Fondamentaux du langage – Types de base, variables, constantes

### Introduction

Go propose plusieurs **types de base** pour représenter des valeurs numériques, booléennes, et textuelles. Par ailleurs, le langage dispose d’une **inférence automatique** de type efficace qui simplifie la déclaration des variables tout en garantissant la sécurité statique du typage.

---

## 1. Les types de base en Go

### 1.1 Types entiers

Go distingue plusieurs variantes d’entiers signés et non signés avec des tailles précises :

| Type     | Taille          | Description                     |
|----------|-----------------|--------------------------------|
| `int`    | 32 ou 64 bits*  | Taille dépend de l’architecture|
| `int8`   | 8 bits          | Entier signé                   |
| `int16`  | 16 bits         | Entier signé                   |
| `int32`  | 32 bits         | Entier signé                   |
| `int64`  | 64 bits         | Entier signé                   |
| `uint`   | 32 ou 64 bits*  | Entier non signé               |
| `uint8`  | 8 bits          | Alias de byte                 |
| `uint16` | 16 bits         | Entier non signé              |
| `uint32` | 32 bits         | Entier non signé              |
| `uint64` | 64 bits         | Entier non signé              |

\* Sur les architectures 64 bits, `int` et `uint` sont 64 bits, sinon 32 bits.

**Exemple :**

```go
var a int = 42
var b uint8 = 255
```

### 1.2 Types flottants (floats)

Pour représenter les nombres à virgule flottante, Go utilise deux types :

- `float32` : précision simple (32 bits)
- `float64` : précision double (64 bits), recommandé pour la plupart des usages.

Exemple :

```go
var f1 float32 = 3.14
var f2 float64 = 2.718281828
```

### 1.3 Type booléen

Il ne peut prendre que deux valeurs : `true` ou `false`.

```go
var flag bool = true
```

### 1.4 Type string

Une chaîne de caractères UTF-8 immuable.

```go
var s string = "Bonjour, Go!"
```

---

## 2. Inférence automatique de type

Go propose une syntaxe de déclaration simplifiée via l’opérateur `:=` qui infère automatiquement le type depuis la valeur à droite.

**Exemple :**

```go
a := 42           // int
b := 3.14         // float64
flag := true      // bool
s := "Hello"      // string
```

Cette déclaration ne peut être utilisée qu’à l’intérieur d’une fonction.

---

## 3. Constantes

Les constantes sont déclarées avec `const`, elles doivent être évaluées à la compilation et ne peuvent être modifiées.

```go
const Pi = 3.14159
const Greeting = "Bonjour"
const IsActive = true
```

---

## 4. Exemple complet

```go
package main

import "fmt"

func main() {
    var i int = 10
    j := 20.5               // inférence float64
    var isReady bool = true
    msg := "Bonjour Go"

    fmt.Println("i =", i)
    fmt.Println("j =", j)
    fmt.Println("isReady =", isReady)
    fmt.Println("msg =", msg)
}
```

---

## 5. Diagramme Mermaid : catégorisation des types de base Go

```mermaid
graph TD
    A[Types de base Go]
    A --> B[Entiers]
    A --> C[Float]
    A --> D[Bool]
    A --> E[String]

    B --> B1[int (32/64 bits)]
    B --> B2[int8, int16, int32, int64]
    B --> B3[uint, uint8 (byte), uint16, uint32, uint64]

    C --> C1[float32]
    C --> C2[float64]

    D --> D1[true]
    D --> D2[false]
```

---

## 6. Références

- [Go by Example - Variables](https://gobyexample.com/variables)
- [Go Language Spec - Basic Types](https://golang.org/ref/spec#Types)
- [Effective Go - Variables](https://go.dev/doc/effective_go#variables)
- [Tour of Go - Values](https://tour.golang.org/basics/1)
- [Go Wiki - FAQ Type System](https://github.com/golang/go/wiki/FAQ#what-kind-of-typing-does-go-use)

---

Cet article vous donne une vision claire des types primitifs classiques en Go ainsi que des mécanismes d’inférence typée et de constantes. Ces bases sont indispensables pour écrire un code sûr, clair et performant.