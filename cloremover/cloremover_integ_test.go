// +build integration

package cloremover

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFindClonesRealFS(t *testing.T) {

	expectedSlice := []fileDataReduced{
		{dir: ".", fileName: "clone1"},
		{dir: ".", fileName: "clone2"},
		{dir: "test-folder1", fileName: "clone1"},
		{dir: "test-folder1", fileName: "clone2"},
		{dir: "test-folder2", fileName: "clone1"},
		{dir: "test-folder2", fileName: "clone3"},
		{dir: "test-folder3", fileName: "clone1"},
		{dir: "test-folder3", fileName: "clone2"},
		{dir: "test-folder3/test-folder4", fileName: "clone2"},
		{dir: "test-folder3/test-folder4", fileName: "clone3"},
	}

	conf := &ConfigType{}
	conf.useStab()

	logger := logrus.New()
	logger.Out = ioutil.Discard

	createTestFiles()
	fileSystem := os.DirFS("test-folder")
	fileSlice, _ := FindClones(conf, logger, fileSystem)

	outFileSlice := make([]fileDataReduced, len(fileSlice))
	for i := range fileSlice {
		outFileSlice[i].dir = fileSlice[i].dir
		outFileSlice[i].fileName = fileSlice[i].fileName
	}

	assert.ElementsMatch(t, outFileSlice, expectedSlice, "File data slice is not valid")

}

func TestEnumDirsRealFS(t *testing.T) {

	conf := &ConfigType{}
	conf.useStab()

	createTestFiles()
	fileSystem := os.DirFS("test-folder")
	stringSlice, _ := enumDirs(conf, fileSystem)

	//	stringSlice, _ := enumDirs(conf, mapFS)

	expectedStringSlice := []string{".", "test-folder1", "test-folder2", "test-folder3", "test-folder3/test-folder4"}
	assert.ElementsMatch(t, stringSlice, expectedStringSlice, "slices not equal")
}
