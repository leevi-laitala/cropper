package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	v "github.com/AlexEidt/Vidio"
)

const (
	outSuffix                      = "_cropped"

	videoFormats				   = `(?i)(\.mp4|\.mkv|\.webm)`

	left, right, top, bottom uint8 = 0x1, 0x2, 0x4, 0x8
)

var (
	screenWidth int32 = 1280
	screenHeight int32 = 720

	videoFiles []string
	curVideoIndex int

	curVideo *v.Video
	curFrame int = 0
	curFrameTex rl.Texture2D

	frameBegin int32
	frameEnd   int32
	muted      bool = false
)
