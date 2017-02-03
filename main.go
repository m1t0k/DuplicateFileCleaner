package main

func main() {
	settings := readSettings()
	findDuplicates(settings)
	showResults()
}
