package app

import (
	"fmt"
	"log"
	"net/http"

	"app/datastore"
	"app/handler"
)

const projectId = "empty"

func Listen(port string) {

	datastore.ProjectId = projectId

	handler.Register()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
