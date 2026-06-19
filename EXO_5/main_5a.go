package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Stockage en mémoire
var items = []Item{
	{ID: "1", Name: "Clavier", Description: "Clavier mécanique"},
	{ID: "2", Name: "Souris", Description: "Souris sans fil"},
}

var nextID = 3 // pour générer un ID unique au nouvel élément

// itemsHandler gère "/items" : GET (liste) et POST (création).
func itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)

	case http.MethodPost:
		var nouvel Item
		if err := json.NewDecoder(r.Body).Decode(&nouvel); err != nil {
			http.Error(w, "JSON invalide", http.StatusBadRequest)
			return
		}

		nouvel.ID = fmt.Sprintf("%d", nextID) // génération de l'ID
		nextID++
		items = append(items, nouvel)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(nouvel)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func itemByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/items/")
	if id == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		for _, it := range items {
			if it.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(it)
				return
			}
		}
		http.Error(w, "Élément non trouvé", http.StatusNotFound)

	case http.MethodPut:
		var maj Item
		if err := json.NewDecoder(r.Body).Decode(&maj); err != nil {
			http.Error(w, "JSON invalide", http.StatusBadRequest)
			return
		}
		for i := range items {
			if items[i].ID == id {
				items[i].Name = maj.Name
				items[i].Description = maj.Description
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(items[i])
				return
			}
		}
		http.Error(w, "Élément non trouvé", http.StatusNotFound)

	case http.MethodDelete:
		for i := range items {
			if items[i].ID == id {
				// Suppression : on recolle l'avant et l'après de l'index i.
				items = append(items[:i], items[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "Élément non trouvé", http.StatusNotFound)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/items", itemsHandler)
	mux.HandleFunc("/items/", itemByIDHandler)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
