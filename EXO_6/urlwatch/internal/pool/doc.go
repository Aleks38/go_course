// Package pool implémente le cœur concurrent : un worker pool borné qui vérifie
// un lot d'URLs en parallèle (fan-out vers les workers, fan-in des résultats),
// avec respect du timeout et de l'annulation via context, et sans fuite de
// goroutine.
package pool
