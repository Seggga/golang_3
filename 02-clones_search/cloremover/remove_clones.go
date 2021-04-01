package cloremover

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func Remove(fileSlice []fileData, conf *ConfigType, outputMap map[uint32]uint32) {
	// removeFlag was not set
	if !conf.RemoveFlag {
		return
	}

	userChoice, err := chooseFile(uint32(len(outputMap)))
	if err != nil {
		fmt.Println(err)
		return
	}

	// print data about chosen file
	onceFlag := false
	dirMap := make(map[uint32]uint32)
	var showCounter uint32
	for i, fileData := range fileSlice {
		if fileData.id == outputMap[uint32(userChoice)] {
			if !onceFlag {
				fmt.Println()
				fmt.Printf("Chosen file: %s - %d bytes:\n", fileData.fileName, fileData.sizeInBytes)
				onceFlag = true
			}

			showCounter += 1
			fmt.Printf("  - %3d - %s\n", showCounter, fileData.dir)
			dirMap[showCounter] = uint32(i)
		}
	}

	userChoice, err = chooseDir(uint32(len(dirMap)))
	if err != nil {
		fmt.Println(err)
		return
	}

	fileToDelete := fileSlice[dirMap[userChoice]]
	// delete the file !!!!!!!!!!
	if conf.ConfirmFlag == "on" {
		confirmation := confirmRemove(fileToDelete)
		if !confirmation {
			fmt.Println("file was not deleted")
			return
		}
	}

	filePath := filepath.Join(fileToDelete.dir, fileToDelete.fileName)
	if err := os.Remove(filePath); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("file %s was successfully deleted", filePath)
}

// user dialog to choose a file
func chooseFile(numOfFiles uint32) (uint32, error) {

	var userInput string
	fmt.Print("Please, choose a number of file you want to delete (for example '3'). For exit enter 'q': ")
	_, err := fmt.Scanln(&userInput)
	if err != nil {
		return 0, fmt.Errorf("There is an error entering data.\n%v", err)
	}
	// user want to quit
	if userInput == "q" || userInput == "Q" {
		return 0, fmt.Errorf("Entered 'q', no files to be deleted. Program exit.")
	}
	// user want to view specified file
	userChoice, err := strconv.ParseInt(userInput, 0, 32)
	if err != nil {
		return 0, fmt.Errorf("Entered data cannot be recognized as number, %v", err)
	}
	// check user input
	if userChoice > int64(numOfFiles) || userChoice < 1 {
		return 0, fmt.Errorf("Invalid input: expected a number (1...%d), got %d", numOfFiles, userChoice)
	}

	return uint32(userChoice), nil
}

// user dialog to choose directory
func chooseDir(numOfDirs uint32) (uint32, error) {
	var userInput string
	// choose and remove a file
	fmt.Print("Please, choose a number of directory you want the file to be removed from (for example '2'). For exit enter 'q': ")
	_, err := fmt.Scanln(&userInput)
	if err != nil {
		return 0, fmt.Errorf("There is an error entering data.\n%v", err)
	}
	// user want to quit
	if userInput == "q" || userInput == "Q" {
		return 0, fmt.Errorf("Entered 'q', no files to be deleted. Program exit.")
	}
	// user want to delete a file frome specified directory
	userChoice, err := strconv.ParseInt(userInput, 0, 32)
	if err != nil {
		return 0, fmt.Errorf("Entered data cannot be recognized as number, %v", err)
	}
	// check user input
	if userChoice > int64(numOfDirs) || userChoice <= 0 {
		return 0, fmt.Errorf("Invalid input: expected a number (0...%d), got %d", numOfDirs, userChoice)
	}

	return uint32(userChoice), nil
}

// user dialog to confirm dile removing
func confirmRemove(fileToDelete fileData) bool {
	fmt.Printf("you are about to delete the file %s from directory %s, %d bytes in size", fileToDelete.fileName, fileToDelete.dir, fileToDelete.sizeInBytes)
	var userInput string
	// choose and remove a file
	fmt.Print("Please, enter 'yes' to confirm file removing. For exit enter 'q': ")
	_, err := fmt.Scanln(&userInput)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// user want to quit
	if userInput == "q" || userInput == "Q" {
		fmt.Println("Entered 'q', no files to be deleted. Program exit.")
		return false
	}
	// user want to delete a file frome specified directory
	if userInput == "yes" {
		return true
	}
	return false

}
