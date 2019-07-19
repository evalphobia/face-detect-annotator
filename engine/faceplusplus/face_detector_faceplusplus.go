package faceplusplus

import (
	"github.com/evalphobia/go-face-plusplus/config"
	"github.com/evalphobia/go-face-plusplus/face"

	"github.com/evalphobia/face-detect-annotator/engine"
)

type FacePlusPlusFaceDetector struct {
	client *face.FaceService
}

func (d *FacePlusPlusFaceDetector) Init(_ engine.Config) error {
	svc, err := face.New(config.Config{})
	if err != nil {
		return err
	}

	d.client = svc
	return nil
}

func (d FacePlusPlusFaceDetector) String() string {
	return "face++"
}

func (d FacePlusPlusFaceDetector) Detect(imgPath string) (engine.FaceResult, error) {
	emptyResult := engine.FaceResult{}
	imgWidth, imgHeight, err := engine.GetImageSize(imgPath)
	if err != nil {
		return emptyResult, err
	}

	resp, err := d.client.DetectFromFile(imgPath)
	if err != nil {
		return emptyResult, err
	}

	faces := make([]engine.FaceData, len(resp.Faces))
	for i, f := range resp.Faces {
		r := f.FaceRectangle
		x := r.Left
		y := r.Top
		w := r.Width
		h := r.Height

		faces[i] = engine.FaceData{
			X:             x,
			Y:             y,
			Width:         w,
			Height:        h,
			PercentWidth:  float64(w) / float64(imgWidth),
			PercentHeight: float64(h) / float64(imgHeight),
		}
	}

	return engine.FaceResult{
		EngineName: d.String(),
		Faces:      faces,
	}, nil
}
