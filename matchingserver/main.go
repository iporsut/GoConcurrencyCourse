package main

import (
	"bufio"
	"fmt"
	"io"
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

	waitingChan := make(chan *CopyPipeline, 1)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		cp := &CopyPipeline{}
		select {
		case waitingChan <- cp:
			log.Println("New Waiting")
			go func(conn net.Conn, cp *CopyPipeline, waitingChan chan *CopyPipeline) {
				defer func() {
					if w := cp.Writer(); w != nil {
						w.Close()
					} else {
						if len(waitingChan) > 0 {
							<-waitingChan
						}
					}
					conn.Close()
				}()
				cp.SetReader(conn)
				cp.StartCopy()
			}(conn, cp, waitingChan)
		case cp := <-waitingChan:
			go func(conn net.Conn, cp *CopyPipeline) {
				defer func() {
					if w := cp.Writer(); w != nil {
						w.Close()
					}
					conn.Close()
				}()
				log.Println("Match")
				cp.SetWriter(conn)
				w := cp.Reader()
				cp = &CopyPipeline{}
				cp.SetReader(conn)
				cp.SetWriter(w)
				cp.StartCopy()
			}(conn, cp)
			waitingChan = make(chan *CopyPipeline, 1)
		}
	}
}

type CopyPipeline struct {
	r io.ReadWriteCloser
	w io.ReadWriteCloser

	mu sync.RWMutex
}

func (cp *CopyPipeline) Reader() io.ReadWriteCloser {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.r
}

func (cp *CopyPipeline) Writer() io.ReadWriteCloser {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.w
}

func (cp *CopyPipeline) SetReader(r io.ReadWriteCloser) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.r = r
}

func (cp *CopyPipeline) SetWriter(w io.ReadWriteCloser) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.w = w
}

func (cp *CopyPipeline) StartCopy() {
	r := cp.Reader()
	if r == nil {
		return
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)

		w := cp.Writer()
		if w != nil {
			_, err := fmt.Fprintln(w, text)
			if err != nil {
				break
			}
		}
	}
	log.Println("Disconnected")
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
