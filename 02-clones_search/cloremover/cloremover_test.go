package cloremover

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCloremover(t *testing.T) {

	os.MkdirAll("test-folder", os.ModePerm)
	os.MkdirAll(filepath.Join("test-folder", "test-folder1"), os.ModePerm)
	os.MkdirAll(filepath.Join("test-folder", "test-folder2"), os.ModePerm)
	os.MkdirAll(filepath.Join("test-folder", "test-folder3"), os.ModePerm)
	os.Create(filepath.Join("test-folder", "test-folder1", "clone1"))
	os.Create(filepath.Join("test-folder", "test-folder1", "clone2"))
	os.Create(filepath.Join("test-folder", "test-folder1", "unique1"))
	os.Create(filepath.Join("test-folder", "test-folder2", "clone1"))
	os.Create(filepath.Join("test-folder", "test-folder2", "clone1"))
	os.Create(filepath.Join("test-folder", "test-folder2", "unique2"))
	os.Create(filepath.Join("test-folder", "test-folder3", "clone1"))
	os.Create(filepath.Join("test-folder", "test-folder3", "clone2"))
	os.Create(filepath.Join("test-folder", "test-folder3", "unique3"))
	os.Create(filepath.Join("test-folder", "test-folder3", "unique4"))
	os.Create(filepath.Join("test-folder", "clone1"))
	os.Create(filepath.Join("test-folder", "clone2"))

	workingDir, _ := os.Getwd()
	expectedSlice := []fileData{
		{dir: filepath.Join(workingDir, "test-folder"), fileName: "clone1", sizeInBytes: 0, id: 1},
		{dir: filepath.Join(workingDir, "test-folder"), fileName: "clone2", sizeInBytes: 0, id: 2},
		{dir: filepath.Join(workingDir, "test-folder", "test-folder1"), fileName: "clone1", sizeInBytes: 0, id: 3},
		{dir: filepath.Join(workingDir, "test-folder", "test-folder1"), fileName: "clone2", sizeInBytes: 0, id: 4},
		{dir: filepath.Join(workingDir, "test-folder", "test-folder2"), fileName: "clone1", sizeInBytes: 0, id: 5},
		{dir: filepath.Join(workingDir, "test-folder", "test-folder3"), fileName: "clone1", sizeInBytes: 0, id: 6},
		{dir: filepath.Join(workingDir, "test-folder", "test-folder3"), fileName: "clone2", sizeInBytes: 0, id: 7},
	}

	conf := &ConfigType{}
	_ = ReadFlags(conf)
	print(conf)
	fileSlice, _ := FindClones(conf, nil)

	var counter uint8
	for _, wantData := range expectedSlice {
		for _, gotData := range fileSlice {
			if compareStructs(wantData, gotData) {
				counter += 1
				break
			}
		}
	}

	if counter != uint8(len(expectedSlice)) {
		t.Errorf("test FAILED !!! Expected %d identical slice elements, got %d", len(expectedSlice), counter)
	}
}

func compareStructs(want, got fileData) bool {
	return want.dir == got.dir && want.fileName == got.fileName
}
