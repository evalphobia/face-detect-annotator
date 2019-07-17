package rekognition

import (
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/rekognition"

	"github.com/evalphobia/face-detect-annotator/engine"
)

type RekognitionFaceDetector struct {
	client *rekognition.Rekognition
}

func (d *RekognitionFaceDetector) Init(_ engine.Config) error {
	cli, err := rekognition.New(config.Config{})
	if err != nil {
		return err
	}

	d.client = cli
	return nil
}

func (d RekognitionFaceDetector) String() string {
	return "rekognition"
}

func (d RekognitionFaceDetector) Detect(imgPath string) (engine.FaceResult, error) {
	emptyResult := engine.FaceResult{}
	imgWidth, imgHeight, err := engine.GetImageSize(imgPath)
	if err != nil {
		return emptyResult, err
	}

	resp, err := d.client.DetectFacesFromLocalFile(imgPath)
	if err != nil {
		return emptyResult, err
	}

	faces := make([]engine.FaceData, len(resp.List))
	for i, r := range resp.List {
		x := r.BoundingLeft * float64(imgWidth)
		y := r.BoundingTop * float64(imgHeight)
		pw := r.BoundingWidth
		ph := r.BoundingHeight

		faces[i] = engine.FaceData{
			X:             int(x),
			Y:             int(y),
			Width:         int(float64(imgWidth) * pw),
			Height:        int(float64(imgHeight) * ph),
			PercentWidth:  pw,
			PercentHeight: ph,
			Confidence:    r.FaceConfidence,
		}
	}

	return engine.FaceResult{
		EngineName: d.String(),
		Faces:      faces,
	}, nil
}
