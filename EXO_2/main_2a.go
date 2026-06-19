package main

import "fmt"

func main() {
	// Exercice 1
	var nomUtilisateur string = "42"
	var ageUtilisateur uint = 255
	var estConnecte bool = true
	var soldeCompte float64 = 255

	fmt.Printf("nomUtilisateur = %v, type = %T\n", nomUtilisateur, nomUtilisateur)
	fmt.Printf("ageUtilisateur = %v, type = %T\n", ageUtilisateur, ageUtilisateur)
	fmt.Printf("estConnecte = %v, type = %T\n", estConnecte, estConnecte)
	fmt.Printf("soldeCompte = %v, type = %T\n", soldeCompte, soldeCompte)

	// Exercice 2
	villeResidence := "Paris"
	codePostal := 75000
	tauxRemise := float32(0.15)

	fmt.Printf("villeResidence = %v, type = %T\n", villeResidence, villeResidence)
	fmt.Printf("codePostal = %v, type = %T\n", codePostal, codePostal)
	fmt.Printf("tauxRemise = %v, type = %T\n", tauxRemise, tauxRemise)

	// Exercice 3
	const PI = 3.14159
	const NOM_APPLICATION = "Gestionnaire Go"
	const ANNEE_LANCEMENT = 2023

	rayon := 10.5
	circonference := 2 * PI * rayon

	fmt.Printf("Pour un rayon de %.1f, la circonférence est %.2f\n", rayon, circonference)
	fmt.Printf("Nom de l'application : %s, Année de lancement : %d\n", NOM_APPLICATION, ANNEE_LANCEMENT)

	// Exercice 4
	ageUtilisateur = 30
	var message string
	var compteur int
	var estValide bool
	var prix float64

	fmt.Printf("Valeur par défaut de 'message' (string non initialisée) : '%s'\n", message)
	fmt.Printf("Valeur par défaut de 'compteur' (int non initialisé) : %d\n", compteur)
	fmt.Printf("Valeur par défaut de 'estValide' (bool non initialisé) : %t\n", estValide)
	fmt.Printf("Valeur par défaut de 'prix' (float64 non initialisé) : %.1f\n", prix)
	fmt.Println("")
}

