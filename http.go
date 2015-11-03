package main

import (
	"fmt"
	"net/http"
)

const (
	HTTP_DELETE  = "DELETE"
	HTTP_GET     = "GET"
	HTTP_HEAD    = "HEAD"
	HTTP_POST    = "POST"
	HTTP_PUT     = "PUT"
	HTTP_OPTIONS = "OPTIONS"
)

var HTTPVerbs = []string{
	HTTP_DELETE,
	HTTP_GET,
	HTTP_HEAD,
	HTTP_POST,
	HTTP_PUT,
	HTTP_OPTIONS,
}

type HTTPMuxer interface {
	Handle(string, http.Handler)
	Handler(r *http.Request) (http.Handler, string)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type ErrNoVerbsFound string

func (e ErrNoVerbsFound) Error() string {
	return fmt.Sprintf("No HTTP verb folders found within: %s", string(e))
}

type ErrNoDefinitionsFound string

func (e ErrNoDefinitionsFound) Error() string {
	return fmt.Sprintf("No definition files found within: %s", string(e))
}
