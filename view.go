package main

import (
	"log"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	v "github.com/AlexEidt/Vidio"
)

func panAndZoom(cam *rl.Camera2D) {
	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		mouseWorld := rl.GetScreenToWorld2D(rl.GetMousePosition(), *cam)
		cam.Offset = rl.GetMousePosition()
		cam.Target = mouseWorld

		scaleFactor := 1 + (0.25 * math.Abs(float64(wheel)))

		if wheel < 0 {
			scaleFactor = 1 / scaleFactor
		}

		cam.Zoom = float32(math.Max(0.05, math.Min(float64(cam.Zoom) * scaleFactor, 10)))
	}

	if rl.IsMouseButtonDown(1) {
		mouseCurPos := rl.GetMousePosition()
		delta := rl.GetMouseDelta()
		delta = rl.Vector2Scale(delta, -1/cam.Zoom)
		cam.Target = rl.Vector2Add(cam.Target, delta)

		if int32(mouseCurPos.X) < 0 {
			rl.SetMousePosition(int(mouseCurPos.X) + int(screenWidth), int(mouseCurPos.Y))
		}

		if int32(mouseCurPos.X) > screenWidth {
			rl.SetMousePosition(int(mouseCurPos.X) - int(screenWidth), int(mouseCurPos.Y))
		}

		if int32(mouseCurPos.Y) < 0 {
			rl.SetMousePosition(int(mouseCurPos.X), int(mouseCurPos.Y) + int(screenHeight))
		}

		if int32(mouseCurPos.Y) > screenHeight {
			rl.SetMousePosition(int(mouseCurPos.X), int(mouseCurPos.Y) - int(screenHeight))
		}
	}
}

func resetCamera(cam *rl.Camera2D) {
	cam.Target = rl.Vector2{
		X: float32(curFrameTex.Width) / 2,
		Y: float32(curFrameTex.Height) / 2,
	}
	cam.Offset = rl.Vector2{
		X: float32(screenWidth) / 2,
		Y: float32(screenHeight) / 2,
	}
	cam.Zoom = float32(math.Min(
		float64(screenWidth)/float64(curFrameTex.Width),
		float64(screenHeight)/float64(curFrameTex.Height),
	) * 0.9)
}

func initNewVideo() error {
	if curVideo != nil {
		curVideo.Close()
	}

	curFrame = 0

	var err error = nil

	curVideo, err = v.NewVideo(videoFiles[curVideoIndex])
	if err != nil {
		return fmt.Errorf("error loading new video: %w", err)
	}

	err = loadVideoFrameToTexture()
	if err != nil {
		return fmt.Errorf("error loading frame to texture: %w", err)
	}

	return nil
}

func run() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(screenWidth, screenHeight, "Cropper")
	rl.SetTargetFPS(100)

	err := initNewVideo()
	if err != nil {
		log.Fatalf("failed to load first video: %w", err)
	}

	cam := rl.Camera2D{}
	resetCamera(&cam)

	rect := rectArea{}
	resetAreaRect(&rect)

	seeker := seekerRect{}
	resetSeeker(&seeker)

	for !rl.WindowShouldClose() {
		panAndZoom(&cam)

		updateAreaRect(&rect, &cam)

		frameUpdated := getSeekValue()
		if frameUpdated {
			updateSeeker(&seeker)

			err := loadVideoFrameToTexture()
			if err != nil {
				curVideo.Close()
				log.Fatalf("error loading frame to texture: %w", err)
			}
		}

		// Enter key exports current video and loads in next video
		if rl.IsKeyPressed(rl.KeyEnter) {
			// Export current video
			go exportCroppedVideo(videoFiles[curVideoIndex], muted, frameBegin, frameEnd, int32(curVideo.Frames()), curVideo.FPS(), rect)

			curVideoIndex++
			if curVideoIndex >= len(videoFiles) {
				log.Println("All videos cropped")

				break
			}

			resetAreaRect(&rect)
			resetCamera(&cam)

			initNewVideo()

			resetSeeker(&seeker)
		}

		if rl.IsKeyPressed(rl.KeyC) {
			resetCamera(&cam)
		}

		if rl.IsWindowResized() {
			screenWidth = int32(rl.GetScreenWidth())
			screenHeight = int32(rl.GetScreenHeight())

			resetCamera(&cam)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Beige)
		rl.BeginMode2D(cam)

		rl.DrawTexture(curFrameTex, 0, 0, rl.White)
		drawAreaRect(&rect)

		rl.EndMode2D()

		drawSeeker(&seeker, &cam)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
