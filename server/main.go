package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Cat struct {
	Name  string `json:"name"`
	Breed string `json:"breed"`
}

func main() {

	// catweight will parse the given cat JSON and return the cat's weight (length of name + breed)
	http.HandleFunc("/catweight", func(w http.ResponseWriter, r *http.Request) {
		// create a variable to hold our cat
		var cat Cat
		// create a json decoder from our request body and decode into the cat var
		err := json.NewDecoder(r.Body).Decode(&cat)

		// close the request body because it's a ReadCloser
		r.Body.Close()

		// check for an error and respond appropriately
		if err != nil {
			http.Error(w, fmt.Sprint("bad request: ", err), 400)
			return
		}

		// log our cat 
		log.Printf("We've got a cat! %+v", cat)

		// cats weight == len(name + breed)
		weight := len(cat.Name + cat.Breed)

		// write our success response
		fmt.Fprint(w, weight)
	})

	http.Handle("/", http.FileServer(http.Dir("./")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
