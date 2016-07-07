package main

import "fmt"

func gen() chan int {
	ch := make(chan int)
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func main() {
	ch := gen()
	for i := range ch {
		fmt.Println(i)
	}
}
