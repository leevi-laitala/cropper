package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type videofile struct {
	fname, ssfname string
	currentFrame   int16
	ssTex          rl.Texture2D
}

func exportCroppedVideo(video *videofile, rect *rectArea) {
	geometry := fmt.Sprintf("%d:%d:%d:%d",
		int32(rect.maxx.X-rect.minn.X),
		int32(rect.maxx.Y-rect.minn.Y),
		int32(rect.minn.X),
		int32(rect.minn.Y))

	fnameWithoutExtension := strings.ReplaceAll(filepath.Base(video.fname),
		filepath.Ext(video.fname), "")

	outfname := fmt.Sprintf("%s%s%s",
		fnameWithoutExtension,
		outSuffix,
		filepath.Ext(video.fname))

	err := ffmpeg.Input(video.fname, ffmpeg.KwArgs{}).
		Output(outfname, ffmpeg.KwArgs{
			"filter:v": "crop=" + geometry,
		}).OverWriteOutput().Run()

    if err != nil {
	    log.Fatalf("could not crop video '%s': %w", video.fname, err)
    }
}

func getFrame(video *videofile) error {
	var ti time.Time

	d := time.Duration(video.currentFrame) * time.Second
	t := ti.Add(d)

	timestamp := t.Format("00:04:05")

	video.ssfname = fmt.Sprintf("%s%s%s", video.fname, ssSuffix, ssFormat)
	log.Printf("creating screenshot from '%s' at %s, saving to '%s'\n",
		video.fname, timestamp, video.ssfname)

	err := ffmpeg.Input(video.fname, ffmpeg.KwArgs{"ss": timestamp}).
		Output(video.ssfname, ffmpeg.KwArgs{"frames:v": 1, "q:v": 2}).
		OverWriteOutput().Run()
	if err != nil {
		return fmt.Errorf("could not save screenshot: %w", err)
	}

	return nil
}

func loadFrameToTexture(video *videofile) {
	src := rl.LoadImage(video.ssfname)
	defer rl.UnloadImage(src)

	rl.UnloadTexture(video.ssTex)
	video.ssTex = rl.LoadTextureFromImage(src)
}

func updateImageFrame(video *videofile, forceReload bool) error {
	var err error

	if rl.IsKeyPressed(rl.KeyLeft) {
		video.currentFrame--

		if video.currentFrame < 0 {
			video.currentFrame = 0
		}

		err = reloadImage(video)
	} else if rl.IsKeyPressed(rl.KeyRight) {
		video.currentFrame++

		err = reloadImage(video)
	}

	if forceReload {
		err = reloadImage(video)
	}

	if err != nil {
		return fmt.Errorf("could not update image on viewport: %w", err)
	}

	rl.SetWindowTitle("Cropper - " + video.fname)

	return nil
}

func reloadImage(video *videofile) error {
	err := getFrame(video)
	if err != nil {
		return fmt.Errorf("could not create screenshot from '%s': %w",
			video.fname, err)
	}

	loadFrameToTexture(video)

	return nil
}

func cleanScreenshots(videoFiles []videofile) {
	for _, video := range videoFiles {
		var err error

		if video.ssfname != "" {
			err = os.Remove(video.ssfname)
		}

		if err != nil {
			log.Fatalf("could not remove file '%s': %v", video.ssfname, err)
		}
	}
}

func populateVideoFiles() ([]videofile, error) {
	var files []videofile

	allFiles, err := os.ReadDir("./")
	if err != nil {
		return files, fmt.Errorf("could not read dir: %w", err)
	}

	for _, file := range allFiles {
		for _, format := range strings.Split(videoformats, videoformatsSep) {
			if filepath.Ext(file.Name()) == format {
				files = append(files, videofile{
					fname: file.Name(),
				})

				break
			}
		}
	}

	if len(files) == 0 {
		return files, errors.New("no video files found")
	}

	return files, nil
}
