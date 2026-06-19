package main

import (
	"fmt"
	"math/rand" // Pour générer des délais aléatoires
	"time"      // Pour time.Sleep et time.NewTicker
)

// dataProducer envoie des mesures régulières sur dataChan.
// Il écoute aussi sur quitChan pour s'arrêter proprement.
func dataProducer(dataChan chan<- string, quitChan <-chan struct{}) {
	defer fmt.Println("[PRODUCER DATA] Arrêt de la goroutine de production de données.")
	for {
		select {
		case <-quitChan: // Reçoit le signal d'arrêt
			return // Quitte la goroutine
		case <-time.After(time.Duration(rand.Intn(3)+1) * time.Second): // Envoie toutes les 1 à 3 secondes
			dataChan <- fmt.Sprintf("Température: %.1f°C (Humidité: %.2f%%)", 20.0+rand.Float64()*10, 50.0+rand.Float64()*20)
		}
	}
}

// alertProducer envoie des alertes critiques sur alertChan.
// Il écoute aussi sur quitChan pour s'arrêter proprement.
func alertProducer(alertChan chan<- string, quitChan <-chan struct{}) {
	defer fmt.Println("[PRODUCER ALERT] Arrêt de la goroutine de production d'alertes.")
	for {
		select {
		case <-quitChan: // Reçoit le signal d'arrêt
			return // Quitte la goroutine
		case <-time.After(time.Duration(rand.Intn(6)+5) * time.Second): // Envoie toutes les 5 à 10 secondes
			alertChan <- "Niveau critique atteint! Intervention requise!"
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Initialisation du générateur de nombres aléatoires

	fmt.Println("--- Système de Surveillance Démarré ---")

	// Étape 1 : Définition des Channels
	dataChannel := make(chan string)
	alertChannel := make(chan string)
	quitChannel := make(chan struct{}) // Channel pour signaler l'arrêt aux producteurs et au moniteur

	// time.NewTicker crée un Ticker qui envoie des signaux sur son channel .C à intervalles réguliers.
	ticker := time.NewTicker(2 * time.Second) // Pour les vérifications de statut périodiques
	defer ticker.Stop()                       // S'assure que le ticker est arrêté à la fin de main pour libérer les ressources

	// Étape 2 : Goroutines Productrices d'Événements
	go dataProducer(dataChannel, quitChannel)
	go alertProducer(alertChannel, quitChannel)

	// Étape 4 : Déclenchement de l'Arrêt
	// Goroutine pour envoyer le signal d'arrêt après un délai (ici 15 secondes)
	go func() {
		time.Sleep(15 * time.Second)
		fmt.Println("\n[MONITOR] Envoi du signal d'arrêt aux producteurs et au moniteur...")
		// Fermer le canal est une manière idiomatique de signaler la fin.
		// Tous les récepteurs qui écoutent ce canal avec `<-chan` ou `for range`
		// recevront la "zero value" une fois le canal vide et fermé, ou sortiront de la boucle `for range`.
		close(quitChannel)
	}()

	// Étape 3 : La Boucle select du Moniteur
	fmt.Println("[MONITOR] En attente d'événements...")
	for { // Boucle infinie pour écouter les événements continuellement
		select {
		case data := <-dataChannel: // Un message est reçu sur dataChannel
			fmt.Printf("[MESURE] %s\n", data)
		case alert := <-alertChannel: // Un message est reçu sur alertChannel
			fmt.Printf("[ALERTE CRITIQUE] !!! %s !!!\n", alert)
		case <-ticker.C: // Un signal est reçu du ticker (vérification de statut)
			fmt.Println("[STATUS] Vérification système...")
		case <-quitChannel: // Un signal est reçu sur quitChannel
			fmt.Println("[MONITOR] Signal d'arrêt reçu. Arrêt du système.")
			return // Sort de la boucle et de la fonction main, terminant le programme
		}
	}
}

