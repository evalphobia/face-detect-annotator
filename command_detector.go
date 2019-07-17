package fda

import (
	"fmt"
	"strings"
	"sync"

	"github.com/evalphobia/face-detect-annotator/engine"
	"github.com/mkideal/cli"
	"github.com/pkg/errors"
)

// detector command
type detectorT struct {
	cli.Helper
	Input        string `cli:"*i,input" usage:"image dir path --input='/path/to/image_dir'"`
	Output       string `cli:"*o,output" usage:"output TSV file path --output='./output.tsv'" dft:"./output.tsv"`
	UseAllEngine bool   `cli:"a,all" usage:"use all engines"`
	Engines      string `cli:"e,engine" usage:"comma separate Face Detect Engines --engine='opencv,dlib,tensorflow,rekognition,google,azure'" dft:"opencv,dlib,tensorflow"`
}

var detector = &cli.Command{
	Name: "detect",
	Desc: "Detect faces from image file or csv list",
	Argv: func() interface{} { return new(detectorT) },
	Fn:   execDetector,
}

var (
	baseDir    string
	pathPrefix string
)

func execDetector(ctx *cli.Context) error {
	argv := ctx.Argv().(*detectorT)
	conf := NewConfig()
	conf.setInputPath(argv.Input)
	conf.setOutputPath(argv.Output)
	for _, e := range strings.Split(argv.Engines, ",") {
		if err := conf.setUseEngineFromName(e); err != nil {
			return errors.Wrap(err, "[ERROR] setUseEngineFromName")
		}
	}

	engines, err := initEngines(conf, enabledEngines)
	if err != nil {
		return errors.Wrap(err, "[ERROR] initEngines")
	}

	switch {
	case conf.isCSVFilePath():
		return detectFromCSV(engines, conf)
	default:
		return detectFromImage(engines, conf)
	}
}

func detectFromImage(engines []engine.Engine, conf Config) error {
	for _, e := range engines {
		faceResult, err := e.Detect(conf.InputPath)
		if err != nil {
			return fmt.Errorf("[ERROR] %s\n", err.Error())
		}
		fmt.Printf("%s\t%s\n", e, faceResult.ShowOutput())
	}
	return nil
}

func detectFromCSV(engines []engine.Engine, conf Config) error {
	f, err := NewCSVHandler(conf.InputPath)
	if err != nil {
		return err
	}

	w, err := NewFileHandler(conf.OutputPath)
	if err != nil {
		return err
	}

	maxReq := make(chan struct{}, 10)

	lines, err := f.ReadAll()
	if err != nil {
		return err
	}

	result := make([]string, len(lines))
	var wg sync.WaitGroup
	for i, line := range lines {
		wg.Add(1)
		go func(i int, line map[string]string) {
			maxReq <- struct{}{}
			defer func() {
				<-maxReq
				wg.Done()
			}()

			fmt.Printf("exec #: [%d]\n", i)
			imgPath := line["path"]

			row := make([]string, len(engines)+2)
			row[0] = imgPath
			row[1] = line["count"]
			for i, e := range engines {
				faceResult, err := e.Detect(imgPath)
				if err != nil {
					fmt.Printf("[ERROR] %s\n", err.Error())
					continue
				}
				row[i+2] = faceResult.ShowOutput()

			}
			result[i] = strings.Join(row, "\t")
		}(i, line)
	}
	wg.Wait()

	result = append([]string{getHeader(engines)}, result...)
	return w.WriteAll(result)
}

func getHeader(engines []engine.Engine) string {
	header := []string{
		"path",
		"count",
	}
	for _, e := range engines {
		s := e.String()
		header = append(header, s+":count", s+":detail")
	}
	return strings.Join(header, "\t")
}
