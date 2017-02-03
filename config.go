package main

import (
	"flag"
	"log"
	"strings"
)

type configurationSettings struct {
	dirList      []string
	filePatterns []string
}

func getConfigSettings() configurationSettings {
	var settings configurationSettings
	var dirs string
	var files string
	flag.StringVar(&dirs, "dir", "", "-dir=dir,dir2")
	flag.StringVar(&files, "ext", "*.*", "-ext=*.ext1,*.ext2")
	flag.Parse()

	if len(dirs) <= 0 {
		log.Fatal("Root directory is not set:-dir=C:\\dir1,dir2")
	}

	settings.dirList = strings.Split(dirs, ",")
	settings.filePatterns = strings.Split(files, ",")

	isSet := false
	for _, dir := range settings.dirList {
		if len(dir) > 0 && isSet == false {
			isSet = true
		}
	}
	if isSet == false {
		log.Fatal("Root directory is not set:-dir=C:\\dir1,dir2")
	}

	return settings
}
