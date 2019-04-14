package app

import (
	"fmt"
	"net/http"

	"app/handler"
)

type Env struct {
	Port     string
	Static   string
	Template string
}

func Listen(env Env) error {
	handler.Register(env.Template, env.Static)
	return http.ListenAndServe(fmt.Sprintf(":%s", env.Port), nil)
}
