package main

import "fmt"

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
