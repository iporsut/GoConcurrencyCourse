package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

func main() {
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	waitingChan := make(chan *MatchingHandler, 1)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		mh := &MatchingHandler{}
		mh.SetConn1(conn)
		select {
		case waitingChan <- mh:
			log.Println("New Waiting")
			mh.SetWaitingChan(waitingChan)
			go mh.Conn1Handler()
		case mh := <-waitingChan:
			log.Println("Match")
			mh.SetConn2(conn)
			go mh.Conn2Handler()
			waitingChan = make(chan *MatchingHandler, 1)
		}
	}
}

type MatchingHandler struct {
	conn1       net.Conn
	conn2       net.Conn
	waitingChan chan *MatchingHandler

	muConn1       sync.RWMutex
	muConn2       sync.RWMutex
	muWaitingChan sync.RWMutex
}

func (h *MatchingHandler) SetConn1(conn net.Conn) {
	h.muConn1.Lock()
	defer h.muConn1.Unlock()
	h.conn1 = conn
}

func (h *MatchingHandler) SetConn2(conn net.Conn) {
	h.muConn2.Lock()
	defer h.muConn2.Unlock()
	h.conn2 = conn
}

func (h *MatchingHandler) SetWaitingChan(ch chan *MatchingHandler) {
	h.muWaitingChan.Lock()
	defer h.muWaitingChan.Unlock()
	h.waitingChan = ch
}

func (h *MatchingHandler) Conn1Handler() {
	defer func() {
		h.muWaitingChan.RLock()
		if len(h.waitingChan) > 0 {
			<-h.waitingChan
			log.Println("Flush waiting channel")
		}
		h.muWaitingChan.RUnlock()

		h.muConn1.RLock()
		if h.conn1 != nil {
			h.conn1.Close()
		}
		h.muConn1.RUnlock()

		h.muConn2.RLock()
		if h.conn2 != nil {
			log.Println("Close 2 from 1")
			h.conn2.Close()
		}
		h.muConn2.RUnlock()
	}()

	h.muConn1.RLock()
	scanner := bufio.NewScanner(h.conn1)
	h.muConn1.RUnlock()

	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)

		h.muConn2.RLock()
		if h.conn2 != nil {
			_, err := fmt.Fprintln(h.conn2, text)
			if err != nil {
				h.muConn2.RUnlock()
				break
			}
		}
		h.muConn2.RUnlock()
	}
	log.Println("Disconnected")
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func (h *MatchingHandler) Conn2Handler() {
	defer func() {
		h.muConn2.RLock()
		if h.conn2 != nil {
			h.conn2.Close()
		}
		h.muConn2.RUnlock()

		h.muConn1.RLock()
		if h.conn1 != nil {
			h.conn1.Close()
		}
		h.muConn1.RUnlock()
	}()

	h.muConn2.RLock()
	scanner := bufio.NewScanner(h.conn2)
	h.muConn2.RUnlock()

	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)

		h.muConn1.RLock()
		if h.conn1 != nil {
			_, err := fmt.Fprintln(h.conn1, text)
			if err != nil {
				h.muConn1.RUnlock()
				break
			}
		}
		h.muConn1.RUnlock()

	}
	log.Println("Disconnected")
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
