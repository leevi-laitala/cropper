package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"math"
)

type rectArea struct {
	minn, maxx rl.Vector2
	attach     uint8
}

func evaluateAttachedHandle(rect *rectArea, mousepos *rl.Vector2) {
	// Max distance from handle where is considered attached
	treshold := float64(16)

	// Check which handle cursor is attached to
	if math.Abs(float64(rect.minn.X)-float64(mousepos.X)) < treshold {
		if rect.attach&right == 0 {
			rect.attach |= left
		}
	}

	if math.Abs(float64(rect.maxx.X)-float64(mousepos.X)) < treshold {
		if rect.attach&left == 0 {
			rect.attach |= right
		}
	}

	if math.Abs(float64(rect.minn.Y)-float64(mousepos.Y)) < treshold {
		if rect.attach&bottom == 0 {
			rect.attach |= top
		}
	}

	if math.Abs(float64(rect.maxx.Y)-float64(mousepos.Y)) < treshold {
		if rect.attach&top == 0 {
			rect.attach |= bottom
		}
	}
}

func reevaluateAreaCorners(rect *rectArea) {
	realMin := rl.Vector2{
		X: float32(math.Min(float64(rect.minn.X), float64(rect.maxx.X))),
		Y: float32(math.Min(float64(rect.minn.Y), float64(rect.maxx.Y))),
	}
	realMax := rl.Vector2{
		X: float32(math.Max(float64(rect.minn.X), float64(rect.maxx.X))),
		Y: float32(math.Max(float64(rect.minn.Y), float64(rect.maxx.Y))),
	}

	rect.minn = realMin
	rect.maxx = realMax
}

func resizeAreaRect(rect *rectArea, mousepos, constraintMin, constraintMax *rl.Vector2) {
	// If attached, move rectangle
	if rect.attach&left != 0 {
		rect.minn.X = float32(math.Min(
			math.Max(float64(mousepos.X), float64(constraintMin.X)),
			float64(constraintMax.X)))
	}

	if rect.attach&right != 0 {
		rect.maxx.X = float32(math.Min(
			math.Max(float64(mousepos.X), float64(constraintMin.X)),
			float64(constraintMax.X)))
	}

	if rect.attach&top != 0 {
		rect.minn.Y = float32(math.Min(
			math.Max(float64(mousepos.Y), float64(constraintMin.Y)),
			float64(constraintMax.Y)))
	}

	if rect.attach&bottom != 0 {
		rect.maxx.Y = float32(math.Min(
			math.Max(float64(mousepos.Y), float64(constraintMin.Y)),
			float64(constraintMax.Y)))
	}
}

func updateAreaRect(rect *rectArea, cam *rl.Camera2D) {
	if rl.IsMouseButtonDown(0) {
		mousepos := rl.GetScreenToWorld2D(rl.GetMousePosition(), *cam)

		constraintMin := rl.Vector2{X: 0, Y: 0}
		constraintMax := rl.Vector2{
			X: float32(curFrameTex.Width),
			Y: float32(curFrameTex.Height),
		}

		evaluateAttachedHandle(rect, &mousepos)
		resizeAreaRect(rect, &mousepos, &constraintMin, &constraintMax)
	}

	if rl.IsMouseButtonReleased(0) {
		rect.attach = 0
		reevaluateAreaCorners(rect)
	}

	if rl.IsKeyPressed(rl.KeyR) {
		resetAreaRect(rect)
	}
}

func resetAreaRect(rect *rectArea) {
	rect.minn = rl.Vector2{
		X: 0,
		Y: 0,
	}
	rect.maxx = rl.Vector2{
		X: float32(curFrameTex.Width),
		Y: float32(curFrameTex.Height),
	}
}

func drawAreaRect(rect *rectArea) {
	rl.DrawLine(int32(rect.minn.X), int32(rect.minn.Y),
		int32(rect.maxx.X), int32(rect.minn.Y), rl.Red)
	rl.DrawLine(int32(rect.minn.X), int32(rect.minn.Y),
		int32(rect.minn.X), int32(rect.maxx.Y), rl.Red)
	rl.DrawLine(int32(rect.maxx.X), int32(rect.minn.Y),
		int32(rect.maxx.X), int32(rect.maxx.Y), rl.Red)
	rl.DrawLine(int32(rect.minn.X), int32(rect.maxx.Y),
		int32(rect.maxx.X), int32(rect.maxx.Y), rl.Red)
}
