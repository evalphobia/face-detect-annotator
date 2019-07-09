package main

import (
	"errors"
	"fmt"
	"sync"

	"gocv.io/x/gocv"
)

type OpenCVFaceDetector struct {
	classifier gocv.CascadeClassifier
	mu         sync.Mutex
}

func NewOpenCVFaceDetector(xmlPath string) (*OpenCVFaceDetector, error) {
	c := gocv.NewCascadeClassifier()
	if !c.Load(xmlPath) {
		return nil, fmt.Errorf("Error reading cascade file: [%s]", xmlPath)
	}

	return &OpenCVFaceDetector{
		classifier: c,
	}, nil
}

func (d OpenCVFaceDetector) String() string {
	return "opencv"
}

func (d *OpenCVFaceDetector) Detect(imgPath string) (FaceResult, error) {
	imgWidth, imgHeight, err := GetImageSize(imgPath)
	if err != nil {
		return FaceResult{}, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	img := gocv.IMRead(imgPath, gocv.IMReadGrayScale)
	if img.Empty() {
		return FaceResult{}, errors.New("Empty image")
	}

	rects := d.classifier.DetectMultiScale(img)
	faces := make([]FaceData, len(rects))
	for i, r := range rects {
		x := r.Min.X
		y := r.Min.Y
		w := r.Dx()
		h := r.Dy()

		faces[i] = FaceData{
			X:             x,
			Y:             y,
			Width:         w,
			Height:        h,
			PercentWidth:  float64(w) / float64(imgWidth),
			PercentHeight: float64(h) / float64(imgHeight),
		}
	}

	return FaceResult{
		EngineName: d.String(),
		Faces:      faces,
	}, nil
}
