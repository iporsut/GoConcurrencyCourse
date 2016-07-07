package main

import (
	"fmt"
	"sync"
)

func distribute(in <-chan int, outs ...chan<- int) {
	go func() {
		var wg sync.WaitGroup
		outbuf := make(chan chan<- int, len(outs))
		for _, out := range outs {
			outbuf <- out
		}
		for i := range in {
			wg.Add(1)
			go func(i int) {
				out := <-outbuf
				out <- i
				outbuf <- out
				wg.Done()
			}(i)
		}
		wg.Wait()
		for _, out := range outs {
			close(out)
		}
	}()
}

func broadcast(in <-chan int, outs ...chan<- int) {
	go func() {
		for i := range in {
			for _, out := range outs {
				out <- i
			}
		}
		for _, out := range outs {
			close(out)
		}
	}()
}

func merge(ins ...<-chan int) chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(ins))
	for _, in := range ins {
		go func(in <-chan int) {
			for i := range in {
				out <- i
			}
			wg.Done()
		}(in)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	/*for i := range square(merge(gen(1, 2), gen(3, 4))) {
		fmt.Println(i)
	}*/

	out1 := make(chan int)
	out2 := make(chan int)
	out3 := make(chan int)
	// broadcast(gen(3, 5), out1, out2)
	distribute(gen(3,10), out1, out2, out3)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		for i := range square(square(out1)) {
			fmt.Println("^4",i)
		}
		wg.Done()
	}()
	go func() {
		for i := range square(out2) {
			fmt.Println("^2",i)
		}
		wg.Done()
	}()
	go func() {
		for i := range out3 {
			fmt.Println(i)
		}
		wg.Done()
	}()
	wg.Wait()
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
