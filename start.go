package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/braintree/manners"
)

func NewServerApp(port string, rootDir string) *ServerApp {
	return &ServerApp{
		RootDir: rootDir,
		Port:    port,
	}
}

type ServerApp struct {
	RootDir       string
	Port          string
	EndpointCache EndpointCache
}

/////////////////////////////////////////////////////
// Check the basic settings and register handlers
/////////////////////////////////////////////////////

func (s *ServerApp) Setup() error {

	// Check that the path exists
	path, err := os.Stat(s.RootDir)

	if os.IsNotExist(err) {
		return fmt.Errorf("no such file or directory: %s", s.RootDir)
	}

	if !path.IsDir() {
		return fmt.Errorf("Path is not a directory %s", s.RootDir)
	}

	http.HandleFunc("/", s.Handle)

	return nil
}

/////////////////////////////////////////////////////
// Given an exit channel, run everything
/////////////////////////////////////////////////////

func (s *ServerApp) Run(exit chan interface{}) {
	go manners.ListenAndServe(fmt.Sprintf(":%s", s.Port), nil)

	<-exit
}

/////////////////////////////////////////////////////
// This is the main handler
/////////////////////////////////////////////////////

func (s *ServerApp) Handle(w http.ResponseWriter, r *http.Request) {

	ep, err := s.fetchEndpoint(r.URL.Path[1:], r.Method)

	if err != nil {
		w.Header().Set("Ersatz-Error", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := json.Marshal(ep.Body)

	if err != nil {
		w.Header().Set("Ersatz-Error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))

	for k, v := range ep.Headers {
		w.Header().Set(k, v)
	}

	w.Write(body)
}

/////////////////////////////////////////////////////
// Fetch a given endpoint by method and url
/////////////////////////////////////////////////////

func (s *ServerApp) fetchEndpoint(url, method string) (*Endpoint, error) {

	variant := "default"

	// Check to see if the request is in the cache
	if c, exists := s.EndpointCache[EndpointIndex{url, method, variant}]; exists {
		return c, nil
	}

	file, e := ioutil.ReadFile(filepath.Join(s.RootDir, url, method, variant+".json"))

	if e != nil {
		return nil, e
	}

	ep := NewEndpoint()

	if e := json.Unmarshal(file, ep); e != nil {
		return nil, e
	}

	return ep, nil
}
