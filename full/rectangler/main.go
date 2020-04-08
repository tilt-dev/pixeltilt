package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/windmilleng/pixeltilt/render/api"
)

type data struct {
	Status      string `json:"status"`
	Predictions []struct {
		LabelID      string    `json:"label_id"`
		Label        string    `json:"label"`
		Probability  float64   `json:"probability"`
		DetectionBox []float64 `json:"detection_box"`
	} `json:"predictions"`
}

func main() {
	Benchmark()
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	HandleWithNoSubpath("/", api.HttpRenderHandler(render))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func HandleWithNoSubpath(path string, f func(http.ResponseWriter, *http.Request)) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != path {
			http.NotFound(w, req)
			return
		}
		f(w, req)
	}

	http.HandleFunc(path, handler)
}

func render(req api.RenderRequest) (api.RenderReply, error) {
	// Run detection on original image
	detected, err := sendPostRequest("http://max-object-detector:5000/model/predict?threshold=0.7", "image", req.OriginalImage)
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "max-object-detector")
	}

	data := data{}
	err = json.Unmarshal(detected, &data)
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "max-object-detector")
	}

	input, err := png.Decode(bytes.NewReader(req.Image))
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "decoding image")
	}

	overlayImg, err := rectangler(data, input.Bounds().Size())
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "rectangling")
	}

	originalImage, err := png.Decode(bytes.NewReader(req.Image))
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "decoding image")
	}

	merged := merge(originalImage, overlayImg)
	var buf bytes.Buffer
	err = png.Encode(&buf, merged)
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "encoding image")
	}

	return api.RenderReply{Image: buf.Bytes()}, nil
}

// type data struct {
// 	Status      string `json:"status"e`
// 	Predictions []struct {
// 		LabelID      string `json:"label_id"`
// 		Label        string `json:"label"`
// 		Probability  int    `json:"probability"`
// 		DetectionBox []int  `json:"detection_box"`
// 	} `json:"predictions"`
// }

// [ymin, xmin, ymax, xmax]
func rectangler(data data, pt image.Point) (*image.RGBA, error) {
	pic := gg.NewContext(pt.X, pt.Y)
	pic.SetRGBA(0, 0, 0, 0)
	pic.Clear()
	w, h := float64(pt.X), float64(pt.Y)

	fontfile, err := ioutil.ReadFile("Inconsolata-Bold.ttf")
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(fontfile)
	if err != nil {
		return nil, err
	}
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
	}

	return pic.Image().(*image.RGBA), nil
}

func sendPostRequest(url string, name string, image []byte) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", name)
	if err != nil {
		return nil, err
	}

	_, err = part.Write(image)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("accept", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func merge(img1, img2 image.Image) *image.RGBA {
	img3 := image.NewRGBA(img1.Bounds())
	draw.Draw(img3, img3.Bounds(), img1, image.ZP, draw.Src)
	draw.Draw(img3, img3.Bounds(), img2, image.ZP, draw.Over)
	return img3
}
