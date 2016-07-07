package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	autoID  int
	connMap map[int]net.Conn = make(map[int]net.Conn)
	mu      sync.Mutex
)

func AddConnection(conn net.Conn) int {
	mu.Lock()
	defer mu.Unlock()
	autoID++
	connMap[autoID] = conn
	return autoID
}

func DeleteConnection(id int) {
	mu.Lock()
	defer mu.Unlock()
	delete(connMap, id)
}

func Broadcast(fromID int, msg string) {
	mu.Lock()
	defer mu.Unlock()
	for id, conn := range connMap {
		if id != fromID {
			go fmt.Fprintln(conn, msg)
		}
	}
}

func main() {
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// go ChatHandler(conn)
		go EchoHandler(conn)
	}
}

func ChatHandler(conn net.Conn) {
	id := AddConnection(conn)
	defer DeleteConnection(id)
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		Broadcast(id, text)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func EchoHandler(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		fmt.Fprintln(conn, text)
	}
	log.Println("Disconnected")
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
