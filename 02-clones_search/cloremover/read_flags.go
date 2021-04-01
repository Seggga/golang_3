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
}

func ReadFlags(conf *ConfigType) error {

	var (
		removeFlag  = flag.Bool("remove", false, "set the flag if you wish to delete some clone files")
		confirmFlag = flag.String("confirm", "on", "confirmation before removing files")
		showFlag    = flag.Int("files", 10, "specifies the amount of different clones to be displayed")
		limitFlag   = flag.Int("dirs", 0, "determines maximum number of directories to be displayed, default is 'no limit'")
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

	conf.DirPath = flag.Arg(0)
	conf.RemoveFlag = *removeFlag
	conf.ConfirmFlag = *confirmFlag
	conf.ShowFiles = uint8(*showFlag)
	conf.DirLimit = uint8(*limitFlag)

	return nil

}
