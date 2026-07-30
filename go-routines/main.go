package main 

import (
	"sync"
	"fmt"
)

func sayHello(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Hello %s!\n", name)
}

func main(){
	var wg sync.WaitGroup

	wg.Add(1)
	go sayHello("Alice", &wg)

	wg.Add(1)
	go sayHello("Giridhar", &wg)

	wg.Wait()
	fmt.Println("All goroutines done!")
}