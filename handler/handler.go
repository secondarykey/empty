package handler

import (
	"net/http"
)

func Register(template, static string) {

	registerStaticHandler(static)

	templateDirectory = template

	r := NewTemplateRouter(false)
	r.Add("/", indexHandler)
	r.Add("/login", loginHandler)
	r.Add("/callback", callbackHandler)
	http.Handle("/", r)

	ar := NewTemplateRouter(true)
	ar.Add("/admin/", adminHandler)
	http.Handle("/admin/", ar)

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
	p.SetTemplate("index.tmpl")
	return nil
}

func adminHandler(p *Parameter) error {
	p.SetTemplate("admin/index.tmpl")
	return nil
}

func loginHandler(p *Parameter) error {
	err := redirectLogin(p.Res, p.Req, "/callback")
	if err != nil {
		return err
	}
	p.Direct = true
	return nil
}

const OAuthLoginCookieName = "OAuthLoggined"

func callbackHandler(p *Parameter) error {

	err := authorization(p.Res, p.Req, "roles/owner")
	if err != nil {
		return err
	}

	sc := http.Cookie{
		Name:   OAuthLoginCookieName,
		Value:  "true",
		MaxAge: 60 * 60 * 24,
		Path:   "/admin/",
	}
	http.SetCookie(p.Res, &sc)

	p.Redirect("/admin/", http.StatusFound)
	return nil
}
