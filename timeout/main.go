package main

func Query(ch chan string) {
	time.Sleep(10 * time.Second)
	ch <- "Hello"
}

func main() {
	ch := make(chan string)
	go Query(ch)
	select {
	case s :=<-ch:
		fmt.Println(s)
	case <-time.After(5 * time.Second):
		fmt.Println("Timeout")
	}
}
