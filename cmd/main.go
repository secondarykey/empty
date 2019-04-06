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

func main() {

	env := app.Env{
		Project: "empty",
	}

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
