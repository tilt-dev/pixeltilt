package main

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sug0/go-glitch"
)

func TestRender(t *testing.T) {
	input, err := os.Open("../frontend/public/baby-bear.png")
	if err != nil {
		panic(err)
	}
	wanted, err := ioutil.ReadFile("test/glitch.png")
	if err != nil {
		panic(err)
	}
	pic, err := png.Decode(input)

	glitchExpr, err = glitch.CompileExpression("(255 - (140 ? c)) & c")
	if err != nil {
		panic(err)
	}
	new, err := glitchExpr.JumblePixels(pic)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, new)
	if err != nil {
		panic(err)
	}
	if bytes.Compare(buf.Bytes(), wanted) != 0 {
		t.Errorf("glitch filter got != wanted")
	}
}
