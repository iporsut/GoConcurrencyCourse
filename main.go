package main

import (
	"fmt"
	"time"
)

func HelloWorld() {
	for i := 1; i <= 10; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("Hello World")
	}
}

func Sawasdee() {
	panic("PANIC!")
	for i := 1; i <= 10; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("สวัสดี")
	}
}

func SendValue(ch chan int) {
	ch <- 10
}
func main() {
	var chInt chan int = make(chan int)
	go SendValue(chInt)
	i := <-chInt
	fmt.Println(i)
}
