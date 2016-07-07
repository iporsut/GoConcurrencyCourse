package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	in := make(chan struct{})

	proxyPassChans := make([]chan struct{}, len(os.Args[2:]))
	proxyOutChans := make([]chan string, len(os.Args[2:]))
	for i, url := range os.Args[2:] {
		proxyPassChans[i] = make(chan struct{})
		proxyOutChans[i] = ProxyPass(proxyPassChans[i], url)
	}
	distribute(in, proxyPassChans...)
	out := merge(proxyOutChans...)
	http.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			in <- struct{}{}
		}()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, <-out)
	})
	http.ListenAndServe(os.Args[1], nil)
}

func ProxyPass(in chan struct{}, url string) chan string {
	out := make(chan string)
	go func() {
		for _ = range in {
			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			resp.Body.Close()
			out <- string(b)
		}
		close(out)
	}()
	return out
}

func distribute(in <-chan struct{}, outs ...chan struct{}) {
	go func() {
		var wg sync.WaitGroup
		outbuf := make(chan chan<- struct{}, len(outs))
		for _, out := range outs {
			outbuf <- out
		}
		for _ = range in {
			wg.Add(1)
			go func() {
				out := <-outbuf
				out <- struct{}{}
				outbuf <- out
				wg.Done()
			}()
		}
		wg.Wait()
		for _, out := range outs {
			close(out)
		}
	}()
}

func merge(ins ...chan string) chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(ins))
	for _, in := range ins {
		go func(in <-chan string) {
			for i := range in {
				out <- i
			}
			wg.Done()
		}(in)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
