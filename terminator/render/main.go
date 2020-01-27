package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

var jsonInput = []byte(`{"status":"ok","predictions":[{"label_id":"5","label":"spaceship","probability":0.993535578250885,"detection_box":[0.3047112226486206,0.1421125829219818,0.7982025146484375,0.907085657119751]}]}`)

type data struct {
	Status      string `json:"status"`
	Predictions []struct {
		LabelID      string    `json:"label_id"`
		Label        string    `json:"label"`
		Probability  float64   `json:"probability"`
		DetectionBox []float64 `json:"detection_box"`
	} `json:"predictions"`
}

func iferr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("\nStarting renderer...")
	port := "8084"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	http.HandleFunc("/", render)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

var myData = data{}

func render(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Render request received!")

	file, header, err := r.FormFile("image")
	iferr(err)
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\tFile Size: %+v\tMIME: %+v\n", header.Filename, header.Size, header.Header)
	fileBytes, err := ioutil.ReadAll(file)
	iferr(err)
	if header.Filename == "data" {
		err = json.Unmarshal(fileBytes, &myData)
		iferr(err)
		w.Write([]byte("ok"))
		return
	}
	if header.Filename == "pic" {
		reader := bytes.NewReader(fileBytes)
		pic, _, err := image.Decode(reader)
		iferr(err)
		red := red(pic)
		text := rectangler(myData, red.Bounds().Size())
		output := merge(red, text)
		w.Header().Set("Content-Type", "image/png")
		err = png.Encode(w, output)
		iferr(err)
	}
}

func merge(img1, img2 *image.RGBA) *image.RGBA {
	img3 := image.NewRGBA(img1.Bounds())
	draw.Draw(img3, img3.Bounds(), img1, image.ZP, draw.Src)
	draw.Draw(img3, img3.Bounds(), img2, image.ZP, draw.Over)
	return img3
}

func red(pic image.Image) *image.RGBA {
	picSize := pic.Bounds().Size()
	newPic := image.NewRGBA(image.Rect(0, 0, picSize.X, picSize.Y))

	for x := 0; x < picSize.X; x++ {
		for y := 0; y < picSize.Y; y++ {
			pixel := color.RGBAModel.Convert(pic.At(x, y)).(color.RGBA)
			new := color.RGBA{
				R: uint8(float64(pixel.R) * 0.8),
				G: 0,
				B: 0,
				A: pixel.A,
			}
			newPic.Set(x, y, new)
		}
	}

	// newFile, err := os.Create("output.png")
	// iferr(err)
	// defer newFile.Close()
	// err = png.Encode(newFile, newPic)
	// iferr(err)

	return newPic
}

// type data struct {
// 	Status      string `json:"status"`
// 	Predictions []struct {
// 		LabelID      string `json:"label_id"`
// 		Label        string `json:"label"`
// 		Probability  int    `json:"probability"`
// 		DetectionBox []int  `json:"detection_box"`
// 	} `json:"predictions"`
// }

// [ymin, xmin, ymax, xmax]
func rectangler(data data, pt image.Point) *image.RGBA {
	pic := gg.NewContext(pt.X, pt.Y)
	pic.SetRGBA(0, 0, 0, 0)
	pic.Clear()
	w, h := float64(pt.X), float64(pt.Y)

	fontfile, err := ioutil.ReadFile("render/Inconsolata-Bold.ttf")
	font, err := truetype.Parse(fontfile)
	// font, err := truetype.Parse(fontface.TTF)
	iferr(err)
	fontSize := (w + h) / 70
	face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
	pic.SetFontFace(face)

	for i := 0; i < len(data.Predictions); i++ {
		values := data.Predictions[i].DetectionBox
		ymin, xmin, ymax, xmax := values[0], values[1], values[2], values[3]
		pic.SetRGBA(1, 1, 1, 1)
		pic.SetLineWidth((w + h) / 500)

		pic.DrawRectangle(xmin*w, ymin*h, xmax*w-xmin*w, ymax*h-ymin*h)
		pic.Stroke()

		object := fmt.Sprintf("OBJECT:      %s", strings.ToUpper(data.Predictions[i].Label))
		pic.DrawString(object, xmin*w, ymin*h-fontSize/2-fontSize*1.2)
		probability := fmt.Sprintf("PROBABILITY: %f", data.Predictions[i].Probability)
		pic.DrawString(probability, xmin*w, ymin*h-fontSize/2)

		// pic.SavePNG("rect.png")
	}

	return pic.Image().(*image.RGBA)
}
