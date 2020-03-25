package main

import (
	"bytes"
	"fmt"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sug0/go-glitch"
	"github.com/windmilleng/enhance/render/api"

	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// syntax described here: https://github.com/sug0/go-glitch/blob/master/res/cool.txt
// various equations here: https://github.com/sug0/go-glitch/blob/master/res/cool.txt
var glitchExprStr = "b ^ (r | (s : (x # B)))"
var glitchExpr *glitch.Expression

func main() {
	called := time.Unix(0, 1585136507678862367)
	current := time.Now()
	elapsed := current.Sub(called)
	fmt.Println("\nStarting glitch...")
	fmt.Println(elapsed.Round(time.Millisecond))

	port := "8085"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	var err error
	glitchExpr, err = glitch.CompileExpression(glitchExprStr)
	if err != nil {
		log.Fatal(err)
	}

	HandleWithNoSubpath("/", api.HttpRenderHandler(render))
	http.HandleFunc("/set_expr", setExpr)
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

// just a dumb form to let people play around with different expressions without restarting
func setExpr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write([]byte(`
<form method="post">
expr: <input type="text" name="expr">
<input type="submit" value="Submit">
</form>
`))
	if err != nil {
		log.Printf("error writing response: %v\n", err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing request: %v", err), http.StatusBadRequest)
		return
	}

	expr := r.PostForm.Get("expr")
	if expr != "" {
		g, err := glitch.CompileExpression(expr)
		if err != nil {
			_, err := w.Write([]byte(fmt.Sprintf("<br>Invalid expression: %s<br>", err.Error())))
			if err != nil {
				log.Printf("error writing response: %v\n", err)
				return
			}
		}
		glitchExpr = g
		log.Printf("changed expression to %s\n", expr)
		_, err = w.Write([]byte(fmt.Sprintf("<br>Expression set to %s<br>", expr)))
		if err != nil {
			log.Printf("error writing response: %v\n", err)
			return
		}
	}
}

func render(req api.RenderRequest) (api.RenderReply, error) {
	input, err := png.Decode(bytes.NewReader(req.Image))
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "decoding image")
	}

	output, err := glitchExpr.JumblePixels(input)
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "glitching image")
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, output)
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "encoding image")
	}

	return api.RenderReply{Image: buf.Bytes()}, nil
}
