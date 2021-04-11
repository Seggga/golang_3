package cloremover

import (
	"fmt"
)

// PrintClones prints given data according to selected sort and filter options
func PrintClones(fileSlice []fileData, conf *ConfigType) map[uint32]uint32 {
	// determine the number of clones of each id
	cloneMap := make(map[uint32]uint32)
	for _, fileData := range fileSlice {
		cloneMap[fileData.id] += 1
	}
	// first output
	fmt.Printf("Directory - %s - contains clone-files as folows\n", conf.DirPath)

	var id, showCounter, limitCounter uint32
	outputMap := make(map[uint32]uint32)
	for _, someData := range fileSlice {
		// ID mismatch - means start of another group of clones with the new id
		if id != someData.id {
			showCounter += 1
			if showCounter > uint32(conf.ShowFiles) {
				return outputMap // that's enough to print clone-files
			}

			id = someData.id
			outputMap[showCounter] = id

			fmt.Println()
			fmt.Printf("[%2d]: %s - %d bytes, %3d clones:\n", showCounter, someData.fileName, someData.sizeInBytes, cloneMap[id])
			limitCounter = 1
		}
		if conf.DirLimit > 0 && limitCounter > uint32(conf.DirLimit) {
			continue // that's enough to print directories for theese clone-files
		}
		fmt.Println(someData.dir)
		limitCounter += 1
	}

	return outputMap
}
