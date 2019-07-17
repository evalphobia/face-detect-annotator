package google

import (
	"io/ioutil"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/vision"
	SDK "google.golang.org/api/vision/v1"

	"github.com/evalphobia/face-detect-annotator/engine"
)

type GoogleVisionFaceDetector struct {
	client *vision.Vision
}

func (d *GoogleVisionFaceDetector) Init(_ engine.Config) error {
	cli, err := vision.New(config.Config{})
	if err != nil {
		return err
	}

	d.client = cli
	return nil
}

func (d GoogleVisionFaceDetector) String() string {
	return "google"
}

func (d GoogleVisionFaceDetector) Detect(imgPath string) (engine.FaceResult, error) {
	emptyResult := engine.FaceResult{}
	imgWidth, imgHeight, err := engine.GetImageSize(imgPath)
	if err != nil {
		return emptyResult, err
	}

	byt, err := ioutil.ReadFile(imgPath)
	if err != nil {
		return emptyResult, err
	}

	resp, err := d.client.Face(byt)
	if err != nil {
		return emptyResult, err
	}

	var list []*SDK.FaceAnnotation
	for _, r := range resp.Responses {
		list = append(list, r.FaceAnnotations...)
	}
	faces := make([]engine.FaceData, len(list))
	for i, r := range list {
		v := r.FdBoundingPoly.Vertices
		if len(v) < 3 {
			continue
		}
		x := v[0].X
		y := v[0].Y
		w := v[2].X - v[0].X
		h := v[2].Y - v[0].Y

		faces[i] = engine.FaceData{
			X:             int(x),
			Y:             int(y),
			Width:         int(w),
			Height:        int(h),
			PercentWidth:  float64(w) / float64(imgWidth),
			PercentHeight: float64(h) / float64(imgHeight),
			Confidence:    r.DetectionConfidence * 100,
		}
	}

	return engine.FaceResult{
		EngineName: d.String(),
		Faces:      faces,
	}, nil
}
