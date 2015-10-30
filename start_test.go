package main

import (
	"io/ioutil"
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

	err := NewStartApp(9999, dirName).Run()

	assert.EqualError(t, err, expectedErr.Error())
}

func Test_ItReturnsErrorIfNoHTTPVerbsFoundInSubDIRs(t *testing.T) {
	dirName, cleanupFn := setupRootDir(t)
	defer cleanupFn()
	setupSubFolders(t, dirName, [][]string{
		{"endpoint1"},
		{"endpoint2", "subpoint1"},
	})

	err := NewStartApp(9999, dirName).Run()

	assert.Error(t, err, ErrNoVerbsFound(dirName).Error())
}

func Test_ItBuildsAMapOfPathToHTTPVerbs(t *testing.T) {
	dirName, cleanupFn := setupRootDir(t)
	defer cleanupFn()
	setupSubFolders(t, dirName, [][]string{
		{"endpoint1", "POST"},
		{"endpoint1", "GET"},
		{"endpoint2", "GET"},
		{"endpoint2", "subpoint1", "PUT"},
		{"endpoint2", "subpoint1", "DELETE"},
	})
	expectedPathMap := map[string][]string{
		"endpoint1": {
			"GET",
			"POST",
		},
		"endpoint2": {
			"GET",
		},
		"endpoint2/subpoint1": {
			"DELETE",
			"PUT",
		},
	}

	startApp := NewStartApp(9999, dirName)
	startApp.Run()

	assert.Equal(t, expectedPathMap, startApp.PathToVerb)
}

func Test_ItErrorsIfNoDefinitionFilesFound(t *testing.T) {
	dirName, cleanupFn := setupRootDir(t)
	defer cleanupFn()

	setupSubFolders(t, dirName, [][]string{
		{"endpoint1", "POST"},
		{"endpoint1", "GET"},
		{"endpoint2", "GET"},
		{"endpoint2", "subpoint1", "PUT"},
		{"endpoint2", "subpoint1", "DELETE"},
	})

	err := NewStartApp(9999, dirName).Run()

	assert.Error(t, err, ErrNoDefinitionsFound(dirName).Error())
}

func Test_ItBuildsMapOfPathVerbsAndDefinitionFiles(t *testing.T) {
	dirName, cleanupFn := setupRootDir(t)
	defer cleanupFn()

	setupSubFolders(t, dirName, [][]string{
		{"endpoint1", "POST"},
		{"endpoint1", "GET"},
		{"endpoint2", "GET"},
		{"endpoint2", "subpoint1", "PUT"},
		{"endpoint2", "subpoint1", "DELETE"},
	})

	_, err := os.Create(filepath.Join(dirName, "endpoint1", "POST", "default.json"))
	assert.Nil(t, err)
	_, err = os.Create(filepath.Join(dirName, "endpoint1", "POST", "variation-1.json"))
	assert.Nil(t, err)
	_, err = os.Create(filepath.Join(dirName, "endpoint1", "POST", "variation-2.json"))
	assert.Nil(t, err)
	_, err = os.Create(filepath.Join(dirName, "endpoint1", "GET", "default.json"))
	assert.Nil(t, err)
	_, err = os.Create(filepath.Join(dirName, "endpoint2", "subpoint1", "PUT", "default.json"))
	assert.Nil(t, err)
	_, err = os.Create(filepath.Join(dirName, "endpoint2", "subpoint1", "DELETE", "default.json"))
	assert.Nil(t, err)
	_, err = os.Create(filepath.Join(dirName, "endpoint2", "subpoint1", "DELETE", "variation-1.json"))
	assert.Nil(t, err)

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

	startApp := NewStartApp(9999, dirName)
	err = startApp.Run()
	assert.Nil(t, err)

	assert.Equal(t, expectedMap, startApp.PathVerbToDefinition)
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
