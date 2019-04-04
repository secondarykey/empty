package handler

import (
	"net/http"
	//"app/datastore"
)

const PublicDirectory = "public"

func Register() {

	registerStaticHandler()

	publicRouter := NewRouter(Template)
	publicRouter.Add("/", indexHandler)
	http.Handle("/", publicRouter)

	//apiRouter := NewRouter(JSON)
	//apiRouter.Add("/api/", apiHandler)
	//http.Handle("/api/", apiRouter)
}

func registerStaticHandler() {

	pub := http.Dir(PublicDirectory)
	s := http.FileServer(pub)

	http.Handle("/js/", s)
	http.Handle("/css/", s)
	http.Handle("/images/", s)
	http.Handle("/favicon.ico", s)
}

func indexHandler(p *Parameter) error {

	if p.Req.URL.Path != "/" {
		return fmt.Errorf("Bad Request")
	}

	//err := datastore.Put("test")
	//if err != nil {
	//return err
	//}

	//fmt.Fprint(p.Res, "Hello, Go112!")
	//p.Direct = true
	p.SetTemplate("admin/index.tmpl")

	return nil
}
