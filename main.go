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
	"strings"
	"sync"
)

// FileItems struct holds all files with the same hash
type FileItems struct {
	Files []string
}

// GlobalDuplicateFileList holds all information about duplicate files
var GlobalDuplicateFileList = struct {
	sync.RWMutex
	FileStorage map[string]FileItems
}{FileStorage: make(map[string]FileItems)}

var GlobalFileList = struct {
	sync.RWMutex
	FileList []string
}{}

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
	GlobalDuplicateFileList.Lock()
	item, ok := GlobalDuplicateFileList.FileStorage[hash]
	if ok == true {
		item.Files = append(item.Files, fileName)
		GlobalDuplicateFileList.FileStorage[hash] = item
	} else {
		var fileItems FileItems
		fileItems.Files = append(GlobalDuplicateFileList.FileStorage[hash].Files, fileName)
		GlobalDuplicateFileList.FileStorage[hash] = fileItems
	}
	GlobalDuplicateFileList.Unlock()
}

func scanDirectory(dir string, pattern string, wg *sync.WaitGroup) {
	defer (*wg).Done()

	var files []string
	subDir := [2]string{"**", ""}

	for _, sub := range subDir {
		fullDir := dir + "\\" + sub + "\\" + pattern
		dirFiles, err := filepath.Glob(fullDir)
		if err != nil {
			//log.Fatalf("Error: can't scan folder %s with pattern %s.\n Error:%v", directory, filePattern, err)
			return
		}
		log.Printf("Found %d files in %s", len(dirFiles), fullDir)
		files = append(files, dirFiles...)
	}
	GlobalFileList.Lock()
	GlobalFileList.FileList = append(GlobalFileList.FileList, files...)
	GlobalFileList.Unlock()

}

func findDuplicates(dirList []string, filePatterns []string) {
	var dirWg sync.WaitGroup
	dirWg.Add(len(dirList) * len(filePatterns))
	for _, dir := range dirList {
		for _, pattern := range filePatterns {
			go scanDirectory(dir, pattern, &dirWg)
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

func showResults() {
	GlobalDuplicateFileList.Lock()
	fmt.Printf("%d\n", len(GlobalDuplicateFileList.FileStorage))
	for _, v := range GlobalDuplicateFileList.FileStorage {
		if len(v.Files) > 1 {
			fmt.Printf("value[%s]\n", v)
		}
	}
	GlobalDuplicateFileList.Unlock()
}

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

func main() {
	settings := getConfigSettings()
	findDuplicates(settings.dirList, settings.filePatterns)
	showResults()
}
