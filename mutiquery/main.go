package main

import (
	"fmt"
	"log"
	"time"
)

type ContactHistory struct {
	SubscriberNumber string
	CustomerName     string
}

func MongoQuery() ([]ContactHistory, error) {
	result := []ContactHistory{
		{"66812345789", "Weerasak"},
		{"66842545789", "Twin"},
	}
	time.Sleep(2 * time.Second)
	return result, nil
}

func SOAPQuery() ([]ContactHistory, error) {
	/*	result := []ContactHistory{
		{"66812367789", "Wingyplus"},
		{"66842511234", "Roofimon"},
	}*/
	time.Sleep(5 * time.Second)
	return nil, fmt.Errorf("Cannot Query SOAP")
}

func SharepointQuery() ([]ContactHistory, error) {
	time.Sleep(5 * time.Second)
	return nil, fmt.Errorf("Cannot Query Sharepoint")
}

type QueryFunc func() ([]ContactHistory, error)

func RunQuery(fn QueryFunc, resultChan chan []ContactHistory, errorChan chan error) {
	result, err := fn()
	if err != nil {
		errorChan <- err
		return
	}
	resultChan <- result
}

func AllQuery() ([]ContactHistory, error) {
	resultChan := make(chan []ContactHistory)
	errorChan := make(chan error)
	go RunQuery(MongoQuery, resultChan, errorChan)
	go RunQuery(SOAPQuery, resultChan, errorChan)
	go RunQuery(SharepointQuery, resultChan, errorChan)
	results := make([]ContactHistory, 0)
	var err error
	for i := 1; i <= 3; i++ {
		select {
		case rs := <-resultChan:
			results = append(results, rs...)
		case e := <-errorChan:
			if err == nil {
				err = e
			} else {
				err = fmt.Errorf("%s,%s", err, e)
			}
		}
	}
	return results, err
}

func main() {
	result, err := AllQuery()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
