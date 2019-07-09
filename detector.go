package main

type FaceDetector interface {
	String() string
	Detect(path string) (FaceResult, error)
}
