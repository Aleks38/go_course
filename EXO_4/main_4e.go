package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Ressources partagées globales
var compteur int         // le compteur partagé (int classique)
var compteurAtomic int64 // version pour sync/atomic (doit être int64)
var mu sync.Mutex        // le verrou qui protège l'accès à "compteur"

const (
	nbGoroutines           = 100
	incrementsParGoroutine = 1000
)

// ÉTAPE 1 : incrémentation SANS protection -> race condition.
func incrementerCompteurNonSynchro(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < incrementsParGoroutine; i++ {
		compteur++ // accès concurrent NON protégé : des incréments seront perdus
	}
}

// ÉTAPE 2 : incrémentation protégée par un Mutex.
func incrementerCompteurSynchro(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < incrementsParGoroutine; i++ {
		mu.Lock()   // une seule goroutine entre ici à la fois
		compteur++  // section critique : courte (juste l'incrément)
		mu.Unlock() // on libère immédiatement pour que les autres avancent
	}
}

// ÉTAPE 3 : alternative atomique (sans mutex) pour une simple incrémentation.
func incrementerCompteurAtomic(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < incrementsParGoroutine; i++ {
		atomic.AddInt64(&compteurAtomic, 1) // opération atomique gérée par le CPU
	}
}

func main() {
	attendu := nbGoroutines * incrementsParGoroutine
	fmt.Printf("Résultat attendu : %d\n\n", attendu)

	// ---------- ÉTAPE 1 : version NON synchronisée ----------
	compteur = 0
	var wg1 sync.WaitGroup
	start := time.Now()
	for i := 0; i < nbGoroutines; i++ {
		wg1.Add(1)
		go incrementerCompteurNonSynchro(&wg1)
	}
	wg1.Wait()
	fmt.Printf("[NON SYNCHRO] compteur = %d (perdu %d) en %v\n",
		compteur, attendu-compteur, time.Since(start))

	// ---------- ÉTAPE 2 : version avec Mutex ----------
	compteur = 0
	var wg2 sync.WaitGroup
	start = time.Now()
	for i := 0; i < nbGoroutines; i++ {
		wg2.Add(1)
		go incrementerCompteurSynchro(&wg2)
	}
	wg2.Wait()
	fmt.Printf("[MUTEX]       compteur = %d en %v\n",
		compteur, time.Since(start))

	// ---------- ÉTAPE 3 : version atomic ----------
	atomic.StoreInt64(&compteurAtomic, 0)
	var wg3 sync.WaitGroup
	start = time.Now()
	for i := 0; i < nbGoroutines; i++ {
		wg3.Add(1)
		go incrementerCompteurAtomic(&wg3)
	}
	wg3.Wait()
	fmt.Printf("[ATOMIC]      compteur = %d en %v\n",
		atomic.LoadInt64(&compteurAtomic), time.Since(start))
}
