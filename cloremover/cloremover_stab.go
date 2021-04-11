package cloremover

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing/fstest"
	"time"
)

/*

test-folder (.)
├── clone1
├── clone2
├── test-folder1
│   ├── clone1
│   ├── clone2
│	└── unique
├── test-folder2
│   ├── clone1
│   ├── clone3
│	└── unique
└── test-folder3
    ├── clone1
    ├── clone2
    ├── unique
	├── unique2
 	└── test-folder4
	    ├── clone2
    	└── clone3
*/

var mapFS = fstest.MapFS{
	// root-folder
	"clone1": {
		Data:    []byte("clone1"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"clone2": {
		Data:    []byte("clone2 file"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	//test-folder1
	"test-folder1": {
		Mode:    fs.ModeDir,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder1/clone1": {
		Data:    []byte("clone1"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder1/clone2": {
		Data:    []byte("clone2 file"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder1/unique": {
		Data:    []byte("unique file from test-folder1"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	//test-folder2
	"test-folder2": {
		Mode:    fs.ModeDir,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder2/clone1": {
		Data:    []byte("clone1"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder2/clone3": {
		Data:    []byte("clone3 test file"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder2/unique": {
		Data:    []byte("one more unique file"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	//test-folder3
	"test-folder3": {
		Mode:    fs.ModeDir,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder3/clone1": {
		Data:    []byte("clone1"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder3/clone2": {
		Data:    []byte("clone2 file"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder3/unique": {
		Data:    []byte("it's an ordinary unique file"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder3/unique2": {
		Data:    []byte("one more unique file"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	//test-folder4
	"test-folder3/test-folder4": {
		Mode:    fs.ModeDir,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder3/test-folder4/clone2": {
		Data:    []byte("clone2 file"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
	"test-folder3/test-folder4/clone3": {
		Data:    []byte("clone3 test file"),
		Mode:    fs.ModeIrregular,
		ModTime: time.Now(),
		Sys:     &sysVal,
	},
}

var sysVal int

func createTestFiles() {

	os.MkdirAll("test-folder", os.ModePerm)
	os.MkdirAll(filepath.Join("test-folder", "test-folder1"), os.ModePerm)
	os.MkdirAll(filepath.Join("test-folder", "test-folder2"), os.ModePerm)
	os.MkdirAll(filepath.Join("test-folder", "test-folder3"), os.ModePerm)
	os.MkdirAll(filepath.Join("test-folder", "test-folder3", "test-folder4"), os.ModePerm)
	os.Create(filepath.Join("test-folder", "test-folder1", "clone1"))
	os.Create(filepath.Join("test-folder", "test-folder1", "clone2"))
	os.Create(filepath.Join("test-folder", "test-folder1", "unique1"))
	os.Create(filepath.Join("test-folder", "test-folder2", "clone1"))
	os.Create(filepath.Join("test-folder", "test-folder2", "clone3"))
	os.Create(filepath.Join("test-folder", "test-folder2", "unique2"))
	os.Create(filepath.Join("test-folder", "test-folder3", "clone1"))
	os.Create(filepath.Join("test-folder", "test-folder3", "clone2"))
	os.Create(filepath.Join("test-folder", "test-folder3", "unique3"))
	os.Create(filepath.Join("test-folder", "test-folder3", "unique4"))
	os.Create(filepath.Join("test-folder", "test-folder3", "test-folder4", "clone2"))
	os.Create(filepath.Join("test-folder", "test-folder3", "test-folder4", "clone3"))
	os.Create(filepath.Join("test-folder", "clone1"))
	os.Create(filepath.Join("test-folder", "clone2"))

}

// stab for test with simple config
func (c *ConfigType) useStab() {
	c.DirPath = "test-folder"
	c.RemoveFlag = false
	c.ConfirmFlag = "on"
	c.ShowFiles = 3
	c.DirLimit = 5
	c.LogFile = "not_set"
	c.LogLevel = "info"

}

// cropped fileData-type for testing only
type fileDataReduced struct {
	dir      string
	fileName string
}
