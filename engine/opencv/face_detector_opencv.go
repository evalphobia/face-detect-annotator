package opencv

import (
	"errors"
	"fmt"
	"sync"

	"gocv.io/x/gocv"

	"github.com/evalphobia/face-detect-annotator/engine"
)

type Config interface {
	GetOpenCVCascadeFile() string
}

type OpenCVFaceDetector struct {
	mu         sync.Mutex
	classifier gocv.CascadeClassifier
}

func (d *OpenCVFaceDetector) Init(conf engine.Config) error {
	c, ok := conf.(Config)
	if !ok {
		return errors.New("Incompatible config type for DlibFaceDetector")
	}

	classifier := gocv.NewCascadeClassifier()
	if !classifier.Load(c.GetOpenCVCascadeFile()) {
		return fmt.Errorf("Error reading cascade file: [%s]", c.GetOpenCVCascadeFile())
	}

	d.classifier = classifier
	return nil
}

func (d OpenCVFaceDetector) String() string {
	return "opencv"
}

func (d *OpenCVFaceDetector) Detect(imgPath string) (engine.FaceResult, error) {
	imgWidth, imgHeight, err := engine.GetImageSize(imgPath)
	if err != nil {
		return engine.FaceResult{}, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	img := gocv.IMRead(imgPath, gocv.IMReadGrayScale)
	if img.Empty() {
		return engine.FaceResult{}, errors.New("Empty image")
	}

	rects := d.classifier.DetectMultiScale(img)
	faces := make([]engine.FaceData, len(rects))
	for i, r := range rects {
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
