package app

import (
	"fmt"
	"net/http"

	"app/datastore"
	"app/handler"
)

type Env struct {
	Project  string
	Port     string
	Static   string
	Template string
}

func Listen(env Env) error {
	datastore.ProjectId = env.Project
	handler.Register(env.Template, env.Static)
	return http.ListenAndServe(fmt.Sprintf(":%s", env.Port), nil)
}
