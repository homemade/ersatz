package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ItReturnsErrorIfCantReadFromRootDIR(t *testing.T) {
	dirName, cleanupFn := setupRootDir(t)
	defer cleanupFn()
	os.Chmod(dirName, 0333)
	expectedErr := &os.PathError{Op: "open", Path: dirName, Err: os.ErrPermission}

	err := NewStartApp("9999", dirName).Setup()

	assert.EqualError(t, err, expectedErr.Error())
}

func Test_ItReturnsErrorIfNoHTTPVerbsFoundInSubDIRs(t *testing.T) {
	dirName, cleanupFn := setupRootDir(t)
	defer cleanupFn()
	setupSubFolders(t, dirName, [][]string{
		{"endpoint1"},
		{"endpoint2", "subpoint1"},
	})

	err := NewStartApp("9999", dirName).Setup()

	assert.EqualError(t, err, ErrNoVerbsFound(dirName).Error())
}

func Test_ItBuildsAMapOfPathToHTTPVerbs(t *testing.T) {
	dirName, cleanupFn := setupRootDir(t)
	defer cleanupFn()
	setupSubFolders(t, dirName, [][]string{
		{"endpoint1", HTTP_POST},
		{"endpoint1", HTTP_GET},
		{"endpoint2", HTTP_GET},
		{"endpoint2", "subpoint1", HTTP_PUT},
		{"endpoint2", "subpoint1", HTTP_DELETE},
	})
	expectedPathMap := map[string][]string{
		"endpoint1": {
			HTTP_GET,
			HTTP_POST,
		},
		"endpoint2": {
			HTTP_GET,
		},
		"endpoint2/subpoint1": {
			HTTP_DELETE,
			HTTP_PUT,
		},
	}

	startApp := NewStartApp("9999", dirName)
	startApp.Setup()

	assert.Equal(t, expectedPathMap, startApp.PathToVerb)
}

func Test_ItErrorsIfNoDefinitionFilesFound(t *testing.T) {
	dirName, cleanupFn := setupRootDir(t)
	defer cleanupFn()

	setupSubFolders(t, dirName, [][]string{
		{"endpoint1", HTTP_POST},
		{"endpoint1", HTTP_GET},
		{"endpoint2", HTTP_GET},
		{"endpoint2", "subpoint1", HTTP_PUT},
		{"endpoint2", "subpoint1", HTTP_DELETE},
	})

	err := NewStartApp("9999", dirName).Setup()

	assert.EqualError(t, err, ErrNoDefinitionsFound(dirName).Error())
}

func Test_ItBuildsMapOfPathVerbsAndDefinitionFiles(t *testing.T) {
	dirName, cleanupFn := setupRootDir(t)
	defer cleanupFn()

	setupSubFolders(t, dirName, [][]string{
		{"endpoint1", HTTP_POST},
		{"endpoint1", HTTP_GET},
		{"endpoint2", HTTP_GET},
		{"endpoint2", "subpoint1", HTTP_PUT},
		{"endpoint2", "subpoint1", HTTP_DELETE},
	})

	setupDefinitionFiles(t, dirName, [][]string{
		{"endpoint1", HTTP_POST, "default.json"},
		{"endpoint1", HTTP_POST, "variation-1.json"},
		{"endpoint1", HTTP_POST, "variation-2.json"},
		{"endpoint1", HTTP_GET, "default.json"},
		{"endpoint2", "subpoint1", HTTP_PUT, "default.json"},
		{"endpoint2", "subpoint1", HTTP_DELETE, "default.json"},
		{"endpoint2", "subpoint1", HTTP_DELETE, "variation-1.json"},
	})

	expectedMap := map[string][]string{
		"endpoint1/POST": {
			"default.json",
			"variation-1.json",
			"variation-2.json",
		},
		"endpoint1/GET": {
			"default.json",
		},
		"endpoint2/subpoint1/PUT": {
			"default.json",
		},
		"endpoint2/subpoint1/DELETE": {
			"default.json",
			"variation-1.json",
		},
	}

	startApp := NewStartApp("9999", dirName)
	startApp.HTTPMuxer = http.NewServeMux()
	err := startApp.Setup()
	assert.Nil(t, err)

	assert.Equal(t, expectedMap, startApp.PathVerbToDefinition)
}

func Test_ItSetupsAHandlerForEachEndpointOnTheMux(t *testing.T) {
	dirName, cleanupFn := setupRootDir(t)
	defer cleanupFn()

	setupSubFolders(t, dirName, [][]string{
		{"endpoint1", HTTP_POST},
		{"endpoint1", HTTP_GET},
		{"endpoint2", HTTP_OPTIONS},
		{"endpoint2", HTTP_HEAD},
		{"endpoint2", "subpoint1", HTTP_PUT},
		{"endpoint2", "subpoint1", HTTP_DELETE},
	})

	setupDefinitionFiles(t, dirName, [][]string{
		{"endpoint1", HTTP_POST, "default.json"},
		{"endpoint1", HTTP_GET, "default.json"},
		{"endpoint2", "subpoint1", HTTP_PUT, "default.json"},
		{"endpoint2", "subpoint1", HTTP_DELETE, "default.json"},
	})

	startApp := NewStartApp("9999", dirName)
	startApp.HTTPMuxer = http.NewServeMux()

	err := startApp.Setup()
	assert.Nil(t, err)

	for endpoint, verbs := range startApp.PathToVerb {
		for _, verb := range verbs {
			req, err := http.NewRequest(verb, "/"+endpoint, nil)
			assert.Nil(t, err)

			_, path := startApp.HTTPMuxer.Handler(req)

			assert.Equal(t, "/"+endpoint, path)
		}
	}
}

func setupDefinitionFiles(t *testing.T, rootDir string, filePathParts [][]string) {
	for _, pathParts := range filePathParts {
		_, err := os.Create(filepath.Join(append([]string{rootDir}, pathParts...)...))
		assert.Nil(t, err)
	}
}

func setupSubFolders(t *testing.T, rootDir string, subDirs [][]string) {
	for _, subFilePath := range subDirs {
		os.MkdirAll(filepath.Join(append([]string{rootDir}, subFilePath...)...), 0755)
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
