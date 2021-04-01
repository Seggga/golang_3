package cloremover

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

// Clones looks for clone-files in the given directory and it's subdirectories
func FindClones(conf *ConfigType) ([]fileData, error) {

	//enumerate all subdirectories in the given dirPath
	dirSlice, err := enumDirs(conf)
	if err != nil {
		return nil, err
	}

	var (
		// pool of workers for each subdirectory
		pool = make(chan struct{}, len(dirSlice))
		// ch transfers data from workers to dataCollector
		ch = make(chan fileData, 1)
		// slice to hold files's data
		fileSlice []fileData
	)

	// start pool manager
	go manager(pool, ch)
	// start ch-writers - one goroutine for each subdirectory
	for _, someDir := range dirSlice {
		go func(someDir string) {
			enumFiles(someDir, ch)
			pool <- struct{}{}
		}(someDir)
	}

	// start ch-reader
	for someData := range ch {
		fileSlice = append(fileSlice, someData)
	}

	// obtain slice of clone-files only
	return filterUnique(fileSlice), nil
}

// manager waits for all the enumFiles functions to end their work
func manager(pool <-chan struct{}, ch chan fileData) {
	for i := 0; i < cap(pool); i++ {
		<-pool
	}
	close(ch)
}

// enumFiles enumerates files in the given subdirectory and sends fileData structure about
// each file via the channel
func enumFiles(dirPath string, ch chan<- fileData) {
	files, _ := ioutil.ReadDir(dirPath)
	for _, someFile := range files {
		if !someFile.IsDir() {

			someFileData := new(fileData)
			someFileData.dir = dirPath
			someFileData.fileName = someFile.Name()
			someFileData.sizeInBytes = uint64(someFile.Size())

			ch <- *someFileData
		}
	}
}

// enumDirs enumerates subdirectories in the given folder.
func enumDirs(conf *ConfigType) ([]string, error) {
	// get absolute directory path
	dirPath, err := filepath.Abs(conf.DirPath)
	if err != nil {
		return nil, err
	}

	// check for dirPath existance
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Entered directory was not found ( %s )\n", dirPath)
		} else {
			return nil, err
		}
	}

	conf.DirPath = dirPath

	// get a slice of subdirectories
	var dirSlice []string
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirSlice = append(dirSlice, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirSlice, nil
}

// filterUnique looks throw all the given []fileData and sets "id" field of each element.
// Clone-files are identical if their name and size are equal.
// filterUnique produces a slice with only data about clones and a map with a number of each ID in the given slice.
func filterUnique(fileSlice []fileData) []fileData /*, map[uint16]uint16*/ {
	// search for clones and mark the fileData structures with id that corresponds to the pair Name-Size
	index := uint32(1)
	idMap := make(map[string]uint32)
	cloneMap := make(map[uint32]uint32)
	for i := 0; i < len(fileSlice); i += 1 {
		cloneID := fileSlice[i].fileName + fmt.Sprint(fileSlice[i].sizeInBytes)
		id, ok := idMap[cloneID]
		if !ok {
			//this file is unique
			idMap[cloneID] = index
			fileSlice[i].id = index
			cloneMap[index] += 1
			index += 1
			continue
		}
		fileSlice[i].id = id
		cloneMap[id] += 1
	}
	// count capacity of slice to store only clone's data
	var capacity uint32
	for _, num := range cloneMap {
		if num > 1 {
			capacity += num
		}
	}
	// fill the slice of clones with the data
	fileSliceClones := make([]fileData, capacity)
	i := 0
	for _, someFileData := range fileSlice {
		if cloneMap[someFileData.id] > 1 {
			fileSliceClones[i] = someFileData
			i += 1
		}
	}

	sortData(fileSliceClones)

	return fileSliceClones
}

// sort according to flags user has set
func sortData(fileSlice []fileData) {
	sort.Slice(fileSlice, func(i, j int) bool {
		return fileSlice[i].id < fileSlice[j].id
	})
}
