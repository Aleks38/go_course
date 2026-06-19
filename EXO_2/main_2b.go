package main

import (
	"errors"
	"fmt"
)

func main()  {
	somme, count, moyenne := CalculerStatistiquesBase(10, 20, 30, 40)
	fmt.Printf("Somme: %d, Count: %d, Moyenne: %.2f\n", somme, count, moyenne)

	sommeVide, countVide, moyenneVide := CalculerStatistiquesBase()
	fmt.Printf("Somme (vide): %d, Count (vide): %d, Moyenne (vide): %.2f\n", sommeVide, countVide, moyenneVide)

	min, max, sum, avg, count, err := CalculerStatistiquesCompletes(1.5, 2.8, 0.7, 3.1)
	if err != nil {
		fmt.Println("Erreur:", err)
	} else {
		fmt.Printf("Min: %.2f, Max: %.2f, Somme: %.2f, Moyenne: %.2f, Count: %d\n", min, max, sum, avg, count)
	}

	_, _, _, _, _, errVide := CalculerStatistiquesCompletes()
	if errVide != nil {
		fmt.Println("Erreur pour arguments vides:", errVide)
	}

	minTemp, maxTemp, avgTemp, validCnt, invalidCnt, err := AnalyserDonneesCapteur(22.5, 23.1, -5.0, 101.0, 21.9, 0.0, 24.0)
	if err != nil {
		fmt.Println("Erreur d'analyse:", err)
	} else {
		fmt.Printf("Temp Min: %.2f, Max: %.2f, Moyenne: %.2f, Valides: %d, Invalides: %d\n", minTemp, maxTemp, avgTemp, validCnt, invalidCnt)
	}

	_, _, _, _, _, errToutInvalide := AnalyserDonneesCapteur(-10.0, 105.0, 0.0)
	if errToutInvalide != nil {
		fmt.Println("Erreur pour données toutes invalides:", errToutInvalide)
	}
}

func CalculerStatistiquesBase(nums ...int) (int, int, float64) {
	total := 0
	for _, n := range nums {
		total += n
	}
	if len(nums) == 0 {
		return 0, 0, 0.0
	}
	return total, len(nums), float64(total) / float64(len(nums))
}

func CalculerStatistiquesCompletes(nums ...float64) (float64,float64,float64,float64,int,error) {
	if len(nums) == 0 {
		return 0, 0, 0, 0, 0, errors.New("aucun argument fourni")
	}

	min := nums[0]
	max := nums[0]
	sum := 0.0
	for _, n := range nums {
		sum += n
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	count := len(nums)
	moyenne := sum / float64(count)
	return min, max, sum, moyenne, count, nil
}

func AnalyserDonneesCapteur(nums ...float64) (float64, float64, float64, int, int, error) {
	var validReadings []float64
	invalidCount := 0

	for _, n := range nums {
		if n > 0.0 && n <= 100.0 {
			validReadings = append(validReadings, n)
		} else {
			invalidCount++
		}
	}

	if len(validReadings) == 0 {
		return 0, 0, 0, 0, invalidCount, errors.New("aucun relevé valide trouvé après filtrage")
	}

	minVal, maxVal, _, avg, count, _ := CalculerStatistiquesCompletes(validReadings...)

	return minVal, maxVal, avg, count, invalidCount, nil
}