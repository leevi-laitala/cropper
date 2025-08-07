package main

import (
	"fmt"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	multiplier int32 = 1
)

type seekerRect struct {
	width	   int32
	height	   int32
}

func frameToWidth(frame int32) int32 {
	return int32(rl.GetScreenWidth() * int(frame) / (curVideo.Frames() - 1))
}

func updateSeeker(seeker *seekerRect) {
	seeker.width = frameToWidth(int32(curFrame))
}

func resetSeeker(seeker *seekerRect) {
	frameBegin = 0
	frameEnd = int32(curVideo.Frames() - 1)
	seeker.width = 0
	seeker.height = 16
}

//nolint:unparam
func getTextOffsetX(text string, fontsize int32, x int32) int32 {
	textWidth := rl.MeasureText(text, fontsize)
	
	newx := x - textWidth / 2

	if newx < 16 {
		newx = 16
	}

	if newx > int32(rl.GetScreenWidth()) - textWidth - 16 {
		newx = int32(rl.GetScreenWidth()) - textWidth - 16
	}

	return newx
}

func drawSeeker(seeker *seekerRect) {
	rl.DrawRectangle(0, 0, seeker.width, seeker.height, rl.Blue)

	// Font size
	var fontsize int32 = 24

	// Current frame text
	frametext := strconv.Itoa(curFrame)
	rl.DrawText(frametext, getTextOffsetX(frametext, fontsize, seeker.width),
		seeker.height, fontsize, rl.Blue)

	if frameBegin != 0 {
		rl.DrawLine(frameToWidth(frameBegin), 0, frameToWidth(frameBegin),
			seeker.height, rl.Red)
		frametext = fmt.Sprintf("%d ->", frameBegin)
		rl.DrawText(frametext, getTextOffsetX(frametext, fontsize,
			frameToWidth(frameBegin)), seeker.height, fontsize, rl.Red)
	}

	if frameEnd != int32(curVideo.Frames() - 1) {
		rl.DrawLine(frameToWidth(frameEnd), 0, frameToWidth(frameEnd),
			seeker.height, rl.Green)
		frametext = fmt.Sprintf("<- %d", frameEnd)
		rl.DrawText(frametext, getTextOffsetX(frametext, fontsize,
			frameToWidth(frameEnd)), seeker.height, fontsize, rl.Green)
	}

	frametext = fmt.Sprintf("Muted: %t", muted)
	rl.DrawText(frametext, getTextOffsetX(frametext, fontsize, 0),
		int32(rl.GetScreenHeight()) - seeker.height * 2, fontsize, rl.Blue)
}

// Vim style, type numbers for multiplier for action key
// E.g. by pressing "13<left-key>", seeker moves thirteen frames forward.
func updateMultiplier() {
	lastKey := rl.GetKeyPressed()
	if lastKey >= 48 && lastKey <= 57 { // Number key has been pressed
		if multiplier != 0 {
			multiplier *= 10
		}

		multiplier += 9 - (57 - lastKey)
	}
}

func getMultiplier() int {
	// Reset multiplier after calling this function
	defer func() { multiplier = 0 }()

	// Multiplier cannot be 0 as no action would be done then
	if multiplier == 0 {
		return 1
	}

	return int(multiplier)
}

//nolint:cyclop
func handleKeys() {
	if rl.IsKeyPressed(rl.KeyLeft) { // Seek frames forward
		curFrame -= getMultiplier()
	}

	if rl.IsKeyPressed(rl.KeyRight) { // Seek frames backward
		curFrame += getMultiplier()
	}

	if rl.IsKeyPressed(rl.KeyUp) { // Seek seconds forward
		curFrame += int(curVideo.FPS()) * getMultiplier()
	}

	if rl.IsKeyPressed(rl.KeyDown) { // Seek seconds backward
		curFrame -= int(curVideo.FPS()) * getMultiplier()
	}

	// Seek to first frame
	if rl.IsKeyPressed(rl.KeyB) && rl.IsKeyDown(rl.KeyLeftShift) {
		curFrame = 0
	}

	// Seek to last frame
	if rl.IsKeyPressed(rl.KeyE) && rl.IsKeyDown(rl.KeyLeftShift) {
		curFrame = curVideo.Frames() - 1
	}

	if rl.IsKeyPressed(rl.KeyM) { // Toggle mute
		muted = !muted
	}

	if rl.IsKeyPressed(rl.KeyA) { // Set trim point A
		frameBegin = int32(curFrame)
	}

	// Set trim point B
	if rl.IsKeyPressed(rl.KeyB) && !rl.IsKeyDown(rl.KeyLeftShift) {
		frameEnd = int32(curFrame)
	}
}

func evalTrimPoints() {
	// Beginning and end cannot be at the same frame
	if frameBegin == frameEnd {
		frameBegin = 0
	}

	// Swap if wrong way around, like B being before A
	if frameBegin > frameEnd {
		frameEnd, frameBegin = frameBegin, frameEnd
	}
}

// Returns true if frame was updated.
func getSeekValue() bool {
	updated := false

	updateMultiplier()

	prevFrame := curFrame

	handleKeys()

	evalTrimPoints()

	if curFrame != prevFrame {
		updated = true
	}

	if updated {
		if curFrame < 0 {
			curFrame = 0
		}

		if curFrame >= curVideo.Frames() {
			curFrame = curVideo.Frames() - 1
		}
	}

	return updated
}
