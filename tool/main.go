package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var host = "localhost:8080"

func main() {
	// grab the start time so we can time our operation
	start := time.Now()

	// open the file, log the error and panic if we have one
	input, err := os.Open("./cats.json")
	if err != nil {
		log.Fatal(err)
	}

	// do the work!
	avg := calcAverageWeight(input)

	log.Print("the average weight is ", avg)
	log.Print("work complete in ", time.Since(start))
}

// calcAverageWeight read from an io.ReaderCloser containing Cats and make a request for each line
func calcAverageWeight(file io.ReadCloser) float32 {
	// set up our func to close the file as soon as it completes (just in case)
	defer file.Close()

	// create a buffered channel to move data from the file to our goroutines
	cats := make(chan []byte, 10000)
	// create a buffered channel to gather the results
	weights := make(chan int64, 1000)
	// create a var to hold our average
	var avg float32

	// create a WaitGroup to help us keep track of our goroutines
	var weighers sync.WaitGroup
	for i := 0; i < 10; i++ {
		// increment the waitgroup to register our new goroutine
		weighers.Add(1)
		// kick off a bounce checker in a goroutine!
		go catWeigher(cats, weights, &weighers)
	}

	// create a goroutine to gather the weights for an average
	// use a simple chan to mark complete  bc it's only 1 goroutine
	done := make(chan bool)
	go func() {
		var catCount int64
		var totalWeight int64
		for weight := range weights {
			catCount++
			totalWeight += weight
		}

		avg = float32(totalWeight) / float32(catCount)

		// alert the main thread we're done
		done <- true
	}()

	// scan each line out of our file and toss it on a buffered channel
	go func() {
		// create a new scanner to help us iterate over each line
		s := bufio.NewScanner(file)
		for s.Scan() {
			// allocate bytes into fresh slice to preserve order
			catBytes := append([]byte{}, s.Bytes()...)
			cats <- catBytes
		}

		// close the channel to alert our workers there is no more to read
		close(cats)

		// check for any errors in the scanner after iterating
		if err := s.Err(); err != nil {
			log.Print("scanner error: ", err)
		}

		// close our input file at the end
		file.Close()
	}()

	// wait for the goroutine pool to finish their work and return
	weighers.Wait()
	// tell the output reader that there are no more records to register
	close(weights)
	// wait for result reader to complete.
	<-done
	return avg
}

func catWeigher(cats chan []byte, weights chan int64, wg *sync.WaitGroup) {
	// set the func up to mark itself as done as soon as it completes
	defer wg.Done()

	// iterate over each line passed along through the channel
	for catBytes := range cats {
		// get the cats weight. log and continue on error
		weight, err := getWeight(host, catBytes)
		if err != nil {
			log.Print(err)
			continue
		}
		// pass the weight along to the output channel
		weights <- weight
	}
}

func getWeight(host string, catBytes []byte) (int64, error) {
	// post against our isbounce service with the line as the body
	resp, err := http.Post(fmt.Sprintf("http://%s/catweight", host), "application/json", bytes.NewReader(catBytes))
	if err != nil {
		return 0, err
	}

	// create a var and read from the response body into it
	var weight []byte
	weight, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()

	// parse the response into an int
	return strconv.ParseInt(string(weight), 10, 64)
}
