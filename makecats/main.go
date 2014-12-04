package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
)

var (
	names = []string{
		"Ziggy",
		"Bowie",
		"Felix",
		"Garfield",
		"Heathcliff",
	}

	breeds = []string{
		"Tabby",
		"Maine Coon",
		"Russian Blue",
		"British Shorthair",
		"Siamese",
		"Persian",
		"American Shorthair",
		"Burmese",
	}
)

func main() {
	out, err := os.Create("./cats.json")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	e := json.NewEncoder(out)
	for i := 0; i < 10000; i++ {
		cat := struct {
			Name  string `json:"name"`
			Breed string `json:"breed"`
		}{
			Name:  names[rand.Intn(len(names))],
			Breed: breeds[rand.Intn(len(breeds))],
		}
		if err = e.Encode(cat); err != nil {
			log.Fatal(err)
		}
	}
}
