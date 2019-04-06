package handler_test

import (
	"app/handler"

	"fmt"
	"net/http"

	"net/http/httptest"
	"testing"
)

//Parameter

func TestNewParameter(t *testing.T) {
}

func TestParameterSetTemplate(t *testing.T) {
}

func TestParameterSet(t *testing.T) {
}

func TestParameterGet(t *testing.T) {
}

//Router

func TestRouter(t *testing.T) {

	r := handler.NewRouter(handler.Direct)

	r.Add("/", testIndex)
	r.Add("/test/", test)
	r.Add("/test/{param1}", testParam)
	r.Add("/test/add", testAdd)
	r.Add("/test/{param2}/{param3}", testParam2)

	req := httptest.NewRequest("GET", "/hello", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusNotFound != rec.Code {
		t.Errorf("/hello is not found")
	}

	req = httptest.NewRequest("GET", "/", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/ is found")
	}

	if rec.Body.String() != "Index" {
		t.Errorf("/ write Index")
	}

	req = httptest.NewRequest("GET", "/test/", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/test/ is found")
	}

	if rec.Body.String() != "Test Directory" {
		t.Errorf("/test/ write Index")
	}

	req = httptest.NewRequest("GET", "/test/aaa", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/test/aaa is found")
	}

	if rec.Body.String() != "Test Param[aaa]" {
		t.Errorf("/test/aaa write parameter[aaa]")
	}

	req = httptest.NewRequest("GET", "/test/bbb", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/test/bbb is found")
	}

	if rec.Body.String() != "Test Param[bbb]" {
		t.Errorf("/test/bbb write parameter[bbb]")
	}

	req = httptest.NewRequest("GET", "/test/add", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/test/add is found")
	}

	if rec.Body.String() != "Test add method" {
		t.Errorf("/test/add not call parameter method")
	}

	req = httptest.NewRequest("GET", "/test/aaa/bbb", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/test/aaa/bbb is found")
	}

	if rec.Body.String() != "Test Param[aaa][bbb]" {
		t.Errorf("/test/aaa/bbb double parameter")
	}
}

func testIndex(p *handler.Parameter) error {
	fmt.Fprintf(p.Res, "Index")
	return nil
}

func test(p *handler.Parameter) error {
	fmt.Fprintf(p.Res, "Test Directory")
	return nil
}

func testParam(p *handler.Parameter) error {
	fmt.Fprintf(p.Res, "Test Param[%s]", p.Get("param1"))
	return nil
}

func testAdd(p *handler.Parameter) error {
	fmt.Fprintf(p.Res, "Test add method")
	return nil
}

func testParam2(p *handler.Parameter) error {

	if p.Get("param1") != "" {
		return fmt.Errorf("Error")
	}

	fmt.Fprintf(p.Res, "Test Param[%s][%s]", p.Get("param2"), p.Get("param3"))

	return nil
}
