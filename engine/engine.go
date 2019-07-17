package engine

type Engine interface {
	Init(conf Config) error
	String() string
	Detect(imgPath string) (FaceResult, error)
}

type Config interface{}
