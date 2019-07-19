package main

import (
	fda "github.com/evalphobia/face-detect-annotator"
	"github.com/evalphobia/face-detect-annotator/engine/azure"
	"github.com/evalphobia/face-detect-annotator/engine/faceplusplus"
	"github.com/evalphobia/face-detect-annotator/engine/google"
	"github.com/evalphobia/face-detect-annotator/engine/pigo"
	"github.com/evalphobia/face-detect-annotator/engine/rekognition"
)

func main() {
	fda.AddEngines(
		&azure.AzureVisionFaceDetector{},
		&google.GoogleVisionFaceDetector{},
		&rekognition.RekognitionFaceDetector{},
		&faceplusplus.FacePlusPlusFaceDetector{},
		&pigo.PigoFaceDetector{})
	fda.Run()
}
