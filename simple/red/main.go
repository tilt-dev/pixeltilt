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
	fmt.Println("Red running!")
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
	fmt.Println("red decoded ok")

	newImg := red(pic)
	var buf bytes.Buffer
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
				// R: uint8(float64(pixel.R) * 0.8),
				// G: 0,
				R: 0,
				G: uint8(float64(pixel.G) * 0.8),
				B: 0,
				A: pixel.A,
			}
			newPic.Set(x, y, new)
		}
	}

	return newPic
}
