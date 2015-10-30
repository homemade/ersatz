package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/braintree/manners"
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

type ErrNoVerbsFound string

func (e ErrNoVerbsFound) Error() string {
	return fmt.Sprintf("No HTTP verb folders found within: %s", string(e))
}

type ErrNoDefinitionsFound string

func (e ErrNoDefinitionsFound) Error() string {
	return fmt.Sprintf("No definition files found within: %s", string(e))
}

type HTTPMuxer interface {
	Handle(string, http.Handler)
	Handler(r *http.Request) (http.Handler, string)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func NewStartApp(port string, rootDir string) *StartApp {
	return &StartApp{
		RootDir:              rootDir,
		Port:                 port,
		PathToVerb:           make(map[string][]string),
		PathVerbToDefinition: make(map[string][]string),
	}
}

type StartApp struct {
	RootDir              string
	Port                 string
	PathToVerb           map[string][]string
	PathVerbToDefinition map[string][]string
	HTTPServer           *manners.GracefulServer
	HTTPMuxer
}

func (s *StartApp) Setup() error {
	err := filepath.Walk(s.RootDir, s.walkFn)
	if err != nil {
		return err
	}

	if len(s.PathToVerb) == 0 {
		return ErrNoVerbsFound(s.RootDir)
	}

	if len(s.PathVerbToDefinition) == 0 {
		return ErrNoDefinitionsFound(s.RootDir)
	}

	for endpoint := range s.PathToVerb {
		s.HTTPMuxer.Handle("/"+endpoint, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	}

	s.HTTPServer = manners.NewWithServer(&http.Server{Handler: s.HTTPMuxer, Addr: ":" + s.Port})

	return nil
}

func (s *StartApp) Run() error {
	return s.HTTPServer.ListenAndServe()
}

func (s *StartApp) Stop() bool {
	if s.HTTPServer != nil {
		return s.HTTPServer.BlockingClose()
	}

	return true
}

func (s *StartApp) walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	relativePath, err := filepath.Rel(s.RootDir, path)
	if err != nil {
		return err
	}
	relativeDirPath := filepath.Dir(relativePath)

	if info.IsDir() {
		for _, httpVerb := range HTTPVerbs {
			if info.Name() == httpVerb {
				s.PathToVerb[relativeDirPath] = append(s.PathToVerb[relativeDirPath], httpVerb)
			}
		}
	} else {
		s.PathVerbToDefinition[relativeDirPath] = append(s.PathVerbToDefinition[relativeDirPath], info.Name())
	}

	return nil
}
