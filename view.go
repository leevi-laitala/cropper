package main

import (
	"log"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
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

		cam.Zoom = float32(math.Max(0.05, math.Min(float64(cam.Zoom)*scaleFactor, 10)))
	}

	if rl.IsMouseButtonDown(1) {
		mouseCurPos := rl.GetMousePosition()
		delta := rl.GetMouseDelta()
		delta = rl.Vector2Scale(delta, -1/cam.Zoom)
		cam.Target = rl.Vector2Add(cam.Target, delta)

		if mouseCurPos.X < 0 {
			rl.SetMousePosition(int(mouseCurPos.X)+screenWidth, int(mouseCurPos.Y))
		}

		if mouseCurPos.X > screenWidth {
			rl.SetMousePosition(int(mouseCurPos.X)-screenWidth, int(mouseCurPos.Y))
		}

		if mouseCurPos.Y < 0 {
			rl.SetMousePosition(int(mouseCurPos.X), int(mouseCurPos.Y)+screenHeight)
		}

		if mouseCurPos.Y > screenHeight {
			rl.SetMousePosition(int(mouseCurPos.X), int(mouseCurPos.Y)-screenHeight)
		}
	}
}

func resetCamera(video *videofile, cam *rl.Camera2D) {
	cam.Target = rl.Vector2{
		X: float32(video.ssTex.Width) / 2,
		Y: float32(video.ssTex.Height) / 2,
	}
	cam.Offset = rl.Vector2{
		X: float32(screenWidth) / 2,
		Y: float32(screenHeight) / 2,
	}
	cam.Zoom = float32(math.Min(
		float64(screenWidth)/float64(video.ssTex.Width),
		float64(screenHeight)/float64(video.ssTex.Height),
	) * 0.9)
}

func run(videoFiles []videofile) {
	rl.InitWindow(screenWidth, screenHeight, "Cropper")
	rl.SetTargetFPS(100)

	currentVideo := int(0)

	err := reloadImage(&videoFiles[currentVideo])
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = updateImageFrame(&videoFiles[currentVideo], false)
	if err != nil {
		log.Fatalf("%v", err)
	}

	cam := rl.Camera2D{}
	resetCamera(&videoFiles[currentVideo], &cam)

	rect := rectArea{}
	resetAreaRect(&rect, &videoFiles[currentVideo])

	for !rl.WindowShouldClose() {
		panAndZoom(&cam)

		updateAreaRect(&rect, &videoFiles[currentVideo], &cam)

		err = updateImageFrame(&videoFiles[currentVideo], false)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if rl.IsKeyPressed(rl.KeyEnter) {
			err := exportCroppedVideo(&videoFiles[currentVideo], &rect)
			if err != nil {
				log.Fatalf("%v", err)
			}

			currentVideo++
			if currentVideo >= len(videoFiles) {
				log.Println("All videos cropped")

				break
			}

			err = updateImageFrame(&videoFiles[currentVideo], true)
			if err != nil {
				log.Fatalf("%v", err)
			}

			resetAreaRect(&rect, &videoFiles[currentVideo])
		}

		if rl.IsKeyPressed(rl.KeyC) {
			resetCamera(&videoFiles[currentVideo], &cam)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Beige)
		rl.BeginMode2D(cam)

		rl.DrawTexture(videoFiles[currentVideo].ssTex, 0, 0, rl.White)
		drawAreaRect(&rect)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
