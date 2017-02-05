package main

func main() {
	settings := readSettings()
	setRuntime()
	findDuplicates(settings)
	showResults()
}
