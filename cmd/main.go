package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"app"
)

const static = "public"
const templates = "templates"
const mainPlace = "cmd"

//
// gcloud beta emulators datastore start --project=my-project-id
// $(gcloud beta emulators datastore env-init)
//
// dev_appserver.py -A=my-project-id .
//
func main() {

	env := app.Env{}
	env.Port = os.Getenv("PORT")
	if env.Port == "" {
		env.Port = "8080"
	}

	abs, _ := filepath.Abs("")

	if strings.LastIndex(abs, mainPlace) != len(abs)-len(mainPlace) {
		abs = filepath.Join(abs, mainPlace)
	}

	env.Template = filepath.Join(abs, templates)
	env.Static = filepath.Join(abs, static)

	log.Fatal(app.Listen(env))
}
