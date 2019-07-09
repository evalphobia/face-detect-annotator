package main

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg"
)

type FaceResult struct {
	EngineName string     `json:"engine"`
	Faces      []FaceData `json:"faces"`
}

func (r FaceResult) HasFaces() bool {
	return len(r.Faces) != 0
}

func (r FaceResult) ShowOutput() string {
	faceCount := len(r.Faces)
	byt, _ := json.Marshal(r)
	return fmt.Sprintf("%d\t%s", faceCount, string(byt))
}

type FaceData struct {
	X             int     `json:"x"`
	Y             int     `json:"y"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	PercentWidth  float64 `json:"width_per"`
	PercentHeight float64 `json:"height_per"`
	Confidence    float64 `json:"confidence"`
}

func (d FaceData) String() string {
	return fmt.Sprintf("[X:%d Y:%d W:%d H:%d PW:%f PH:%f Confidence:%f]", d.X, d.Y, d.Width, d.Height, d.PercentWidth, d.PercentHeight, d.Confidence)
}

func (d FaceData) SizeString() string {
	return fmt.Sprintf("%d,%d %d,%d", d.X, d.Y, d.Width, d.Height)
}

func (d FaceData) PercentString() string {
	return fmt.Sprintf("%d,%d", d.PercentWidth, d.PercentHeight)
}

func (d FaceData) ToJson() string {
	byt, _ := json.Marshal(d)
	return string(byt)
}

func (d FaceData) MaxX() int {
	return d.X + d.Width
}

func (d FaceData) MaxY() int {
	return d.Y + d.Height
}
