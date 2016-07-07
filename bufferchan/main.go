package main

import (
	"fmt"
	"math/rand"
	"time"
)

func Download(done chan bool, limit chan int) {
	delay := rand.Intn(10)
	limit <- delay
	fmt.Println("Delay", delay)
	time.Sleep(time.Duration(delay) * time.Second)
	<-limit
	done <- true
}

func main() {
	rand.Seed(time.Now().UnixNano())
	done := make(chan bool)
	downloadLimit := make(chan int, 4)
	for i:= 1; i <= 1000; i++ {
		go Download(done, downloadLimit)
	}
	for i:= 1; i <= 1000; i++ {
		<-done
	}
}
