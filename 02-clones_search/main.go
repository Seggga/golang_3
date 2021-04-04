package main

import (
	"fmt"
	"os"

	"github.com/Seggga/golang_3/02-clones_search/cloremover"
	"github.com/sirupsen/logrus"
)

func main() {

	// read configuration from flags
	conf := &cloremover.ConfigType{}
	if err := cloremover.ReadFlags(conf); err != nil {
		fmt.Println(err)
		return
	}

	// initialize logger
	log := logrus.New()
	f, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.SEEK_CUR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch conf.LogLevel {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	}
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(f)
	log.Info("logging started")
	log.Info("command-line parameters:")
	log.Infof("cl-param directory: %s", conf.DirPath)
	log.Infof("cl-param remove flag: %t", conf.RemoveFlag)
	log.Infof("cl-param confirm flag: %s", conf.ConfirmFlag)
	log.Infof("cl-param show files: %d", conf.ShowFiles)
	log.Infof("cl-param directory limit: %d", conf.DirLimit)
	log.Infof("cl-param logfile: %s", conf.LogFile)

	// collect data
	fileSlice, err := cloremover.FindClones(conf, log) // add debug level
	if err != nil {
		log.Errorf("error in FindClones: %v", err)
		fmt.Println(err)
		return
	}
	// no clone files were found
	if len(fileSlice) == 0 {
		log.Info("no clone files were found in the given directory")
		fmt.Println("no clone files were found in the given directory")
		return
	}

	// display data
	log.Debug("list of clones has been constructed, show results")
	outputMap := cloremover.PrintClones(fileSlice, conf)

	// remove data
	cloremover.Remove(fileSlice, conf, outputMap, log)
	log.Info("Program exit.")
}
