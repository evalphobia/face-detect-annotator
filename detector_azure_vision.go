package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/computervision"
	"github.com/Azure/go-autorest/autorest"
)

type AzureVisionFaceDetector struct {
	client computervision.BaseClient
}

func NewAzureVisionFaceDetector(region, subscriptionKey string) (*AzureVisionFaceDetector, error) {
	endpoint := fmt.Sprintf("https://%s.api.cognitive.microsoft.com", region)
	cli := computervision.New(endpoint)
	authorizer := autorest.NewCognitiveServicesAuthorizer(subscriptionKey)
	cli.Authorizer = authorizer

	return &AzureVisionFaceDetector{
		client: cli,
	}, nil
}

func (d AzureVisionFaceDetector) String() string {
	return "azure"
}

func (d AzureVisionFaceDetector) Detect(imgPath string) (FaceResult, error) {
	imgWidth, imgHeight, err := GetImageSize(imgPath)
	if err != nil {
		return FaceResult{}, err
	}

	f, err := os.Open(imgPath)
	if err != nil {
		return FaceResult{}, err
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
		return FaceResult{}, err
	}

	b, _ := json.Marshal(resp)
	fmt.Printf("======= resp\n%s\n\n", string(b))
	if resp.Faces == nil {
		return FaceResult{
			EngineName: d.String(),
		}, nil
	}

	respFaces := *resp.Faces
	faces := make([]FaceData, len(respFaces))
	for i, f := range respFaces {
		r := f.FaceRectangle
		x := *r.Left
		y := *r.Top
		w := *r.Width
		h := *r.Height

		faces[i] = FaceData{
			X:             int(x),
			Y:             int(y),
			Width:         int(w),
			Height:        int(h),
			PercentWidth:  float64(w) / float64(imgWidth),
			PercentHeight: float64(h) / float64(imgHeight),
			// Confidence:    r.FaceConfidence,
		}
	}

	return FaceResult{
		EngineName: d.String(),
		Faces:      faces,
	}, nil
}
