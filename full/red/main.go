package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/windmilleng/pixeltilt/render/api"
)

func main() {
	Benchmark()
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	http.HandleFunc("/", api.HttpRenderHandler(render))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func render(req api.RenderRequest) (api.RenderReply, error) {
	pic, err := png.Decode(bytes.NewReader(req.Image))
	if err != nil {
		return api.RenderReply{}, err
	}

	// newImg := red(pic)
	newImg := scaleColor(pic, 1.0, 1.0, 0.0)
	var buf bytes.Buffer
	x := 0
	os.Exit(1)
	log.Printf("hrm %d", 5/x)
	err = png.Encode(&buf, newImg)
	if err != nil {
		return api.RenderReply{}, err
	}
	resp := api.RenderReply{
		Image: buf.Bytes(),
	}

	return resp, nil
}

func red(pic image.Image) *image.RGBA {
	picSize := pic.Bounds().Size()
	newPic := image.NewRGBA(image.Rect(0, 0, picSize.X, picSize.Y))

	for x := 0; x < picSize.X; x++ {
		for y := 0; y < picSize.Y; y++ {
			pixel := color.RGBAModel.Convert(pic.At(x, y)).(color.RGBA)
			new := color.RGBA{
				R: uint8(float64(pixel.R) * 0.8),
				G: uint8(float64(pixel.G) * 0.8),
				B: 0,
				A: pixel.A,
			}
			newPic.Set(x, y, new)
		}
	}

	return newPic
}

func scaleColor(pic image.Image, rScale float64, gScale float64, bScale float64) *image.RGBA {
	picSize := pic.Bounds().Size()
	newPic := image.NewRGBA(image.Rect(0, 0, picSize.X, picSize.Y))

	for x := 0; x < picSize.X; x++ {
		for y := 0; y < picSize.Y; y++ {
			pixel := color.RGBAModel.Convert(pic.At(x, y)).(color.RGBA)
			new := color.RGBA{
				R: uint8(float64(pixel.R) * rScale),
				G: uint8(float64(pixel.G) * gScale),
				B: uint8(float64(pixel.B) * bScale),
				A: pixel.A,
			}
			newPic.Set(x, y, new)
		}
	}

	return newPic
}
