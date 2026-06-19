package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Produit struct {
	ID        int
	Nom       string
	Prix      float64
	Categorie string
}

func main() {
	// PARTIE 1
	categories := []string{"Électronique", "Vêtements", "Livres"}
	categories = append(categories, "Chaussures", "Accessoires")

	for _, category := range categories {
		fmt.Println(category)
	}

	categorieExiste("Électronique", categories)
	categorieExiste("Test", categories)

	categories = supprimerCategorie("Électronique", categories)
	for _, category := range categories {
		fmt.Println(category)
	}

	categories = supprimerCategorie("Test", categories)
	for _, category := range categories {
		fmt.Println(category)
	}

	fmt.Println(len(categories))
	fmt.Println(cap(categories))

	// PARTIE 2: GESTION DES PRODUITS ET DU STOCK
	fmt.Println("\n=== PARTIE 2: GESTION DES PRODUITS ===")

	inventaireProduits := make(map[int]Produit)
	stockProduits := make(map[int]int)

	// Ajout de 3 produits
	inventaireProduits[1] = Produit{1, "Laptop", 899.99, "Électronique"}
	stockProduits[1] = 10

	inventaireProduits[2] = Produit{2, "T-Shirt", 29.99, "Vêtements"}
	stockProduits[2] = 50

	inventaireProduits[3] = Produit{3, "Harry Potter", 15.99, "Livres"}
	stockProduits[3] = 25

	// Modification prix
	p := inventaireProduits[1]
	p.Prix = 799.99
	inventaireProduits[1] = p

	// Mise à jour stock
	stockProduits[2] = 60

	fmt.Println("\nTous les produits:")
	afficherProduits(inventaireProduits, stockProduits)

	// Recherche existant
	fmt.Println("\nRecherche ID 2:")
	prod, stock, ok := obtenirProduit(2, inventaireProduits, stockProduits)
	if ok {
		fmt.Printf("%s (Stock: %d)\n", prod.Nom, stock)
	}

	// Recherche inexistant
	fmt.Println("Recherche ID 999:")
	_, _, ok = obtenirProduit(999, inventaireProduits, stockProduits)
	if !ok {
		fmt.Println("Produit non trouvé")
	}

	// Suppression produit 3
	delete(inventaireProduits, 3)
	delete(stockProduits, 3)
	fmt.Println("\nProduit 3 supprimé")

	// Vente
	fmt.Printf("\nStock produit 1 avant: %d\n", stockProduits[1])
	if vendreProduit(1, 3, stockProduits) {
		fmt.Printf("Vente réussie - Stock après: %d\n", stockProduits[1])
	}

	// Vente échouée
	if !vendreProduit(1, 20, stockProduits) {
		fmt.Println("Vente échouée - stock insuffisant")
	}

	// Réapprovisionnement
	reapprovisionnerProduit(2, 40, stockProduits)
	fmt.Printf("Réapprov produit 2 - Stock: %d\n", stockProduits[2])

	// PARTIE 3: COMBINAISON SLICES ET MAPS & PERFORMANCE
	fmt.Println("\n=== PARTIE 3: PERFORMANCE ===")

	// Index par catégorie
	produitsParCategorie := make(map[string][]int)
	for id, prod := range inventaireProduits {
		produitsParCategorie[prod.Categorie] = append(produitsParCategorie[prod.Categorie], id)
	}

	fmt.Println("\nProduits par catégorie:")
	for cat, ids := range produitsParCategorie {
		fmt.Printf("%s: ", cat)
		listerProduitsParCategorie(cat, inventaireProduits, produitsParCategorie)
	}

	// Performance grand volume sans capacité
	fmt.Println("\nAjout 100k produits SANS capacité initiale:")
	start := time.Now()
	inv1 := make(map[int]Produit)
	stock1 := make(map[int]int)
	for i := 1; i <= 100000; i++ {
		inv1[i] = Produit{i, fmt.Sprintf("P%d", i), float64(i%100), "Cat"}
		stock1[i] = i % 50
	}
	temps1 := time.Since(start)
	fmt.Printf("Temps: %v\n", temps1)

	// Performance grand volume AVEC capacité
	fmt.Println("\nAjout 100k produits AVEC capacité 100000:")
	start = time.Now()
	inv2 := make(map[int]Produit, 100000)
	stock2 := make(map[int]int, 100000)
	for i := 1; i <= 100000; i++ {
		inv2[i] = Produit{i, fmt.Sprintf("P%d", i), float64(i%100), "Cat"}
		stock2[i] = i % 50
	}
	temps2 := time.Since(start)
	fmt.Printf("Temps: %v\n", temps2)

	fmt.Printf("\nGain: %v (%.2f%%)\n", temps1-temps2, float64(temps1-temps2)*100/float64(temps1))
	fmt.Println("→ make avec capacité = moins de réallocations = plus rapide")

	// Recherche 10k aléatoires
	fmt.Println("\n10 000 recherches aléatoires:")
	start = time.Now()
	for i := 0; i < 10000; i++ {
		_ = inv2[rand.Intn(100000)+1]
	}
	tempsRecherche := time.Since(start)
	fmt.Printf("Temps: %v\n", tempsRecherche)

	// Itération complète
	fmt.Println("\nItération 100k produits:")
	start = time.Now()
	count := 0
	for _ = range inv2 {
		count++
	}
	tempsIter := time.Since(start)
	fmt.Printf("Temps: %v\n", tempsIter)

	fmt.Println("\nConclusions:")
	fmt.Println("- Les maps: accès O(1), recherche très rapide")
	fmt.Println("- Capacité initiale améliore ajout massif (évite réallocations)")
	fmt.Println("- Itération sur map est efficace mais plus lente que recherche directe")
}

func categorieExiste(nom string, categories []string) bool {
	for _, category := range categories {
		if category == nom {
			return true
		}
	}
	return false
}

func supprimerCategorie(nom string, categories []string) []string {
	var resultat []string

	for _, categorie := range categories {
		if categorie != nom {
			resultat = append(resultat, categorie)
		}
	}

	return resultat
}

func obtenirProduit(id int, inventaire map[int]Produit, stock map[int]int) (Produit, int, bool) {
	produit, ok := inventaire[id]
	if !ok {
		return Produit{}, 0, false
	}
	return produit, stock[id], true
}

func vendreProduit(id int, quantite int, stock map[int]int) bool {
	if stock[id] >= quantite {
		stock[id] -= quantite
		return true
	}
	return false
}

func reapprovisionnerProduit(id int, quantite int, stock map[int]int) {
	stock[id] += quantite
}

func afficherProduits(inventaire map[int]Produit, stock map[int]int) {
	for id, prod := range inventaire {
		fmt.Printf("ID %d: %s - %.2f€ (%s) - Stock: %d\n", id, prod.Nom, prod.Prix, prod.Categorie, stock[id])
	}
}

func listerProduitsParCategorie(categorie string, inventaire map[int]Produit, produitsParCategorie map[string][]int) {
	ids := produitsParCategorie[categorie]
	for _, id := range ids {
		fmt.Printf("%s ", inventaire[id].Nom)
	}
	fmt.Println()
}
