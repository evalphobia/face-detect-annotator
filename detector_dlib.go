package main

import (
	"github.com/Kagami/go-face"
)

type DlibFaceDetector struct {
	recognizer *face.Recognizer
}

func NewDlibFaceDetector(modelDir string) (*DlibFaceDetector, error) {
	r, err := face.NewRecognizer(modelDir)
	if err != nil {
		return nil, err
	}

	return &DlibFaceDetector{
		recognizer: r,
	}, nil
}

func (d DlibFaceDetector) String() string {
	return "dlib"
}

func (d DlibFaceDetector) Detect(imgPath string) (FaceResult, error) {
	imgWidth, imgHeight, err := GetImageSize(imgPath)
	if err != nil {
		return FaceResult{}, err
	}

	rects, err := d.recognizer.RecognizeFile(imgPath)
	faces := make([]FaceData, len(rects))
	for i, rect := range rects {
		r := rect.Rectangle
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
