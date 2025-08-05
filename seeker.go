package main

import (
	"fmt"

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

func getTextOffsetX(text string, fontsize int32, x int32) int32 {
	textWidth := rl.MeasureText(text, fontsize)
	
	var xx int32 = x - textWidth / 2

	if xx < 16 {
		xx = 16
	}

	if xx > int32(rl.GetScreenWidth()) - textWidth - 16 {
		xx = int32(rl.GetScreenWidth()) - textWidth - 16
	}

	return xx
}

func drawSeeker(seeker *seekerRect, cam *rl.Camera2D) {
	rl.DrawRectangle(0, 0, seeker.width, seeker.height, rl.Blue)

	// Font size
	var fs int32 = 24

	// Current frame text
	t := fmt.Sprintf("%d", curFrame)
	rl.DrawText(t, getTextOffsetX(t, fs, seeker.width), seeker.height, fs, rl.Blue)

	if frameBegin != 0 {
		rl.DrawLine(frameToWidth(frameBegin), 0, frameToWidth(frameBegin), seeker.height, rl.Red)
		t = fmt.Sprintf("%d ->", frameBegin)
		rl.DrawText(t, getTextOffsetX(t, fs, frameToWidth(frameBegin)), seeker.height, fs, rl.Red)
	}

	if frameEnd != int32(curVideo.Frames() - 1) {
		rl.DrawLine(frameToWidth(frameEnd), 0, frameToWidth(frameEnd), seeker.height, rl.Green)
		t = fmt.Sprintf("<- %d", frameEnd)
		rl.DrawText(t, getTextOffsetX(t, fs, frameToWidth(frameEnd)), seeker.height, fs, rl.Green)
	}

	t = fmt.Sprintf("Muted: %t", muted)
	rl.DrawText(t, getTextOffsetX(t, fs, 0), int32(rl.GetScreenHeight()) - seeker.height * 2, fs, rl.Blue)
}

// Vim style, type numbers for multiplier for action key
// E.g. by pressing "13<left-key>", seeker moves thirteen frames forward
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

// Returns true if frame was updated
func getSeekValue() bool {
	updated := false

	updateMultiplier()

	prevFrame := curFrame

	if rl.IsKeyPressed(rl.KeyLeft) {
		curFrame -= getMultiplier()
	} else if rl.IsKeyPressed(rl.KeyRight) {
		curFrame += getMultiplier()
	} else if rl.IsKeyPressed(rl.KeyUp) {
		curFrame += int(curVideo.FPS()) * getMultiplier()
	} else if rl.IsKeyPressed(rl.KeyDown) {
		curFrame -= int(curVideo.FPS()) * getMultiplier()
	} else if rl.IsKeyDown(rl.KeyLeftShift) && rl.IsKeyPressed(rl.KeyE) {
		curFrame = int(curVideo.Frames()) - 1
	} else if rl.IsKeyDown(rl.KeyLeftShift) && rl.IsKeyPressed(rl.KeyB) {
		curFrame = 0
	} else if rl.IsKeyPressed(rl.KeyM) {
		muted = !muted
	}

	if rl.IsKeyPressed(rl.KeyA) {
		frameBegin = int32(curFrame)
	}
	if rl.IsKeyPressed(rl.KeyB) {
		frameEnd = int32(curFrame)
	}

	if frameBegin > frameEnd {
		tmp := frameEnd
		frameEnd = frameBegin
		frameBegin = tmp
	}

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
