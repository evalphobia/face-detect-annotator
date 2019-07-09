package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"os"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
	"golang.org/x/image/bmp"
)

type TensorFlowFaceDetector struct {
	graph   *tf.Graph
	session *tf.Session
}

func NewTensorFlowFaceDetector(modelPath string) (*TensorFlowFaceDetector, error) {
	model, err := ioutil.ReadFile(modelPath)
	if err != nil {
		return nil, err
	}

	graph := tf.NewGraph()
	err = graph.Import(model, "")
	if err != nil {
		return nil, err
	}

	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}

	return &TensorFlowFaceDetector{
		graph:   graph,
		session: session,
	}, nil
}

func (d TensorFlowFaceDetector) String() string {
	return "tensorflow"
}

func (d TensorFlowFaceDetector) Detect(imgPath string) (FaceResult, error) {
	imgWidth, imgHeight, err := GetImageSize(imgPath)
	if err != nil {
		return FaceResult{}, err
	}

	tensor, err := makeTensorFromFile(imgPath)
	if err != nil {
		return FaceResult{}, err
	}

	results, err := d.detectFaces(tensor)
	if err != nil {
		return FaceResult{}, err
	}

	faces := make([]FaceData, len(results))
	for i, r := range results {
		x := r.PercentMinX * float64(imgWidth)
		y := r.PercentMinY * float64(imgHeight)
		w := (r.PercentMaxX - r.PercentMinX) * float64(imgWidth)
		h := (r.PercentMaxY - r.PercentMinY) * float64(imgHeight)
		pw := float64(r.PercentMaxX - r.PercentMinX)
		ph := float64(r.PercentMaxY - r.PercentMinY)

		faces[i] = FaceData{
			X:             int(x),
			Y:             int(y),
			Width:         int(w),
			Height:        int(h),
			PercentWidth:  pw,
			PercentHeight: ph,
			Confidence:    r.Score * 100,
		}
	}

	return FaceResult{
		EngineName: d.String(),
		Faces:      faces,
	}, nil
}

func makeTensorFromFile(imgPath string) (*tf.Tensor, error) {
	f, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = bmp.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	tensor, err := tf.NewTensor(buf.String())
	if err != nil {
		return nil, err
	}
	normalizeGraph, input, output, err := decodeBitmapGraph()
	if err != nil {
		return nil, err
	}
	normalizeSession, err := tf.NewSession(normalizeGraph, nil)
	if err != nil {
		return nil, err
	}
	defer normalizeSession.Close()
	normalized, err := normalizeSession.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, err
	}

	return normalized[0], nil
}

func decodeBitmapGraph() (*tf.Graph, tf.Output, tf.Output, error) {
	s := op.NewScope()
	input := op.Placeholder(s, tf.String)
	output := op.ExpandDims(
		s,
		op.DecodeBmp(s, input, op.DecodeBmpChannels(3)),
		op.Const(s.SubScope("make_batch"), int32(0)))
	graph, err := s.Finalize()
	return graph, input, output, err
}

func (d TensorFlowFaceDetector) detectFaces(tensor *tf.Tensor) ([]tfFace, error) {
	session := d.session
	graph := d.graph

	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			graph.Operation("image_tensor").Output(0): tensor,
		},
		[]tf.Output{
			graph.Operation("detection_boxes").Output(0),
			graph.Operation("detection_scores").Output(0),
			// graph.Operation("detection_classes").Output(0),
			// graph.Operation("num_detections").Output(0),
		},
		nil)
	if err != nil {
		return nil, fmt.Errorf("Error running session: %v", err)
	}

	boxes := output[0].Value().([][][]float32)[0]
	scores := output[1].Value().([][]float32)[0]
	results := make([]tfFace, 0, len(scores))
	for i, score := range scores {
		box := boxes[i]
		f := tfFace{
			PercentMinY: float64(box[0]),
			PercentMinX: float64(box[1]),
			PercentMaxY: float64(box[2]),
			PercentMaxX: float64(box[3]),
			Score:       float64(score),
		}
		if f.hasConfidence() {
			results = append(results, f)
		}
	}
	return results, nil
}

type tfFace struct {
	PercentMinX float64
	PercentMinY float64
	PercentMaxX float64
	PercentMaxY float64
	Score       float64
}

func (f tfFace) hasConfidence() bool {
	const confidenceBorder = 0.5
	return f.Score > confidenceBorder
}
