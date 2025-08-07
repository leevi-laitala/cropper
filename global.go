package main

import (
	v "github.com/AlexEidt/Vidio"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	outSuffix = "_cropped"
	ssSuffix  = "_screenshot"

	videoFormats = `(?i)(\.mp4|\.mkv|\.webm)`

	left, right, top, bottom uint8 = 0x1, 0x2, 0x4, 0x8
)

var (
	screenWidth  int32 = 1280
	screenHeight int32 = 720

	videoFiles    []string
	curVideoIndex int

	curVideo    *v.Video
	curFrame    int
	curFrameTex rl.Texture2D

	frameBegin int32
	frameEnd   int32
	muted      = false
)
