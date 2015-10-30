package main

import (
	"net/http"
	"os"
)

type Endpoint struct {
	Path               string
	VerbsToDefinitions map[string][]*os.File
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for verb := range e.VerbsToDefinitions {
		if r.Method == verb {

		}
	}

	http.NotFound(w, r)
}
