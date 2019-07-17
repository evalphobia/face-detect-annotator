package azure

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/computervision"
	"github.com/Azure/go-autorest/autorest"

	"github.com/evalphobia/face-detect-annotator/engine"
)

type Config interface {
	GetAzureRegion() string
	GetAzureSubscriptionKey() string
}

type AzureVisionFaceDetector struct {
	client computervision.BaseClient
}

func (d *AzureVisionFaceDetector) Init(conf engine.Config) error {
	c, ok := conf.(Config)
	if !ok {
		return errors.New("Incompatible config type for AzureVisionFaceDetector")
	}

	endpoint := fmt.Sprintf("https://%s.api.cognitive.microsoft.com", c.GetAzureRegion())
	cli := computervision.New(endpoint)
	authorizer := autorest.NewCognitiveServicesAuthorizer(c.GetAzureSubscriptionKey())
	cli.Authorizer = authorizer

	d.client = cli
	return nil
}

func (d AzureVisionFaceDetector) String() string {
	return "azure"
}

func (d AzureVisionFaceDetector) Detect(imgPath string) (engine.FaceResult, error) {
	emptyResult := engine.FaceResult{}
	imgWidth, imgHeight, err := engine.GetImageSize(imgPath)
	if err != nil {
		return emptyResult, err
	}

	f, err := os.Open(imgPath)
	if err != nil {
		return emptyResult, err
	}
	defer f.Close()

	ctx := context.Background()
	resp, err := d.client.AnalyzeImageInStream(
		ctx,
		f,
		[]computervision.VisualFeatureTypes{computervision.VisualFeatureTypesFaces},
		nil,
		"",
	)
	if err != nil {
		return emptyResult, err
	}

	if resp.Faces == nil {
		return engine.FaceResult{
			EngineName: d.String(),
		}, nil
	}

	respFaces := *resp.Faces
	faces := make([]engine.FaceData, len(respFaces))
	for i, f := range respFaces {
		r := f.FaceRectangle
		x := *r.Left
		y := *r.Top
		w := *r.Width
		h := *r.Height

		faces[i] = engine.FaceData{
			X:             int(x),
			Y:             int(y),
			Width:         int(w),
			Height:        int(h),
			PercentWidth:  float64(w) / float64(imgWidth),
			PercentHeight: float64(h) / float64(imgHeight),
			// Confidence:    r.FaceConfidence,
		}
	}

	return engine.FaceResult{
		EngineName: d.String(),
		Faces:      faces,
	}, nil
}
