package main

import "os"

func main() {
	os.Exit(0) // want "calling os.Exit in the main function is forbidden in package main"
}
