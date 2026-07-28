package main

import (
	"fmt"
	"sync"
)

func counter() {
	var wg sync.WaitGroup
	var mu sync.RWMutex
	counter := 0

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			counter++
		} ()
	}
	wg.Wait()
	fmt.Println("Counter:", counter)
}

func main() {
	counter()
}