package fda

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/mkideal/cli"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"

	"github.com/evalphobia/face-detect-annotator/engine"
)

var (
	colorRed  = color.RGBA{255, 0, 0, 255}
	colorBlue = color.RGBA{0, 0, 255, 255}
)

var (
	fontFace18 font.Face
	fontFace36 font.Face
)

func init() {
	f, err := truetype.Parse(gobold.TTF)
	if err != nil {
		panic("[ERROR] cannot load font gobold.TTF")
	}

	fontFace18 = truetype.NewFace(f, &truetype.Options{
		Size: 18,
	})

	fontFace36 = truetype.NewFace(f, &truetype.Options{
		Size: 36,
	})
}

// annotator command
type annotatorT struct {
	cli.Helper
	Input string `cli:"*i,input" usage:"detector's output tsv file --input='/path/to/output.tsv'"`
}

var annotator = &cli.Command{
	Name: "annotate",
	Desc: "Annotate faces of image from --input TSV file",
	Argv: func() interface{} { return new(annotatorT) },
	Fn:   execAnnotator,
}

func execAnnotator(ctx *cli.Context) error {
	argv := ctx.Argv().(*annotatorT)

	f, err := NewCSVHandler(argv.Input)
	if err != nil {
		return err
	}

	lines, err := f.ReadAll()
	if err != nil {
		return err
	}

	const colSuffix = ":detail"
	var engines []string
	for _, h := range f.header {
		if strings.HasSuffix(h, colSuffix) {
			engines = append(engines, strings.TrimSuffix(h, colSuffix))
		}
	}
	fmt.Printf("engines:%+v\n", engines)

	for _, line := range lines {
		rawJsonData := make([]string, len(engines))
		for i, e := range engines {
			rawJsonData[i] = line[e+colSuffix]
		}
		imgPath := line["path"]
		err := annotateImage(imgPath, rawJsonData...)
		if err != nil {
			fmt.Printf("[ERROR] path:%s\terr:%s\n", imgPath, err.Error())
		}
	}
	return nil
}

func annotateImage(path string, targets ...string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	srcImg, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := srcImg.Bounds()
	images := make([]image.Image, len(targets))
	for i, rawJsonBody := range targets {
		img := image.NewRGBA(bounds)
		draw.Draw(img, bounds, srcImg, image.Pt(0, 0), draw.Src)
		images[i] = img

		data := engine.FaceResult{}
		err := json.Unmarshal([]byte(rawJsonBody), &data)
		if err != nil {
			fmt.Printf("[ERROR] JSON path:%s\t\terr:%s\n", path, err.Error())
			continue
		}
		drawString(img, image.Pt(10, 40), colorBlue, fontFace36, data.EngineName)

		if !data.HasFaces() {
			continue
		}
		for _, f := range data.Faces {
			drawRectBounds(img, image.Rect(f.X, f.Y, f.MaxX(), f.MaxY()), colorRed)
			var strList []string
			if f.Confidence > 0 {
				score := strconv.FormatFloat(f.Confidence, 'f', 2, 64)
				strList = append(strList, fmt.Sprintf("[%s%%]", score))
			}
			if f.PercentWidth > 0 || f.PercentHeight > 0 {
				pw := strconv.FormatFloat(f.PercentWidth, 'f', 2, 64)
				ph := strconv.FormatFloat(f.PercentHeight, 'f', 2, 64)
				strList = append(strList, fmt.Sprintf("[W:%s,H:%s]", pw, ph))
			}
			if len(strList) != 0 {
				drawString(img, image.Pt(f.X, f.Y-8), colorRed, fontFace18, strings.Join(strList, " "))

			}
		}
	}

	height := bounds.Max.Y
	bounds.Max.Y = height * len(images)
	canvas := image.NewRGBA(bounds)
	for i, img := range images {
		bounds.Min.Y = height * i
		bounds.Max.Y = height * (i + 1)
		draw.Draw(canvas, bounds, img, image.Pt(0, 0), draw.Src)
	}

	out, err := os.Create(getAnnotatedPath(path))
	if err != nil {
		return err
	}
	defer out.Close()

	err = jpeg.Encode(out, canvas, nil)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func getAnnotatedPath(p string) string {
	const prefix = "_annotated_"
	dir, file := filepath.Split(p)
	return filepath.Join(dir, prefix+file)
}

func drawString(img *image.RGBA, p image.Point, c color.Color, f font.Face, s string) {
	point := fixed.Point26_6{fixed.Int26_6(p.X * 64), fixed.Int26_6(p.Y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: f,
		Dot:  point,
	}
	d.DrawString(s)
}

func drawRectBounds(img *image.RGBA, r image.Rectangle, c color.Color) {
	minX, maxX := r.Min.X, r.Max.X
	minY, maxY := r.Min.Y, r.Max.Y

	// write lines of top and bottom
	for x := minX; x <= maxX; x++ {
		img.Set(x, minY-1, c)
		img.Set(x, minY, c)
		img.Set(x, maxY, c)
		img.Set(x, maxY+1, c)
	}

	// write lines of left and right
	for y := minY; y <= maxY; y++ {
		img.Set(minX-1, y, c)
		img.Set(minX, y, c)
		img.Set(maxX, y, c)
		img.Set(maxX+1, y, c)
	}
}
