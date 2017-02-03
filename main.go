package main

func main() {
	settings := getConfigSettings()
	findDuplicates(settings.dirList, settings.filePatterns)
	showResults()
}
