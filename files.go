package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	_ "log"
	"os"
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

/// GlobalFileList stores all files found by pattern
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
