package handler

import (
	//"app/datastore"
	//"fmt"
	"net/http"
)

func Register(template, static string) {

	registerStaticHandler(static)

	templateDirectory = template

	r := NewRouter(Template)
	r.Add("/admin/", adminHandler)
	r.Add("/", indexHandler)

	http.Handle("/", r)

	//apiRouter := NewRouter(JSON)
	//apiRouter.Add("/api/", apiHandler)
	//http.Handle("/api/", apiRouter)
}

func registerStaticHandler(static string) {

	pub := http.Dir(static)
	s := http.FileServer(pub)

	http.Handle("/js/", s)
	http.Handle("/css/", s)
	http.Handle("/images/", s)
	http.Handle("/favicon.ico", s)
}

func indexHandler(p *Parameter) error {

	//err := datastore.Put("test")
	//if err != nil {
	//return err
	//}

	//fmt.Fprint(p.Res, "Hello, Go112!")
	//p.Direct = true
	p.SetTemplate("index.tmpl")

	return nil
}

func adminHandler(p *Parameter) error {

	//err := datastore.Put("test")
	//if err != nil {
	//return err
	//}

	//fmt.Fprint(p.Res, "Hello, Go112!")
	//p.Direct = true
	p.SetTemplate("admin/index.tmpl")

	return nil
}
