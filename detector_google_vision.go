package main

import (
	"io/ioutil"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/vision"
	SDK "google.golang.org/api/vision/v1"
)

type GoogleVisionFaceDetector struct {
	client *vision.Vision
}

func NewGoogleVisionFaceDetector() (*GoogleVisionFaceDetector, error) {
	cli, err := vision.New(config.Config{})
	if err != nil {
		return nil, err
	}

	return &GoogleVisionFaceDetector{
		client: cli,
	}, nil
}

func (d GoogleVisionFaceDetector) String() string {
	return "google"
}

func (d GoogleVisionFaceDetector) Detect(imgPath string) (FaceResult, error) {
	imgWidth, imgHeight, err := GetImageSize(imgPath)
	if err != nil {
		return FaceResult{}, err
	}

	byt, err := ioutil.ReadFile(imgPath)
	if err != nil {
		return FaceResult{}, err
	}

	resp, err := d.client.Face(byt)
	if err != nil {
		return FaceResult{}, err
	}

	var list []*SDK.FaceAnnotation
	for _, r := range resp.Responses {
		list = append(list, r.FaceAnnotations...)
	}
	faces := make([]FaceData, len(list))
	for i, r := range list {
		v := r.FdBoundingPoly.Vertices
		if len(v) < 3 {
			continue
		}
		x := v[0].X
		y := v[0].Y
		w := v[2].X - v[0].X
		h := v[2].Y - v[0].Y

		faces[i] = FaceData{
			X:             int(x),
			Y:             int(y),
			Width:         int(w),
			Height:        int(h),
			PercentWidth:  float64(w) / float64(imgWidth),
			PercentHeight: float64(h) / float64(imgHeight),
			Confidence:    r.DetectionConfidence * 100,
		}
	}

	return FaceResult{
		EngineName: d.String(),
		Faces:      faces,
	}, nil
}
