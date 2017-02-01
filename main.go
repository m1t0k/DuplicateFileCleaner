package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// FileItems struct holds all files with the same hash
type FileItems struct {
	Files []string
}

// FileStorage holds all information about duplicate files
type FileStorage map[string]FileItems

func md5sum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Print(err)
		log.Printf("md5sum: can't open file %filePath. Error: %v", err)
		return
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		fmt.Print(err)
		log.Printf("md5sum: can't open file %filePath. Error: %v", err)

		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}

func findDuplicates(directory string, filePattern string, fileStorage FileStorage) {
	files, err := filepath.Glob(directory + "\\" + filePattern)
	if err != nil {
		log.Fatalf("Error: can't scan folder %s with pattern %s.\n Error:%v", directory, filePattern, err)
		return
	}
	log.Printf("Found %d files in %s", len(files), directory)

	for index := 0; index < len(files); index++ {
		fileName := files[index]
		hash, errHash := md5sum(fileName)
		if errHash != nil {
			log.Printf("Error: can't calculate hash for file %s.\n Error:%v", fileName, err)
			continue
		}
		item, ok := fileStorage[hash]
		if ok == true {
			item.Files = append(item.Files, fileName)
			fileStorage[hash] = item
		} else {
			var fileItems FileItems
			fileItems.Files = append(fileStorage[hash].Files, fileName)
			fileStorage[hash] = fileItems
		}
	}
}

func main() {
	fileStorage := make(FileStorage)
	findDuplicates("C:\\Books\\*", "*.pdf", fileStorage)

	fmt.Printf("%d\n", len(fileStorage))

	for k, v := range fileStorage {
		if len(v.Files) > 1 {
			fmt.Printf("key[%s] value[%s]\n", k, v)
		}
	}
}
