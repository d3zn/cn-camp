package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 10)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	// producer
	go func() {
		for range ticker.C {
			ch <- 1
		}
	}()
	//consumer
	for {
		fmt.Printf("%d\n", <-ch)
	}
}
