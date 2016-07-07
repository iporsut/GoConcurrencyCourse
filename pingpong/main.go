package main

import "fmt"

func Ping(msgCh chan string, done chan struct{}) {
	for i := 1; i <= 10; i++ {
		msgCh <- "Ping"
		fmt.Println("Recieve", <-msgCh)
	}
	close(done)
}

func Pong(msgCh chan string, done chan struct{}) {
	for i := 1; i <= 10; i++ {
		fmt.Println("Receive", <-msgCh)
		msgCh <- "Pong"
	}
	close(done)
}

func main() {
	msgCh := make(chan string)
	done := make(chan struct{})
	go Ping(msgCh, done)
	go Pong(msgCh, done)
	<-done
	<-done
}
