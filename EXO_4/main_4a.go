package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func effectuerTache(id int, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	fmt.Printf("Goroutine %d: Début de la tâche...\n", id)

	d := time.Duration(50+rand.Intn(451)) * time.Millisecond
	time.Sleep(d)

	fmt.Printf("Goroutine %d: Tâche terminée.\n", id)
}

func effectuerTacheAvecCanal(id int, wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()

	fmt.Printf("Goroutine %d: Début de la tâche...\n", id)

	d := time.Duration(50+rand.Intn(451)) * time.Millisecond
	time.Sleep(d)

	fmt.Printf("Goroutine %d: Tâche terminée.\n", id)

	ch <- fmt.Sprintf("Goroutine %d a terminé avec succès.", id)
}

func travailleur(id int, taches <-chan int, resultats chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for tache := range taches {
		fmt.Printf("Travailleur %d: traite la tâche %d...\n", id, tache)

		d := time.Duration(50+rand.Intn(451)) * time.Millisecond
		time.Sleep(d)

		resultats <- fmt.Sprintf("Tâche %d traitée par le travailleur %d", tache, id)
	}
}

func main() {
	fmt.Println("EXERCICE 1")
	fmt.Println("////////////////////////////////")
	fmt.Println()

	rand.Seed(time.Now().UnixNano())


	for i := 1; i <= 5; i++ {
		id := i
		go effectuerTache(id, nil)
	}
	fmt.Println("Toutes les goroutines lancées.")

	fmt.Println()
	time.Sleep(600 * time.Millisecond)

	// Exercice 1 - Réponse :
	// Non, rien ne permet de savoir si toutes les goroutines ont terminé.
	// Une goroutines est non bloquantes, donc main lance toutes les goroutines quasiment en même temps.
	// Mais une fois que toutes les goroutines sont lancés rien ne dit d'attendre la fin de leur travail.
	// Donc main rend la main une fois tout lancé.

	fmt.Println()
	fmt.Println("EXERCICE 2")
	fmt.Println("////////////////////////////////")
	fmt.Println()

	var wg2 sync.WaitGroup


	for i := 6; i <= 10; i++ {
		id := i
		wg2.Add(1)
		go effectuerTache(id, &wg2)
	}
	wg2.Wait()
	fmt.Println("Toutes les Goroutines sont terminés")

	// Exercice 2 - Réponse :
	// Oui le comportement a changé
	// Maintenant toutes les goroutines terminent bien leur travail grâce au wg.Wait()
	// Il se bloque tant que les 5 wg.Done() n'ont pas été appelés.

	fmt.Println()
	fmt.Println("EXERCICE 3")
	fmt.Println("////////////////////////////////")
	fmt.Println()

	var wg3 sync.WaitGroup

	ch := make(chan string, 5)

	for i := 1; i <= 5; i++ {
		id := i
		wg3.Add(1)
		go effectuerTacheAvecCanal(id, &wg3, ch)
	}

	wg3.Wait()
	close(ch)

	fmt.Println("\nRésultats reçus :")
	for msg := range ch {
		fmt.Println(msg)
	}

	// Exercice 3 - Réponse :
	// L'ordre d'affichage des messages de fin de tâche est aléatoire
	// Non l'ordre de résultat ne correspond pas à leurs id, mais dans l'ordre du time.Sleep aléatoire.

	fmt.Println()
	fmt.Println("EXERCICE 4")
	fmt.Println("////////////////////////////////")
	fmt.Println()

	var wg4 sync.WaitGroup

	const nbTravailleurs = 3
	const nbTaches = 10

	taches := make(chan int, nbTravailleurs)
	resultats := make(chan string, nbTaches)

	for i := 1; i <= nbTravailleurs; i++ {
		wg4.Add(1)
		go travailleur(i, taches, resultats, &wg4)
	}

	for t := 1; t <= nbTaches; t++ {
		taches <- t
	}
	close(taches)

	wg4.Wait()
	close(resultats)

	fmt.Println("\nRésultats reçus :")
	for msg := range resultats {
		fmt.Println(msg)
	}

	// Exercice 4 - Réponse :
	// L'ordre d'affichage est aléatoire, c'est en fonction du temps aléatoire du time.Sleep.
	// Les taches sont parallélisées en fonction du nombre de travailleurs, car le nombre de goroutines lancées correspond au nombre de travailleurs
	// Plus, il y a de travailleur plus le temps d'exécution global sera petit et inversement.
}