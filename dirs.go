package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
)

func scanDirectory(rootDir string, patterns []string, wg *sync.WaitGroup) {
	defer (*wg).Done()

	dirs, err := ioutil.ReadDir(rootDir)
	if err != nil {
		//log.Fatal(err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			//log.Printf("Start scan in directory %s", dir.Name())
			(*wg).Add(1)
			go scanDirectory(filepath.Join(rootDir, dir.Name()), patterns, wg)
		}
	}

	for _, pattern := range patterns {
		dirFiles, err := filepath.Glob(filepath.Join(rootDir, pattern))
		if err != nil {
			//log.Printf("Error: can't scan folder %s with pattern %s.\n Error:%v", dir+"\\"+sub, pattern, err)
		}

		log.Printf("Found %d files in %s", len(dirFiles), rootDir)

		GlobalFileList.Lock()
		GlobalFileList.FileList = append(GlobalFileList.FileList, dirFiles...)
		GlobalFileList.Unlock()
	}
}

func findDuplicates(settings configSettings) {
	var dirWg sync.WaitGroup

	for _, rootDir := range settings.dirList {
		dirs, err := ioutil.ReadDir(rootDir)
		if err != nil {
			log.Print(err)
		}
		for _, dir := range dirs {
			if dir.IsDir() {
				dirWg.Add(1)
				//log.Printf("Start scan in directory %s", dir.Name())
				go scanDirectory(filepath.Join(rootDir, dir.Name()), settings.filePatterns, &dirWg)
			}
		}

	}
	dirWg.Wait()

	var fileWg sync.WaitGroup
	fileWg.Add(len(GlobalFileList.FileList))
	for index := 0; index < len(GlobalFileList.FileList); index++ {
		go processFile(GlobalFileList.FileList[index], &fileWg)
	}
	fileWg.Wait()
}
