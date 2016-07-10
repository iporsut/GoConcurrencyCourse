package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

func main() {
	SocketChat()
}

func WebSocketChat() {
	waitingChan := make(chan *MessageHandler, 1)
	http.Handle("/echo", websocket.Handler(func(ws *websocket.Conn) {
		Handler(ws, waitingChan)
	}))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func SocketChat() {
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
		go Handler(conn, waitingChan)
	}
}

func Handler(conn io.ReadWriteCloser, waitingChan chan *MessageHandler) {
	mh := &MessageHandler{}
	select {
	case waitingChan <- mh:
		mh.FirstConn(conn, waitingChan)
	case mh = <-waitingChan:
		mh.SecondConn(conn)
	}
}

type MessageHandler struct {
	r io.ReadWriteCloser
	w io.ReadWriteCloser

	mu sync.RWMutex
}

func (mh *MessageHandler) FirstConn(conn io.ReadWriteCloser, waitingChan chan *MessageHandler) {
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
	mh.Start()
}

func (mh *MessageHandler) SecondConn(conn io.ReadWriteCloser) {
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
	mh.Start()
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

func (mh *MessageHandler) Start() {
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
