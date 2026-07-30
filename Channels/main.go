package main

import (
	"fmt"
	"sync"
)

func multiplierWorker (id int, numbers <- chan int, results chan <- int, wg *sync.WaitGroup ) {
	defer wg.Done()
	for num := range numbers {
		fmt.Printf("Worker: %d is working on %d\n", id, num)
		num *= 10
		results <- num
	}
}

func main() {
	numbers := make(chan int, 10)
	results := make(chan int, 10)

	var wg sync.WaitGroup

	for w := 0; w < 2; w++ {
		wg.Add(1)
		go multiplierWorker(w, numbers, results, &wg)
	}

	for i := 1; i < 5; i++{
		numbers <- i
	}
	close(numbers)

	go func(){
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}

	fmt.Println("All done! Worker pool shut down cleanly.")
}