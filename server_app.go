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
		RootDir:                   rootDir,
		Port:                      port,
		EndpointCache:             make(EndpointCache),
		EndpointVariationSchedule: make(EndpointVariationSchedule),
		Mux: http.NewServeMux(),
	}
}

type ServerApp struct {
	RootDir                   string
	Port                      string
	Mux                       *http.ServeMux
	EndpointCache             EndpointCache
	EndpointVariationSchedule EndpointVariationSchedule
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

	s.Mux.HandleFunc("/__ersatz", s.HandleControlRequest)
	s.Mux.HandleFunc("/", s.HandleMockRequest)

	return nil
}

/////////////////////////////////////////////////////
// Given an exit channel, run everything
/////////////////////////////////////////////////////

func (s *ServerApp) Run(exit chan interface{}) {
	go manners.ListenAndServe(fmt.Sprintf(":%s", s.Port), s.Mux)

	<-exit

	manners.Close()
}

/////////////////////////////////////////////////////
// Handle control requests to the ersatz process
/////////////////////////////////////////////////////

func (s *ServerApp) HandleControlRequest(w http.ResponseWriter, r *http.Request) {
}

/////////////////////////////////////////////////////
// This is the main handler for mock requests
/////////////////////////////////////////////////////

func (s *ServerApp) HandleMockRequest(w http.ResponseWriter, r *http.Request) {

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

	if ep.ResponseCode > 0 {
		w.WriteHeader(ep.ResponseCode)
	}

	w.Write(body)
}

/////////////////////////////////////////////////////
// Fetch a given endpoint by method and url
/////////////////////////////////////////////////////

func (s *ServerApp) fetchEndpoint(url, method string) (*Endpoint, error) {

	variant, err := s.fetchVariation(url, method)

	if err != nil {
		return nil, err
	}

	// Check to see if the request is in the cache
	if c, exists := s.EndpointCache[VariableEndpointIndex{EndpointIndex{url, method}, variant}]; exists {
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

/////////////////////////////////////////////////////
// Find the active variation for an endpoint index
/////////////////////////////////////////////////////

func (s *ServerApp) fetchVariation(url, method string) (string, error) {
	return "default", nil
}
