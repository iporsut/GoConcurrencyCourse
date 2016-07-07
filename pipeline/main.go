package main

import "fmt"

func gen(from int, to int) chan int {
	out := make(chan int)
	go func() {
		for i := from; i <= to; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) chan int {
	out := make(chan int)
	go func() {
		for i := range in {
			out <- i * i
		}
		close(out)
	}()
	return out
}

func add(in <-chan int) chan int {
	out := make(chan int)
	go func() {
		match := make(chan int, 1)
		for i := range in {
			select {
			case match <- i:
			case left := <-match:
				out <- left + i
			}
		}
		close(out)
	}()

	return out
}

func add2(left <-chan int, right <-chan int) chan int {
	out := make(chan int)
	go func() {
		for {
			l, lopen := <- left
			r, ropen := <- right
			if lopen && ropen {
				out <- l + r
			} else {
				break
			}
		}
		close(out)
	}()
	return out
}

func main() {
	for i := range add2(gen(1, 8), square(gen(3,10))) {
		fmt.Println(i)
	}
}
