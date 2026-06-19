package main

import (
	"fmt"
	"sync"
	"time"
)

type Result struct {
	number int
	sum    int
}

// sumDivisors calcule la somme de tous les diviseurs d'un nombre n.
// Par exemple, pour n=6, les diviseurs sont 1, 2, 3, 6, et la somme est 12.
func sumDivisors(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		if n%i == 0 {
			sum += i
		}
	}
	return sum
}

func generateNumbers(numJobs int, jobs chan<- int) {
	for n := 1; n <= numJobs; n++ {
		jobs <- n
	}
	close(jobs)
}

func worker(id int, jobs <-chan int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		sum := sumDivisors(job)
		results <- Result{number: job, sum: sum}
	}
}

func main() {
	const numWorkers = 4
	const numJobs = 100

	startTime := time.Now()

	jobs := make(chan int, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup

	go generateNumbers(numJobs, jobs)

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	total := 0
	for result := range results {
		fmt.Printf("Somme des diviseurs de %d = %d\n", result.number, result.sum)
		count++
		total += result.sum
	}

	duration := time.Since(startTime)

	fmt.Println()
	fmt.Printf("%d tâches traitées par %d workers.\n", count, numWorkers)
	fmt.Printf("Somme cumulée de tous les résultats : %d\n", total)
	fmt.Printf("Temps d'exécution total : %s\n", duration)
}