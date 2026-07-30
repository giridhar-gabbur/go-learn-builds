package main

import (
	"context"
	"fmt"
	"time"
)

func doWork (ctx context.Context, name string){
	for{
		select {
		case <- ctx.Done():
			fmt.Printf("%s: stopped, reason: %v\n",name, ctx.Err())
			return
			
		case <- time.After(500*time.Millisecond):
			fmt.Printf("%s working....\n", name)
		}
	}
}

func fetchUserFromDB (ctx context.Context)(string, error){
	for {
		select {
		case <- time.After(1*time.Second):
			return "Giridhar", nil
		case <- ctx.Done():
			return "", ctx.Err()
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go doWork(ctx, "worker1")
	go doWork(ctx, "worker2")

	time.Sleep(2*time.Second)

	fmt.Println("Cancelling....")
	cancel()

	ctx1, cancel1 := context.WithTimeout(
		context.Background(),
		1*time.Second,
	)
	defer cancel1()

	user, err := fetchUserFromDB(ctx1)
	if err != nil{
		fmt.Println("failed: ", err)
		return
	}
	fmt.Println("Got user: ", user)
}