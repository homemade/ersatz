package main

import (
	"fmt"
	"os"
	"path/filepath"
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

func NewStartApp(port int, rootDir string) *StartApp {
	return &StartApp{
		RootDir:              rootDir,
		Port:                 port,
		PathToVerb:           make(map[string][]string),
		PathVerbToDefinition: make(map[string][]string),
	}
}

type StartApp struct {
	RootDir              string
	Port                 int
	PathToVerb           map[string][]string
	PathVerbToDefinition map[string][]string
}

func (s *StartApp) Run() error {
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

	return nil
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
