package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

func FirstConn(conn io.ReadWriteCloser, mh *MessageHandler, waitingChan chan *MessageHandler) {
	defer func() {
		if w := mh.Writer(); w != nil {
			w.Close()
		} else {
			if len(waitingChan) > 0 {
				<-waitingChan
			}
		}
		conn.Close()
	}()
	log.Println("New Waiting")
	mh.SetReader(conn)
	mh.StartCopy()

}

func SecondConn(conn io.ReadWriteCloser, mh *MessageHandler) {
	defer func() {
		if w := mh.Writer(); w != nil {
			w.Close()
		}
		conn.Close()
	}()
	log.Println("Match")
	mh.SetWriter(conn)
	w := mh.Reader()
	mh = &MessageHandler{}
	mh.SetReader(conn)
	mh.SetWriter(w)
	mh.StartCopy()
}

func main() {
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	waitingChan := make(chan *MessageHandler, 1)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		mh := &MessageHandler{}
		select {
		case waitingChan <- mh:
			go FirstConn(conn, mh, waitingChan)
		case mh := <-waitingChan:
			go SecondConn(conn, mh)
			waitingChan = make(chan *MessageHandler, 1)
		}
	}
}

type MessageHandler struct {
	r io.ReadWriteCloser
	w io.ReadWriteCloser

	mu sync.RWMutex
}

func (mh *MessageHandler) Reader() io.ReadWriteCloser {
	mh.mu.RLock()
	defer mh.mu.RUnlock()
	return mh.r
}

func (mh *MessageHandler) Writer() io.ReadWriteCloser {
	mh.mu.RLock()
	defer mh.mu.RUnlock()
	return mh.w
}

func (mh *MessageHandler) SetReader(r io.ReadWriteCloser) {
	mh.mu.Lock()
	defer mh.mu.Unlock()
	mh.r = r
}

func (mh *MessageHandler) SetWriter(w io.ReadWriteCloser) {
	mh.mu.Lock()
	defer mh.mu.Unlock()
	mh.w = w
}

func (mh *MessageHandler) StartCopy() {
	r := mh.Reader()
	if r == nil {
		return
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)

		w := mh.Writer()
		if w != nil {
			_, err := fmt.Fprintln(w, text)
			if err != nil {
				break
			}
		}
	}
	log.Println("Disconnected")
}
