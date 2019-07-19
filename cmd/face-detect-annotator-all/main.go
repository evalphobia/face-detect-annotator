package main

import (
	fda "github.com/evalphobia/face-detect-annotator"
	"github.com/evalphobia/face-detect-annotator/engine/azure"
	"github.com/evalphobia/face-detect-annotator/engine/dlib"
	"github.com/evalphobia/face-detect-annotator/engine/faceplusplus"
	"github.com/evalphobia/face-detect-annotator/engine/google"
	"github.com/evalphobia/face-detect-annotator/engine/opencv"
	"github.com/evalphobia/face-detect-annotator/engine/pigo"
	"github.com/evalphobia/face-detect-annotator/engine/rekognition"
	"github.com/evalphobia/face-detect-annotator/engine/tensorflow"
)

func main() {
	fda.AddEngines(
		&azure.AzureVisionFaceDetector{},
		&google.GoogleVisionFaceDetector{},
		&rekognition.RekognitionFaceDetector{},
		&faceplusplus.FacePlusPlusFaceDetector{},
		&pigo.PigoFaceDetector{},
		&dlib.DlibFaceDetector{},
		&opencv.OpenCVFaceDetector{},
		&tensorflow.TensorFlowFaceDetector{})
	fda.Run()
}
