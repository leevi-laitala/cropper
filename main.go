package main

import (
	"log"
)

func main() {
	err := populateVideoFiles()
	if err != nil {
		log.Fatalln(err)
		return
	}

	run()
}
