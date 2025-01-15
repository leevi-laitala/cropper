package main

import (
	"log"
)

func main() {
	videoFiles, err := populateVideoFiles()
	if err != nil {
		log.Fatalln(err)

		return
	}

	run(videoFiles)

	cleanScreenshots(videoFiles)
}
