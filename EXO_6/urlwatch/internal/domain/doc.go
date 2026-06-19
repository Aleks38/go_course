// Package domain contient les types métier (CheckResult, Batch), les erreurs
// du domaine (sentinelles et erreurs personnalisées) et les interfaces centrales
// (Checker, Store). Les autres packages dépendent de domain (inversion de
// dépendance) ; domain ne dépend d'aucun autre package interne.
package domain
