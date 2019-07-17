package dlib

import (
	"errors"

	"github.com/Kagami/go-face"

	"github.com/evalphobia/face-detect-annotator/engine"
)

type Config interface {
	GetDlibModelDir() string
}

type DlibFaceDetector struct {
	recognizer *face.Recognizer
}

func (d *DlibFaceDetector) Init(conf engine.Config) error {
	c, ok := conf.(Config)
	if !ok {
		return errors.New("Incompatible config type for DlibFaceDetector")
	}

	r, err := face.NewRecognizer(c.GetDlibModelDir())
	if err != nil {
		return err
	}

	d.recognizer = r
	return nil
}

func (d DlibFaceDetector) String() string {
	return "dlib"
}

func (d DlibFaceDetector) Detect(imgPath string) (engine.FaceResult, error) {
	imgWidth, imgHeight, err := engine.GetImageSize(imgPath)
	if err != nil {
		return engine.FaceResult{}, err
	}

	rects, err := d.recognizer.RecognizeFile(imgPath)
	faces := make([]engine.FaceData, len(rects))
	for i, rect := range rects {
		r := rect.Rectangle
		x := r.Min.X
		y := r.Min.Y
		w := r.Dx()
		h := r.Dy()

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
