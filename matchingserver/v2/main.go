package main

import (
	"io"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	waitingChan := make(chan net.Conn, 1)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		select {
		case waitingChan <- conn:
		case waitCon := <-waitingChan:
			go CopyHandler(waitCon, conn)
		}
	}
}

func CopyHandler(conn1 net.Conn, conn2 net.Conn) {
	defer conn1.Close()
	defer conn2.Close()
	done := make(chan struct{}, 1)
	go func() {
		io.Copy(conn1, conn2)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(conn2, conn1)
		done <- struct{}{}
	}()
	<-done
}
