package engine

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg"
)

type FaceResult struct {
	EngineName string     `json:"engine"`
	Faces      []FaceData `json:"faces"`
}

func (r FaceResult) HasFaces() bool {
	return len(r.Faces) != 0
}

func (r FaceResult) ShowOutput() string {
	faceCount := len(r.Faces)
	byt, _ := json.Marshal(r)
	return fmt.Sprintf("%d\t%s", faceCount, string(byt))
}
