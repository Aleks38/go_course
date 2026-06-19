package main

import (
	"context"
	"fmt"
	"time"
)

func effectuerOperationLongue(ctx context.Context, id string) error {
	fmt.Printf("[%s] Début de l'opération...\n", id)

	for i := 1; i <= 5; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("[%s] Opération annulée : %v\n", id, ctx.Err())
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
			fmt.Printf("[%s] Traitement étape %d...\n", id, i)
		}
	}

	fmt.Printf("[%s] Opération terminée avec succès.\n", id)
	return nil
}

func main() {
	fmt.Println("Démarrage du programme principal.")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resultChan := make(chan error, 1)

	go func() {
		err := effectuerOperationLongue(ctx, "Ma Tâche")
		resultChan <- err
	}()

	select {
	case err := <-resultChan:
		if err != nil {
			fmt.Printf("Main: L'opération s'est terminée avec une erreur : %v\n", err)
		} else {
			fmt.Println("Main: L'opération s'est terminée avec succès avant le timeout.")
		}
	case <-ctx.Done():
		fmt.Printf("Main: Timeout atteint ou annulation : %v\n", ctx.Err())
	}

	fmt.Println("Fin du programme principal.")
}
