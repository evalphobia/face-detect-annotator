package main

import (
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/rekognition"
)

type RekognitionFaceDetector struct {
	client *rekognition.Rekognition
}

func NewRekognitionFaceDetector() (*RekognitionFaceDetector, error) {
	cli, err := rekognition.New(config.Config{})
	if err != nil {
		return nil, err
	}

	return &RekognitionFaceDetector{
		client: cli,
	}, nil
}

func (d RekognitionFaceDetector) String() string {
	return "rekognition"
}

func (d RekognitionFaceDetector) Detect(imgPath string) (FaceResult, error) {
	imgWidth, imgHeight, err := GetImageSize(imgPath)
	if err != nil {
		return FaceResult{}, err
	}

	resp, err := d.client.DetectFacesFromLocalFile(imgPath)
	if err != nil {
		return FaceResult{}, err
	}

	faces := make([]FaceData, len(resp.List))
	for i, r := range resp.List {
		x := r.BoundingLeft * float64(imgWidth)
		y := r.BoundingTop * float64(imgHeight)
		pw := r.BoundingWidth
		ph := r.BoundingHeight

		faces[i] = FaceData{
			X:             int(x),
			Y:             int(y),
			Width:         int(float64(imgWidth) * pw),
			Height:        int(float64(imgHeight) * ph),
			PercentWidth:  pw,
			PercentHeight: ph,
			Confidence:    r.FaceConfidence,
		}
	}

	return FaceResult{
		EngineName: d.String(),
		Faces:      faces,
	}, nil
}
