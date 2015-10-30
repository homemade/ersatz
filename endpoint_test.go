package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ItReturns404OnHTTPVerbsItDoesntHaveRegistered(t *testing.T) {
	resp, req := setupHttpClient(t, HTTP_GET)

	endpoint := &Endpoint{Path: "/", VerbsToDefinitions: make(map[string][]*os.File)}

	endpoint.ServeHTTP(resp, req)

	assert.Equal(t, resp.Code, 404)
}

func setupHttpClient(t *testing.T, httpVerb string) (*httptest.ResponseRecorder, *http.Request) {
	resp := httptest.NewRecorder()

	req, err := http.NewRequest(httpVerb, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	return resp, req
}
