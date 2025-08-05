package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"os"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func loadVideoFrameToTexture() error {
	if curFrameTex.Width != 0 {
		rl.UnloadTexture(curFrameTex)
	}

	err := curVideo.ReadFrame(int(curFrame))
	if err != nil {
		return fmt.Errorf("could not read frame '%d' to framebuffer: %w", curFrame, err)
	}

	img := rl.NewImage(curVideo.FrameBuffer(),
					   int32(curVideo.Width()),
					   int32(curVideo.Height()),
					   1,
					   rl.UncompressedR8g8b8a8)

	curFrameTex = rl.LoadTextureFromImage(img)

	return nil
}

// Convert frames to HH:MM:SS.mmm timestamp from ffmpeg
func convertFramesToTimestamp(frames int32, fps float64) string {
	var frameLength float64 = 1 / fps * 1000

	var milliseconds int32 = int32(frameLength * float64(frames))

	var mmm int32 = milliseconds % 1000
	var SS int32 = (milliseconds / 1000) % 60
	var MM int32 = (milliseconds / 1000 / 60) % 60
	var HH int32 = (milliseconds / 1000 / 60 / 60) % 60

	return fmt.Sprintf("%.2d:%.2d:%.2d.%.3d", HH, MM, SS, mmm)
}

func exportCroppedVideo(fname string, mute bool, framesFrom int32, framesTo int32, framesTotal int32, fps float64, rect rectArea) {
	geometry := fmt.Sprintf("%d:%d:%d:%d",
		int32(rect.maxx.X - rect.minn.X),
		int32(rect.maxx.Y - rect.minn.Y),
		int32(rect.minn.X),
		int32(rect.minn.Y))

	fnameWithoutExtension := strings.ReplaceAll(filepath.Base(fname),
		filepath.Ext(fname), "")

	outfname := fmt.Sprintf("%s%s%s",
		fnameWithoutExtension,
		outSuffix,
		filepath.Ext(fname))

	if len(os.Args) == 2 {
		outfname = filepath.Join(os.Args[1], outfname)
	}

	cropArgs := ffmpeg.KwArgs{
		"filter:v": "crop=" + geometry,
	}
	muteArgs := ffmpeg.KwArgs{
		"an": "",
	}
	trimArgs := ffmpeg.KwArgs{
		"ss": convertFramesToTimestamp(framesFrom, fps),
		"to": convertFramesToTimestamp(framesTo, fps),
	}

	usedArgs := cropArgs

	if mute {
		usedArgs = ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{usedArgs, muteArgs})
	}

	if framesFrom != 0 || framesTo != int32(fps) {
		usedArgs = ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{usedArgs, trimArgs})
	}

	err := ffmpeg.Input(fname, ffmpeg.KwArgs{}).
		Output(outfname, usedArgs).OverWriteOutput().Run()

    if err != nil {
	    log.Fatalf("could not crop video '%s': %w", fname, err)
    }
}


// Go through cwd or dir from first cmd argument and find applicable video files
func populateVideoFiles() error {
	var dir string = "./"

	if len(os.Args) == 2 {
		dir = os.Args[1]
	}

	dirContents, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("could not read dir: %w", err)
	}

	fnameValidator := regexp.MustCompile(videoFormats)

	for _, file := range dirContents {
		if !fnameValidator.MatchString(file.Name()) {
			continue
		}

		entry := filepath.Join(dir, file.Name())

		videoFiles = append(videoFiles, entry)
	}

	if len(videoFiles) == 0 {
		return errors.New("no video files found")
	}

	return nil
}
