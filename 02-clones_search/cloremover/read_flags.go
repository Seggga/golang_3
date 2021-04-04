package cloremover

import (
	"flag"
	"fmt"
)

type ConfigType struct {
	DirPath     string
	RemoveFlag  bool
	ConfirmFlag string
	ShowFiles   uint8
	DirLimit    uint8
	LogFile     string
	LogLevel    string
}

func ReadFlags(conf *ConfigType) error {

	var (
		removeFlag  = flag.Bool("remove", false, "set this flag if you wish to delete some clone files")
		confirmFlag = flag.String("confirm", "on", "confirmation before removing files")
		showFlag    = flag.Int("files", 10, "specifies the amount of different clones to be displayed")
		limitFlag   = flag.Int("dirs", 0, "determines maximum number of directories to be displayed, default is 'no limit'")
		logFile     = flag.String("log", "cloremover.log", "the name for log-file")
		logLevel    = flag.String("loglevel", "info", "determines what kind of events should be logged")
	)

	flag.Parse()

	// data validation for confirmFlag
	if *confirmFlag != "off" && *confirmFlag != "on" {
		return fmt.Errorf("Incorrect 'confirm' flag value. Expected on/off, got %s", *confirmFlag)
	}
	if *removeFlag == false && *confirmFlag == "off" {
		return fmt.Errorf("You did not set -remove flag, but entered -confirm 'off'")
	}
	// data validation for showFlag
	if *showFlag < 0 || *showFlag > 255 {
		return fmt.Errorf("Incorrect 'show' flag value. Expected value >=0, got %d", *showFlag)
	}
	// data validation for limitFlag
	if *limitFlag < 0 || *limitFlag > 255 {
		return fmt.Errorf("Incorrect 'limit' flag value. Expected value >=0, got %d", *limitFlag)
	}
	// data validation for loglevel
	if *logLevel != "info" && *logLevel != "error" && *logLevel != "debug" {
		return fmt.Errorf("Incorrect 'loglevel' flag value. Expected debug/info/error, got %s", *logLevel)
	}

	conf.DirPath = flag.Arg(0)
	conf.RemoveFlag = *removeFlag
	conf.ConfirmFlag = *confirmFlag
	conf.ShowFiles = uint8(*showFlag)
	conf.DirLimit = uint8(*limitFlag)
	conf.LogFile = *logFile
	conf.LogLevel = *logLevel

	return nil

}
