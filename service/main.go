package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func RandomHandler(w http.ResponseWriter, r *http.Request) {
	n := rand.Intn(10)
	time.Sleep(time.Duration(n) * time.Second)
	w.Header().Set("Content-Type", "application/json")
	log.Println("Number", n)
	fmt.Fprintf(w, `{"number":%d}`, n)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/random", RandomHandler)
	http.ListenAndServe(os.Args[1], nil)
}
