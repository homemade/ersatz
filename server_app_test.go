package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultJSON = `{ "response_code": 200, "headers": { "header-1": "some value" }, "body": { "a":1, "b":2, "c":3 }}`

type definitionFile struct {
	verb      string
	endpoints []string
	variants  []string
	json      string
}

func Test_ItReturnsErrorIfPathDoesntExist(t *testing.T) {

	if err := NewServerApp("9999", "/this/doesnt/exist").Setup(); err == nil {
		t.Error("Able to read from a non-existent directory")
	}
}

func Test_ItSetsupAHandlerForEachEndpointOnTheMux(t *testing.T) {
	dirName, cleanupFn := setupRootDir(t)
	defer cleanupFn()

	endpoints := []definitionFile{
		{HTTP_POST, []string{"endpoint1"}, []string{"default.json"}, defaultJSON},
		{HTTP_GET, []string{"endpoint1"}, []string{"default.json"}, defaultJSON},
		{HTTP_GET, []string{"endpoint1", "subendpoint1"}, []string{"default.json"}, defaultJSON},
	}

	setupSubFolders(t, dirName, endpoints)
	setupDefinitionFiles(t, dirName, endpoints)

	startApp := NewServerApp("9998", dirName)

	err := startApp.Setup()
	assert.Nil(t, err)

	exit := make(chan interface{})

	go startApp.Run(exit)

	for _, df := range endpoints {

		path := strings.Join(df.endpoints, "/")

		req, err := http.NewRequest(
			df.verb,
			fmt.Sprintf("http://localhost:%s/%s", "9998", path),
			bytes.NewBuffer([]byte("")),
		)

		assert.Nil(t, err)

		client := &http.Client{}
		res, err := client.Do(req)
		assert.Nil(t, err)

		if g, e := res.StatusCode, 200; g != e {
			t.Errorf("Expected status code %d, got %d", e, g)
		}
	}

	exit <- true
}

func setupDefinitionFiles(t *testing.T, rootDir string, dfs []definitionFile) {
	for _, df := range dfs {

		for _, variant := range df.variants {

			path := filepath.Join(rootDir, filepath.Join(df.endpoints...), df.verb) + "/" + variant

			// Create the file itself
			f, err := os.Create(path)
			assert.Nil(t, err)

			// Write some stock data into it
			n, err := f.Write([]byte(df.json))
			assert.Nil(t, err)
			assert.Equal(t, n, len(df.json))
		}
	}
}

func setupSubFolders(t *testing.T, rootDir string, dfs []definitionFile) {
	for _, df := range dfs {

		path := filepath.Join(rootDir, filepath.Join(df.endpoints...), df.verb)
		os.MkdirAll(path, 0755)
	}
}

func setupRootDir(t *testing.T) (string, func()) {
	dirName, err := ioutil.TempDir("", "ersatz")
	assert.Nil(t, err)

	return dirName, func() {
		if err := os.RemoveAll(dirName); err != nil {
			panic(err.Error())
		}
	}
}
