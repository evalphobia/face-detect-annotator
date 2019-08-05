package fda

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/mkideal/cli"
)

// list command
type listT struct {
	cli.Helper
	Input          string `cli:"*i,input" usage:"image dir path --input='/path/to/image_dir'"`
	Output         string `cli:"*o,output" usage:"output CSV file path --output='./list.csv'" dft:"./list.csv"`
	IncludeAllType bool   `cli:"a,all" usage:"use all files"`
	Type           string `cli:"t,type" usage:"comma separate file extensions --type='jpg,jpeg,png,gif'" dft:"jpg,jpeg,png,gif"`
	PathPrefix     string `cli:"d,prefix" usage:"prefix for file path --prefix='/tmp'" dft:""`
}

var list = &cli.Command{
	Name: "list",
	Desc: "Find image files in --input directory and save it to csv file.",
	Argv: func() interface{} { return new(listT) },
	Fn:   execList,
}

func execList(ctx *cli.Context) error {
	argv := ctx.Argv().(*listT)

	f, err := NewFileHandler(argv.Output)
	if err != nil {
		return err
	}

	types := newFileType(strings.Split(argv.Type, ","))
	if argv.IncludeAllType {
		types.setIncludeAll(argv.IncludeAllType)
	}

	pathPrefix = argv.PathPrefix
	baseDir = fmt.Sprintf("%s/", filepath.Clean(argv.Input))
	result := getFilesFromDir(baseDir, types)
	result = append([]string{"path"}, result...)
	return f.WriteAll(result)
}

func getFilesFromDir(dir string, types fileType) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() {
			paths = append(paths, getFilesFromDir(filepath.Join(dir, fileName), types)...)
			continue
		}

		if !types.isTarget(fileName) {
			continue
		}

		path := path.Join(pathPrefix, dir, fileName)
		paths = append(paths, path)
	}

	return paths
}
