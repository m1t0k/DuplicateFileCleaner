package main

import "fmt"

func showResults() {

	fmt.Println("|||============================================|||")
	fmt.Println("\t\tList of duplicates:")
	fmt.Println("|||============================================|||")

	GlobalDuplicateFileList.Lock()
	for _, v := range GlobalDuplicateFileList.FileStorage {
		if len(v.Files) > 1 {
			fmt.Printf("value[%s]\n", v)
		}
	}
	GlobalDuplicateFileList.Unlock()

	fmt.Println("|||============================================|||")
	fmt.Println("|||============================================|||")
}
