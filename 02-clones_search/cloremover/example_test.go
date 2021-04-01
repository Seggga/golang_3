package cloremover_test

import (
	"github.com/Seggga/golang_2/08_clones_search/cloremover"
)

func Example() {
	// read flags
	conf := &cloremover.ConfigType{}
	_ = cloremover.ReadFlags(conf)
	// collect data
	fileSlice, _ := cloremover.FindClones(conf)
	// display data
	outputMap := cloremover.PrintClones(fileSlice, conf)
	// remove data
	cloremover.Remove(fileSlice, conf, outputMap)
}
