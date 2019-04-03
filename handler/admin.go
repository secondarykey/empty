package handler

import (
	"fmt"
	"log"
	"net/http"

	"app/datastore"
)

func init() {
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	err := datastore.Put("test")
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, "Hello, Go112!")
}
