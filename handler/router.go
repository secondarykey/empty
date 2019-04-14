package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

type Router struct {
	behavior Behavior
	login    bool
	pattern  map[string]Handler
	keyMap   map[string][]string
}

type Behavior int

const (
	Direct Behavior = iota
	Template
	JSON
)

func NewRouter(b Behavior, login bool) Router {
	r := Router{}
	r.behavior = b
	r.login = login
	return r
}

func NewTemplateRouter(login bool) Router {
	r := Router{}
	r.login = login
	r.behavior = Template
	return r
}

func NewJSONRouter(login bool) Router {
	r := Router{}
	r.login = login
	r.behavior = JSON
	return r
}

type Handler func(*Parameter) error

func (router *Router) Add(key string, h Handler) error {
	if router.pattern == nil {
		router.pattern = make(map[string]Handler)
		router.keyMap = make(map[string][]string)
	}
	router.pattern[key] = h

	slc := strings.Split(key, "/")
	for _, elm := range slc {
		if strings.Index(elm, "{") == 0 &&
			//TODO Error
			strings.Index(elm, "}") == len(elm)-1 {
			router.keyMap[key] = slc
			break
		}
	}
	return nil
}

func (router *Router) getHandler(p *Parameter) (Handler, error) {

	path := p.Req.URL.Path
	var h Handler
	for key, elm := range router.pattern {
		if key == path {
			return elm, nil
		}
	}

	slc := strings.Split(path, "/")
	for key, elm := range router.keyMap {
		if len(elm) != len(slc) {
			continue
		}

		flag := true
		wkMap := make(map[string]string)

		for idx, buf := range elm {
			if strings.Index(buf, "{") == 0 &&
				strings.Index(buf, "}") == len(buf)-1 {
				wkKey := buf[1 : len(buf)-1]
				wkMap[wkKey] = slc[idx]
			} else {
				if buf != slc[idx] {
					flag = false
					break
				}
			}
		}

		if flag {
			if len(wkMap) != 0 {
				p.values = make(map[string]string)
				for key, val := range wkMap {
					p.values[key] = val
				}
			}
			return router.pattern[key], nil
		}
	}

	return h, xerrors.Errorf("URL Pattern Not Found[%s]", path)
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	p := NewParameter(w, r)

	if router.login {
		sc, err := r.Cookie(OAuthLoginCookieName)
		if err != nil {
			p.Redirect("/login", http.StatusMovedPermanently)
			return
		} else {
			if sc.Value != "true" {
				setError(p, http.StatusUnauthorized, xerrors.Errorf("Login cookie error"))
				return
			}
		}
	}

	if router.behavior == Direct {
		p.Direct = true
	}

	handler, err := router.getHandler(p)
	if err != nil {
		setError(p, http.StatusNotFound, err)
		return
	}

	err = handler(p)
	if err != nil {
		setError(p, http.StatusInternalServerError, err)
		return
	}

	if !p.Direct {
		if router.behavior == Template {
			if p.templates == nil || len(p.templates) == 0 {
				err = xerrors.Errorf("error Behavior->Template but templates not setting")
			} else {
				err = setTemplates(w, p.output, p.templates...)
			}
		} else if router.behavior == JSON {
			err = setJSON(w, p.output)
		} else {
			err = xerrors.Errorf("error Behavior")
		}
	}

	if err != nil {
		setError(p, http.StatusInternalServerError, err)
		return
	}
}

type Parameter struct {
	Res http.ResponseWriter
	Req *http.Request

	values    map[string]string
	output    map[string]interface{}
	templates []string

	//Headers
	Direct bool
}

func NewParameter(w http.ResponseWriter, r *http.Request) *Parameter {
	p := Parameter{}
	p.Req = r
	p.Res = w
	p.Direct = false
	p.values = nil
	p.output = nil
	p.templates = nil
	return &p
}

func (p *Parameter) SetTemplate(tmpls ...string) {
	p.templates = tmpls
}

func (p *Parameter) Redirect(path string, status int) {
	p.Direct = true
	http.Redirect(p.Res, p.Req, path, status)
}

func (p *Parameter) Set(key string, v interface{}) {
	if p.output == nil {
		p.output = make(map[string]interface{})
	}
	p.output[key] = v
}

func (p *Parameter) Get(key string) string {
	if p.values == nil {
		return ""
	}
	return p.values[key]
}

func setJSON(w http.ResponseWriter, v interface{}) error {
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	return err
}

var templateDirectory string

func setTemplates(w io.Writer, param interface{}, files ...string) error {

	tmpls := make([]string, len(files))
	for idx, file := range files {
		tmpls[idx] = filepath.Join(templateDirectory, file)
	}

	tmpl := template.Must(template.ParseFiles(tmpls...))
	return tmpl.Execute(w, param)
}

func setError(p *Parameter, status int, err error) {

	buf := fmt.Sprintf("%+v", err)
	log.Printf("%s\n", buf)

	p.Res.WriteHeader(status)

	dto := struct {
		Status      int
		Title       string
		Description string
	}{status, http.StatusText(status), buf}

	p.Set("Error", dto)
	setTemplates(p.Res, p.output, "error.tmpl")
}
