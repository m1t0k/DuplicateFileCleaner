package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// FileItems struct holds all files with the same hash
type FileItems struct {
	Files []string
}

// GlobalFileStorage holds all information about duplicate files
var GlobalFileStorage = struct {
	sync.RWMutex
	FileStorage map[string]FileItems
}{FileStorage: make(map[string]FileItems)}

func md5sum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		//fmt.Print(err)
		//log.Printf("md5sum: can't open file %filePath. Error: %v", err)
		return
	}
	defer file.Close()
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		//fmt.Print(err)
		//log.Printf("md5sum: can't open file %filePath. Error: %v", err)
		return
	}
	result = hex.EncodeToString(hash.Sum(nil))
	return
}

func processFile(fileName string, wg *sync.WaitGroup) {
	defer (*wg).Done()
	hash, errHash := md5sum(fileName)
	if errHash != nil {
		//log.Printf("Error: can't calculate hash for file %s.\n Error:%v", fileName, errHash)
		return
	}
	GlobalFileStorage.Lock()
	item, ok := GlobalFileStorage.FileStorage[hash]
	if ok == true {
		item.Files = append(item.Files, fileName)
		GlobalFileStorage.FileStorage[hash] = item
	} else {
		var fileItems FileItems
		fileItems.Files = append(GlobalFileStorage.FileStorage[hash].Files, fileName)
		GlobalFileStorage.FileStorage[hash] = fileItems
	}
	GlobalFileStorage.Unlock()
}

func findDuplicates(directory string, filePattern string) {
	subDir := [2]string{"**", ""}
	var files []string
	for _, sub := range subDir {
		dirFiles, err := filepath.Glob(directory + "\\" + sub + "\\" + filePattern)
		if err != nil {
			log.Fatalf("Error: can't scan folder %s with pattern %s.\n Error:%v", directory, filePattern, err)
			return
		}
		log.Printf("Found %d files in %s", len(dirFiles), directory)
		files = append(files, dirFiles...)
	}
	var wg sync.WaitGroup
	wg.Add(len(files))
	for index := 0; index < len(files); index++ {
		go processFile(files[index], &wg)
	}
	wg.Wait()

}

func showResults() {
	GlobalFileStorage.Lock()
	fmt.Printf("%d\n", len(GlobalFileStorage.FileStorage))
	for k, v := range GlobalFileStorage.FileStorage {
		if len(v.Files) > 1 {
			fmt.Printf("key[%s] value[%s]\n", k, v)
		}
	}
	GlobalFileStorage.Unlock()
}

type configurationSettings struct {
	directory   string
	filePattern string
}

func getConfigSettings() configurationSettings {
	var settings configurationSettings

	flag.StringVar(&settings.directory, "dir", "", "-dir=C:\\downloads")
	flag.StringVar(&settings.filePattern, "ext", "*.*", "-ext=*.pdf")
	flag.Parse()

	if len(settings.directory) <= 0 {
		log.Fatal("Root directory is not set:-dir=C:\\downloads")
	}

	return settings
}

func main() {
	settings := getConfigSettings()
	findDuplicates(settings.directory, settings.filePattern)
	showResults()
}
