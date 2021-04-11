// Package cloremover contains a number of functions to find clone-files.
// The aim is to quickly scan the given directory and let user to safely delete some clone files.
//
// ReadFlags function parses user's input when the program starts.
// It returns errors if the input is incorrect.
//
//	ReadFlags(conf *ConfigType) error
//
// FindClones function reads all the files in given directory and it's subdirectories.
// As a result it returns a slice with data about files with same name and size.
//
//	FindClones(conf *ConfigType) ([]fileData, error)
//
// The PrintClones function produces data output according to specified flags.
// A map in the output contains the data about clone files user would probably like to remove.
//
//	PrintClones(fileSlice []fileData, conf *ConfigType) map[uint32]uint32
//
// The Remove function deletes chosen file. It asks the user which file he want to delete.
// If user's input was correct, the function deletes the file.
//
//	Remove(fileSlice []fileData, conf *ConfigType, outputMap map[uint32]uint32)
//
package cloremover

type fileData struct {
	dir         string
	fileName    string
	sizeInBytes uint64
	id          uint32
}

type ConfigType struct {
	DirPath     string
	RemoveFlag  bool
	ConfirmFlag string
	ShowFiles   uint8
	DirLimit    uint8
	LogFile     string
	LogLevel    string
}
