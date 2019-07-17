package fda

import (
	"fmt"

	"github.com/evalphobia/face-detect-annotator/engine"
	"github.com/pkg/errors"
)

var enabledEngines []engine.Engine

func AddEngines(engine ...engine.Engine) {
	enabledEngines = append(enabledEngines, engine...)
}

func initEngines(conf Config, enabledEngines []engine.Engine) ([]engine.Engine, error) {
	engines := make([]engine.Engine, 0, len(enabledEngines))
	for _, e := range enabledEngines {
		switch {
		case e.String() == "azure" && conf.UseEngineAzureVision,
			e.String() == "google" && conf.UseEngineGoogleVision,
			e.String() == "rekognition" && conf.UseEngineRekognition,
			e.String() == "dlib" && conf.UseEngineDlib,
			e.String() == "pigo" && conf.UseEnginePigo,
			e.String() == "opencv" && conf.UseEngineOpenCV,
			e.String() == "tensorflow" && conf.UseEngineTensorFlow:
			err := e.Init(conf)
			if err != nil {
				return nil, err
			}
			engines = append(engines, e)
		}
	}

	if len(engines) == 0 {
		return nil, errors.New("Any face detect engine is specified")
	}

	for _, e := range engines {
		fmt.Printf("[INFO] Use %s\n", e.String())
	}

	return engines, nil
}
