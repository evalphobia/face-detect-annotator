package pigo

import (
	"errors"
	"io/ioutil"
	"sync"

	pigo "github.com/esimov/pigo/core"

	"github.com/evalphobia/face-detect-annotator/engine"
)

type Config interface {
	GetPigoCascadeFile() string
}

type PigoFaceDetector struct {
	mu         sync.Mutex
	classifier *pigo.Pigo

	angle        float64
	iouThreshold float64
	minSize      int
	maxSize      int
	shiftFactor  float64
	scaleFactor  float64
	qThresh      float32
}

func (d *PigoFaceDetector) Init(conf engine.Config) error {
	c, ok := conf.(Config)
	if !ok {
		return errors.New("Incompatible config type for PigoFaceDetector")
	}

	cascadeFile, err := ioutil.ReadFile(c.GetPigoCascadeFile())
	if err != nil {
		return err
	}

	classifier, err := pigo.NewPigo().Unpack(cascadeFile)
	if err != nil {
		return err
	}

	d.classifier = classifier
	d.angle = 0.0
	d.iouThreshold = 0.2
	d.minSize = 20
	d.maxSize = 1000
	d.shiftFactor = 0.1
	d.scaleFactor = 1.1
	d.qThresh = 5.0
	return nil
}

func (d PigoFaceDetector) String() string {
	return "pigo"
}

func (d PigoFaceDetector) Detect(imgPath string) (engine.FaceResult, error) {
	imgWidth, imgHeight, err := engine.GetImageSize(imgPath)
	if err != nil {
		return engine.FaceResult{}, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	src, err := pigo.GetImage(imgPath)
	if err != nil {
		return engine.FaceResult{}, err
	}

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y
	dets := d.classifier.RunCascade(pigo.CascadeParams{
		MinSize:     d.minSize,
		MaxSize:     d.maxSize,
		ShiftFactor: d.shiftFactor,
		ScaleFactor: d.scaleFactor,
		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}, d.angle)

	dets = d.classifier.ClusterDetections(dets, d.iouThreshold)

	faces := make([]engine.FaceData, 0, len(dets))
	for _, det := range dets {
		if det.Q < d.qThresh {
			continue
		}
		x := det.Col - det.Scale/2
		y := det.Row - det.Scale/2
		w := det.Scale
		h := det.Scale

		faces = append(faces, engine.FaceData{
			X:             x,
			Y:             y,
			Width:         w,
			Height:        h,
			PercentWidth:  float64(w) / float64(imgWidth),
			PercentHeight: float64(h) / float64(imgHeight),
			Confidence:    float64(det.Q),
		})
	}

	return engine.FaceResult{
		EngineName: d.String(),
		Faces:      faces,
	}, nil
}
